package handlers

import (
	"context"
	"log"
	"strconv"

	"time"

	"firebase.google.com/go/messaging"
	"google.golang.org/api/iterator"

	"github.com/segmentio/ksuid"
	"github.com/zicops/zicops-notification-server/global"
	"github.com/zicops/zicops-notification-server/graph/model"
	"github.com/zicops/zicops-notification-server/jwt"
)

type message struct {
	Notification *messaging.Notification `json:"notification"`
	To           string                  `json:"to"`
	CreatedAt    int64                   `json:"created_at"`
	Data         string                  `json:"data"`
}

type firebaseData struct {
	M     message
	LspID string
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

	s := &messaging.Notification{
		Title: notification.Title,
		Body:  notification.Body,
	}

	l := len(notification.UserID)
	var flag []int = make([]int, l)
	for k := range flag {
		flag[k] = 0
	}
	//now we need to get fcm-token for given email, i.e., from email we need userID and using that we will get fcm-token
	for k, userId := range notification.UserID {
		//userId := base64.StdEncoding.EncodeToString([]byte(*email))
		if userId == nil {
			continue
		}
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
			if !checkNullValues(tmp) {
				continue
			}
			resp = append(resp, doc.Data())
		}
		//now we have all the instances where userID is of person in the mail, and their fcm token/tokens alongside
		for _, v := range resp {

			m := messaging.Message{
				Token:        v["FCM-token"].(string),
				Notification: s,
			}

			var e string
			code, err := sendToFirebase(m, ctx)
			if err != nil {
				e = err.Error()
			}

			if code == 0 {
				res = append(res, &model.Notification{
					Statuscode: strconv.Itoa(code),
					UserID:     userId,
					Error:      &e,
				})
			}
			if code == 1 {
				res = append(res, &model.Notification{
					Statuscode: strconv.Itoa(code),
					UserID:     userId,
					Error:      nil,
				})
			}

			temp := message{
				Notification: s,
				To:           v["FCM-token"].(string),
				CreatedAt:    time.Now().Unix(),
				Data:         link,
			}
			fbd := firebaseData{
				M:     temp,
				LspID: lsp,
			}

			if code == 1 {
				if flag[k] == 0 {
					//means value has not been added yet, add the value
					sendingToFirestore(fbd, *userId)
					flag[k] = 1
				}
			}
		}

	}
	return res, nil

}

func sendingToFirestore(msg firebaseData, userId string) {

	msgId := ksuid.New()
	tmp := msg.M.Data
	_, err := global.Client.Collection("notification").Doc(msgId.String()).Set(global.Ct, model.FirestoreData{
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

func sendToFirebase(message messaging.Message, ctx context.Context) (int, error) {

	//sending request
	v, err := global.Messanger.Send(ctx, &message)
	if err != nil {
		log.Printf("Got error: %v", err)
	}
	if len(v) == 0 {
		return 0, err
	}
	if v[:29] == "projects/zicops-one/messages/" {
		return 1, nil
	}

	return 0, nil

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

func checkNullValues(tmp map[string]interface{}) bool {
	if tmp["FCM-token"] == nil || tmp["FCM-token"].(string) == "null" {
		return false
	}

	if tmp["LspID"] == nil || tmp["LspID"].(string) == "null" {
		return false
	}

	return true
}
