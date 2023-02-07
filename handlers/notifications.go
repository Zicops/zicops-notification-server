package handlers

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"strconv"

	//"os"

	"time"

	"encoding/json"

	"google.golang.org/api/iterator"

	//"github.com/zicops/contracts/notificationz"
	"github.com/segmentio/ksuid"
	"github.com/zicops/zicops-notification-server/global"
	"github.com/zicops/zicops-notification-server/graph/model"
	"github.com/zicops/zicops-notification-server/jwt"
)

type message struct {
	Notification skeleton `json:"notification"`
	To           string   `json:"to"`
	CreatedAt    int64    `json:"created_at"`
	Data         dat      `json:"data"`
}

type firebaseData struct {
	M     message
	LspID string
}

type dat struct {
	OpenUrl string `json:"openURL"`
}

type skeleton struct {
	Title        string `json:"title"`
	Body         string `json:"body"`
	Click_Action string `json:"click_action"`
}

type results struct {
	MessageId string `json:"message_id"`
}

type respBody struct {
	Multicast_id  int       `json:"multicast_id"`
	Success       int       `json:"success"`
	Failure       int       `json:"failure"`
	Canonical_ids int       `json:"canonical_id"`
	Results       []results `json:"results"`
}

// send notification without link
func SendNotification(ctx context.Context, notification model.NotificationInput) ([]*model.Notification, error) {
	global.Ct = ctx
	var res []*model.Notification

	//s := notificationz.Skeleton and so on
	s := skeleton{
		Title: notification.Title,
		Body:  notification.Body,
	}

	l := len(notification.UserID)
	var flag []int = make([]int, l)
	//now we need to get fcm-token for given email, i.e., from email we need userID and using that we will get fcm-token
	for k, ui := range notification.UserID {
		userId := ui

		var resp []map[string]interface{}
		//using this user id we will get fcm tokens
		iter := global.Client.Collection("tokens").Where("UserID", "==", userId).Documents(ctx)
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to iterate: %v", err)
				return nil, err
			}

			resp = append(resp, doc.Data())
		}
		//now we have all the instances where userID is of person in the mail, and their fcm token/tokens alongside
		for _, v := range resp {
			m := message{
				Notification: s,
				To:           v["FCM-token"].(string),
				CreatedAt:    time.Now().Unix(),
			}

			dataJson, _ := json.Marshal(m)

			code := sendToFirebase(dataJson)

			res = append(res, &model.Notification{
				Statuscode: strconv.Itoa(code),
			})
			if code == 1 {
				if flag[k] == 0 {
					//means value has not been added yet, add the value
					sendingToFirestore(dataJson, *userId)
					flag[k] = 1
				}
			}
		}

	}
	return res, nil

}

// send notification with link
func SendNotificationWithLink(ctx context.Context, notification model.NotificationInput, link string) ([]*model.Notification, error) {
	global.Ct = ctx

	//get claims from context
	claims, err := GetClaimsFromContext(ctx)
	if err != nil {
		log.Printf("Unable to get claims from context: %v", err)
		return nil, err
	}
	lsp := claims["lsp_id"].(string)

	var res []*model.Notification

	tmp := dat{
		OpenUrl: link,
	}
	//s := notificationz.Skeleton and so on
	s := skeleton{
		Title:        notification.Title,
		Body:         notification.Body,
		Click_Action: "action.open.zicops.with.url",
	}

	l := len(notification.UserID)
	var flag []int = make([]int, l)
	for k := range flag {
		flag[k] = 0
	}
	//now we need to get fcm-token for given email, i.e., from email we need userID and using that we will get fcm-token
	for k, userId := range notification.UserID {
		//userId := base64.StdEncoding.EncodeToString([]byte(*email))

		var resp []map[string]interface{}
		//using this user id we will get fcm tokens
		iter := global.Client.Collection("tokens").Where("UserID", "==", *userId).Where("LspID", "==", lsp).Documents(ctx)
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to iterate: %v", err)
				return nil, err
			}
			//check for null values
			tmp := doc.Data()
			if tmp["FCM-token"] == "null" || tmp["LspId"] == "null" {
				continue
			}

			resp = append(resp, doc.Data())
		}
		//now we have all the instances where userID is of person in the mail, and their fcm token/tokens alongside
		for _, v := range resp {
			m := message{
				Notification: s,
				To:           v["FCM-token"].(string),
				CreatedAt:    time.Now().Unix(),
				Data:         tmp,
			}

			dataJson, err := json.Marshal(m)

			if err != nil {
				log.Printf("Unable to convert to JSON: %v", err)
			}

			code := sendToFirebase(dataJson)
			res = append(res, &model.Notification{
				Statuscode: strconv.Itoa(code),
			})
			templink := dat{
				OpenUrl: link,
			}
			temp := message{
				Notification: s,
				To:           v["FCM-token"].(string),
				CreatedAt:    time.Now().Unix(),
				Data:         templink,
			}
			fbd := firebaseData{
				M:     temp,
				LspID: lsp,
			}
			tempJson, _ := json.Marshal(fbd)
			if code == 1 {
				if flag[k] == 0 {
					//means value has not been added yet, add the value
					sendingToFirestore(tempJson, *userId)
					flag[k] = 1
				}
			}
		}

	}
	return res, nil

}

func sendingToFirestore(dataJson []byte, userId string) {
	var msg firebaseData
	err := json.Unmarshal(dataJson, &msg)

	if err != nil {
		log.Println(err)
	}

	msgId := ksuid.New()
	tmp := msg.M.Data.OpenUrl
	_, err = global.Client.Collection("notification").Doc(msgId.String()).Set(global.Ct, model.FirestoreData{
		Title:     msg.M.Notification.Title,
		Body:      msg.M.Notification.Body,
		CreatedAt: int(time.Now().Unix()),
		MessageID: msgId.String(),
		UserID:    userId,
		IsRead:    false,
		Link:      &tmp,
		LspID:     msg.LspID,
	})

	if err != nil {
		log.Fatalf("Failed adding value to cloud firestore: %v", err)
	}
}

func sendToFirebase(dataJson []byte) int {

	body := bytes.NewReader(dataJson)
	//sending request
	req, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", body)
	if err != nil {
		log.Printf("Error while sending request %v", err)
	}

	req.Header.Set("Authorization", "key=AAAAU56xXYc:APA91bHtHX1hjkj8B4u0tSTuuTgURF6PvlqKEzgn3Qv7JR14mwra7rrCCg3bRRJZHxYyK8DHntk4Tc9CsXkqj44vuxFeD1WgRy1nifgbYgi60IAmfApLKK6Rd92Puuj3NPtUNGvdNHvr")
	req.Header.Set("Content-Type", "application/json")

	//getting response
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error while getting response- %v", err)
	}

	//converting response received to bytes
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to receive response to JSON")
	}

	//getting success value from response body
	var successCode respBody
	err = json.Unmarshal(b, &successCode)
	if err != nil {
		//it means that we we don't have data according to respBody struct, i.e., instead of message_id, there are errors
		log.Printf("Unable to send the notification %v", err)
	}
	return successCode.Success

}

func GetClaimsFromContext(ctx context.Context) (map[string]interface{}, error) {
	token := ctx.Value("token").(string)
	claims, err := jwt.GetClaims(token)
	if err != nil {
		return nil, err
	}
	//get lsp-id from context, if already there then okay otherwise put zicops lsp-id
	lspID := "d8685567-cdae-4ee0-a80e-c187848a760e"
	lsp := ctx.Value("tenant").(string)
	if lsp == "" {
		lsp = lspID
	}
	claims["lsp_id"] = lsp
	return claims, err
}
