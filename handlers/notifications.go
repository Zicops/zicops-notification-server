package handlers

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"strconv"

	//"os"
	"sync"
	"time"

	"encoding/json"

	"github.com/allegro/bigcache/v3"
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
}

type skeleton struct {
	Title string `json:"title"`
	Body  string `json:"body"`
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

var cache *bigcache.BigCache

func SendNotification(ctx context.Context, notification model.NotificationInput) ([]*model.Notification, error) {
	global.Ct = ctx
	var res []*model.Notification

	//channel for sending data to cache function
	ch := make(chan []byte, 100)
	var mut sync.Mutex

	cacheVar, err := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	if err != nil {
		log.Printf("Unable to create cache %v", err)
	}
	cache = cacheVar

	//s := notificationz.Skeleton and so on
	s := skeleton{
		Title: notification.Title,
		Body:  notification.Body,
	}

	l := len(notification.UserID)
	var flag []int = make([]int, l)
	//now we need to get fcm-token for given email, i.e., from email we need userID and using that we will get fcm-token
	for k, userId := range notification.UserID {
		//userId := base64.StdEncoding.EncodeToString([]byte(*email))

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
			//log.Println("FCM-token for given userID ", v["FCM-token"].(string))
			dataJson, err := json.Marshal(m)

			if err != nil {
				log.Printf("Unable to convert to JSON: %v", err)
			}
			ch <- dataJson
			go sendToCache(ch, &mut)

			time.Sleep(2 * time.Second)
			data, _ := cache.Get(string(dataJson))
			code := string(data)
			res = append(res, &model.Notification{
				Statuscode: code,
			})
			if code == "1" {
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

func sendingToFirestore(dataJson []byte, userId string) {
	var msg message
	err := json.Unmarshal(dataJson, &msg)
	if err != nil {
		log.Println(err)
	}

	msgId := ksuid.New()
	_, err = global.Client.Collection("notification").Doc(msgId.String()).Set(global.Ct, model.FirestoreData{
		Title:     msg.Notification.Title,
		Body:      msg.Notification.Body,
		CreatedAt: int(time.Now().Unix()),
		MessageID: msgId.String(),
		UserID:    userId,
		IsRead:    false,
	})
	if err != nil {
		log.Fatalf("Failed adding value to cloud firestore: %v", err)
	}
}

func sendToCache(ch chan []byte, mut *sync.Mutex) {
	mut.Lock()

	firebaseCh := make(chan []byte, 100)

	dataJson := <-ch
	data, err := cache.Get(string(dataJson))

	statusCode := ""
	_ = json.Unmarshal(data, &statusCode)

	var m sync.Mutex

	//checking if value has not been seen before, then sending it to firebase
	if err != nil || statusCode == "" {
		firebaseCh <- dataJson
		go sendToFirebase(firebaseCh, &m)
	}
	mut.Unlock()
}

func sendToFirebase(ch chan []byte, m *sync.Mutex) {
	m.Lock()

	dataJson := <-ch
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
	//log.Println(successCode.Results[0].MessageId)
	//log.Println("Key", string(dataJson))
	err = cache.Set(string(dataJson), []byte(strconv.Itoa(successCode.Success)))
	if err != nil {
		log.Printf(" Got error while setting the key %v", err)
	}

	//log.Println(successCode.Success)

	m.Unlock()

}

func GetClaimsFromContext(ctx context.Context) (map[string]interface{}, error) {
	token := ctx.Value("token").(string)
	claims, err := jwt.GetClaims(token)
	return claims, err
}
