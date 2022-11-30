package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	//"os"
	"sync"
	"time"

	"encoding/json"

	"github.com/allegro/bigcache/v3"
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
	MessageId string `json:"message_ids"`
}

type respBody struct {
	Multicast_id  int `json:"multicast_id"`
	Success       int `json:"success"`
	Failure       int `json:"failure"`
	Canonical_ids int `json:"canonical_id"`
	Results       []results
}

var cache *bigcache.BigCache

func SendNotification(ctx context.Context, notification model.NotificationInput) (*model.Notification, error) {
	global.Ct = ctx

	fcm_token := fmt.Sprintf("%s", ctx.Value("fcm-token"))
	//log.Println(fcm_token)

	cacheVar, err := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	if err != nil {
		log.Printf("Unable to create cache %v", err)
	}
	cache = cacheVar

	s := skeleton{
		Title: notification.Title,
		Body:  notification.Body,
	}

	m := message{
		Notification: s,
		To:           fcm_token,
		CreatedAt:    time.Now().Unix(),
	}
	dataJson, err := json.Marshal(m)

	if err != nil {
		log.Printf("Unable to convert to JSON: %v", err)
	}

	//sending data to cache function
	ch := make(chan []byte, 100)
	var mut sync.Mutex

	ch <- dataJson
	go sendToCache(ch, &mut)

	time.Sleep(2 * time.Second)

	data, _ := cache.Get(string(dataJson))
	statusCode := ""
	err = json.Unmarshal(data, &statusCode)
	if err != nil {
		log.Println("Error while converting to json ", err)
	}

	//log.Println(statusCode)
	return &model.Notification{
		Statuscode: statusCode,
	}, nil

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

	code := strconv.Itoa(successCode.Success)

	//marshalling it to send it to cache
	res, err := json.Marshal(code)
	if err != nil {
		//it means success is 0 i.e., unable to send request
		var temp int = 0
		tempBytes, _ := json.Marshal(temp)

		err = cache.Set(string(dataJson), tempBytes)
		if err != nil {
			log.Printf(" Got error while setting the key %v", err)
		}
	}

	_ = cache.Set(string(dataJson), res)
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
