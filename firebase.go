package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/allegro/bigcache/v3"
	//"github.com/gin-gonic/gin"
)

// this is the message that we will send to firebase server
type message struct {
	Notification skeleton `json:"notification"`
	To           string   `json:"to"`
}

// this is the basic skeleton of the message
type skeleton struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

/*

example of data that we will send
{
  "notification": {
    "title": "Portugal vs. Denmark",
    "body": "5 to 1",
    "icon": "firebase-logo.png",
    "click_action": "http://localhost:8081"
  },
  "to": "YOUR-IID-TOKEN"
}

*/

var cacheVar *bigcache.BigCache

func main() {

	ch := make(chan message, 2)

	for i := 0; i < 10; i++ {
		b := "This side notification from firebase " + strconv.Itoa(i)
		m := skeleton{
			Title: "Hello everyone message ",
			Body:  b,
		}

		tokenFromWebsite := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImQ3YjE5MTI0MGZjZmYzMDdkYzQ3NTg1OWEyYmUzNzgzZGMxYWY4OWYiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoiQW5zaCBKb3NoaSIsImlzcyI6Imh0dHBzOi8vc2VjdXJldG9rZW4uZ29vZ2xlLmNvbS96aWNvcHMtb25lIiwiYXVkIjoiemljb3BzLW9uZSIsImF1dGhfdGltZSI6MTY2NzkzMjk3MywidXNlcl9pZCI6IjRpZFowRkRERjZhanNReEpJOHJTcmtZMjRaeTIiLCJzdWIiOiI0aWRaMEZEREY2YWpzUXhKSThyU3JrWTI0WnkyIiwiaWF0IjoxNjY4MDE0ODIxLCJleHAiOjE2NjgwMTg0MjEsImVtYWlsIjoiYW5zaGpvc2hpMDYwN0BnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwicGhvbmVfbnVtYmVyIjoiKzkxOTg3NzMxNjU5NiIsImZpcmViYXNlIjp7ImlkZW50aXRpZXMiOnsicGhvbmUiOlsiKzkxOTg3NzMxNjU5NiJdLCJlbWFpbCI6WyJhbnNoam9zaGkwNjA3QGdtYWlsLmNvbSJdfSwic2lnbl9pbl9wcm92aWRlciI6InBhc3N3b3JkIn19.BnVGRQmQ7qaqz0vicpRj3rw0kQFNmwS4qxClvRSn4a0rYfNuQWm4-WSLbdqssbKZelTtiGZJrk9PF_e7ow9gc8KmbYt-62Uzd-Pd1l8_xkB2asmkbNw6qMgs4Hl1mL-Rsu_9DClsNJ7mg0vxZmIv1TVsBZhDg9mV5woiLxt9ouI14yLKUEkTNyUpwUuMD6YoeDgls9MF4cai2sXoFAIwU_PrLaaFjY0YYuQAF3Wo7mXQpnjO9GuDQsYS1FoQAF3K2qTlsWSvdYRUPGU1MwDHCfMR-_34Nil-PQvpOjBfy4tUdOW6BFguo0tBNqTzW-RjwWw5mEidQAdaRsbZAx9D6w"

		dat := message{
			Notification: m,
			To:           tokenFromWebsite,
		}

		//send this data to cache, from there if no hit then to firebase
		cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))

		if err != nil {
			log.Println("Unable to create the cache")
		}

		//we need a global variable, which will store this cache
		cacheVar = cache

		//we will be sending the message to other function via channels from here
		ch <- dat
		go sendToCache(ch)

		time.Sleep(2 * time.Second)
	}

	//dummy data for confirmation
	bDummy1 := "This side ansh from firebase " + strconv.Itoa(5)
	mDummy1 := skeleton{
		Title: "Hello everyone message ",
		Body:  bDummy1,
	}

	tokenFromWebsiteDummy1 := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImQ3YjE5MTI0MGZjZmYzMDdkYzQ3NTg1OWEyYmUzNzgzZGMxYWY4OWYiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoiQW5zaCBKb3NoaSIsImlzcyI6Imh0dHBzOi8vc2VjdXJldG9rZW4uZ29vZ2xlLmNvbS96aWNvcHMtb25lIiwiYXVkIjoiemljb3BzLW9uZSIsImF1dGhfdGltZSI6MTY2NzkzMjk3MywidXNlcl9pZCI6IjRpZFowRkRERjZhanNReEpJOHJTcmtZMjRaeTIiLCJzdWIiOiI0aWRaMEZEREY2YWpzUXhKSThyU3JrWTI0WnkyIiwiaWF0IjoxNjY4MDE0ODIxLCJleHAiOjE2NjgwMTg0MjEsImVtYWlsIjoiYW5zaGpvc2hpMDYwN0BnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwicGhvbmVfbnVtYmVyIjoiKzkxOTg3NzMxNjU5NiIsImZpcmViYXNlIjp7ImlkZW50aXRpZXMiOnsicGhvbmUiOlsiKzkxOTg3NzMxNjU5NiJdLCJlbWFpbCI6WyJhbnNoam9zaGkwNjA3QGdtYWlsLmNvbSJdfSwic2lnbl9pbl9wcm92aWRlciI6InBhc3N3b3JkIn19.BnVGRQmQ7qaqz0vicpRj3rw0kQFNmwS4qxClvRSn4a0rYfNuQWm4-WSLbdqssbKZelTtiGZJrk9PF_e7ow9gc8KmbYt-62Uzd-Pd1l8_xkB2asmkbNw6qMgs4Hl1mL-Rsu_9DClsNJ7mg0vxZmIv1TVsBZhDg9mV5woiLxt9ouI14yLKUEkTNyUpwUuMD6YoeDgls9MF4cai2sXoFAIwU_PrLaaFjY0YYuQAF3Wo7mXQpnjO9GuDQsYS1FoQAF3K2qTlsWSvdYRUPGU1MwDHCfMR-_34Nil-PQvpOjBfy4tUdOW6BFguo0tBNqTzW-RjwWw5mEidQAdaRsbZAx9D6w"

	datDummy1 := message{
		Notification: mDummy1,
		To:           tokenFromWebsiteDummy1,
	}
	ch <- datDummy1
	go sendToCache(ch)
	time.Sleep(2 * time.Second)
	//even this data is not getting stored in the cache, because response from firebase server is different each time

}

func sendToCache(ch chan message) {
	queries := <-ch

	dataJson, err := json.Marshal(queries)
	if err != nil {
		log.Println("Error while converting to JSON")
	}

	//we will use this channel to call firebase
	firebaseCh := make(chan []byte, 2)

	entry, err := cacheVar.Get(string(dataJson))
	if err != nil {
		//if not found, send to firebase
		firebaseCh <- dataJson
		go sendToFirebase(firebaseCh)

	} else {
		log.Printf("Response  - \n%v", entry)
	}
}

func sendToFirebase(ch chan []byte) {

	dataJson := <-ch
	body := bytes.NewReader(dataJson)

	req, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", body)

	if err != nil {
		log.Printf("Error while sending request-  %v", err)
	}

	req.Header.Set("Authorization", "key=AAAAU56xXYc:APA91bHtHX1hjkj8B4u0tSTuuTgURF6PvlqKEzgn3Qv7JR14mwra7rrCCg3bRRJZHxYyK8DHntk4Tc9CsXkqj44vuxFeD1WgRy1nifgbYgi60IAmfApLKK6Rd92Puuj3NPtUNGvdNHvr")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Printf("Error while getting response-  %v", err)
	}

	//if we can convert it to bytes array, then convert and enter to cache result, else display error
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Gor error while converting io.readCloser to bytes array")
	} else {
		cacheVar.Set(string(dataJson), respBody)
		log.Println(resp.StatusCode)
		log.Println(resp.Body)
	}
	defer resp.Body.Close()

	//return string(dataJson);

}
