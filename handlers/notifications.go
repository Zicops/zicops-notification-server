package handlers

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"

	//"os"
	"sync"
	"time"

	"encoding/base64"
	"encoding/json"

	"github.com/allegro/bigcache/v3"
	"github.com/zicops/zicops-notification-server/global"
	"github.com/zicops/zicops-notification-server/graph/model"
)

type message struct {
	Notification skeleton `json:"notification"`
	To           string   `json:"to"`
}

type skeleton struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

var cache *bigcache.BigCache

func SendNotification(ctx context.Context, notification model.NotificationInput) (*model.Notification, error) {
	global.Ct = ctx
	token := fmt.Sprintf("%s", ctx.Value("token"))
	//log.Println(token)
	cacheVar, err := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	if err != nil {
		log.Printf("Unable to create cache %v", err)
	}
	claims, _ := GetClaimsFromContext(ctx)
	email_creator := claims["email"].(string)
	userId := base64.StdEncoding.EncodeToString([]byte(email_creator))
	cache = cacheVar

	s := skeleton{
		Title: notification.Title,
		Body:  notification.Body,
	}

	m := message{
		Notification: s,
		To:           token,
	}
	AddToDatastore(m, userId)

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
	_ = json.Unmarshal(data, &statusCode)

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
	log.Printf("Status code of our request is %v", resp.Status)

	statusCode, err := json.Marshal(resp.Status)
	if err != nil {
		log.Printf("Got error while converting data to json %v", err)
	}

	//setting the response in cache
	_ = cache.Set(string(dataJson), statusCode)
	if err != nil {
		log.Printf(" Got error while setting the key %v", err)
	}

	m.Unlock()

}

func GetClaimsFromContext(ctx context.Context) (map[string]interface{}, error) {
	customClaims := ctx.Value("zclaims").(map[string]interface{})
	if customClaims == nil {
		return make(map[string]interface{}), fmt.Errorf("custom claims not found. Unauthorized user")
	}
	return customClaims, nil
}
