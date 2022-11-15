package handlers

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"encoding/json"

	"github.com/allegro/bigcache/v3"
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

	cacheVar, err := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	if err != nil {
		log.Printf("Unable to create cache %v", err)
	}

	cache = cacheVar

	ch := make(chan []byte, 100)
	var mut sync.Mutex

	//we need tokens from all the users to whom we have to send notifications

	s := skeleton{
		Title: notification.Title,
		Body:  notification.Body,
	}

	m := message{
		Notification: s,
		To:           notification.Token,
	}

	dataJson, err := json.Marshal(m)

	if err != nil {
		log.Printf("Unable to convert to JSON: %v", err)
	}

	ch <- dataJson
	go sendToCache(ch, &mut)

	//list of notifications to be sent to people

	data, _ := cache.Get(string(dataJson))

	statusCode := 0
	_ = json.Unmarshal(data, &statusCode)

	return &model.Notification{
		Statuscode: strconv.Itoa(statusCode),
	}, nil

}

func sendToCache(ch chan []byte, mut *sync.Mutex) {
	mut.Lock()

	firebaseCh := make(chan []byte)
	var m sync.Mutex

	dataJson := <-ch

	data, err := cache.Get(string(dataJson))

	statusCode := 0
	_ = json.Unmarshal(data, &statusCode)

	if err != nil || statusCode == 404 {
		//log
		firebaseCh <- dataJson
		go sendToFirebase(firebaseCh, &m)
	}
	mut.Unlock()
}

func sendToFirebase(ch chan []byte, m *sync.Mutex) {
	m.Lock()

	dataJson := <-ch
	body := bytes.NewReader(dataJson)
	// for sending to all users,
	// subscribe them all to a topic "all"
	// and then send notifications to this topic
	// FirebaseMessaging.getInstance().subscribeToTopic("all");
	req, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", body)

	if err != nil {
		log.Printf("Error while sending request %v", err)
	}
	req.Header.Set("Authorization", "key=AAAAU56xXYc:APA91bHtHX1hjkj8B4u0tSTuuTgURF6PvlqKEzgn3Qv7JR14mwra7rrCCg3bRRJZHxYyK8DHntk4Tc9CsXkqj44vuxFeD1WgRy1nifgbYgi60IAmfApLKK6Rd92Puuj3NPtUNGvdNHvr")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error while getting response- %v", err)
	}

	statusCode, err := json.Marshal(resp.StatusCode)
	log.Println(statusCode)
	if err != nil {
		log.Printf("Got error while converting data to json %v", err)
	}
	cache.Set(string(dataJson), statusCode)
	m.Unlock()

}
