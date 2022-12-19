package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/zicops/zicops-notification-server/global"
	"google.golang.org/api/iterator"
)

type Payload struct {
	RegistrationIds []string `json:"registration_ids"`
	DryRun          bool     `json:"dry_run"`
}

type Response struct {
	Multicast_id  int     `json:"multicast_id"`
	Success       int     `json:"success"`
	Failure       int     `json:"failure"`
	Canonical_ids int     `json:"canonical_ids"`
	Results       []array `json:"results"`
}

type array struct {
	Message_id string `json:"message_id"`
}

func Auth_tokens(ctx context.Context) (string, error) {

	global.Ct = ctx
	iter := global.Client.Collection("tokens").Documents(global.Ct)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}
		verify(doc.Data(), doc.Ref.ID)
	}

	return "Delteted redundant tokens", nil
}

func verify(token map[string]interface{}, id string) {
	//here we will go through each data items, check their tokens and delete which are non existing
	t := token["FCM-token"].(string)
	data := Payload{
		RegistrationIds: []string{t},
		DryRun:          true,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Println("Got error while marshalling the payload", err)
	}

	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", body)
	if err != nil {
		log.Println("Got error while sending the request ", err)
	}

	req.Header.Set("Authorization", "key=AAAAU56xXYc:APA91bHtHX1hjkj8B4u0tSTuuTgURF6PvlqKEzgn3Qv7JR14mwra7rrCCg3bRRJZHxYyK8DHntk4Tc9CsXkqj44vuxFeD1WgRy1nifgbYgi60IAmfApLKK6Rd92Puuj3NPtUNGvdNHvr")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Got error while receiving the respone")
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to receive response to JSON")
	}
	var r Response
	err = json.Unmarshal(b, &r)
	if err != nil {
		log.Println("Got error while unmarshalling data ", err)
	}
	if r.Success == 0 {
		//delete token
		_, err = global.Client.Collection("tokens").Doc(id).Delete(global.Ct)
		if err != nil {
			log.Println("Got error while deleting the data ", err)
		}
	}

	//log.Println(string(b))
	defer resp.Body.Close()

	//log.Println(token["FCM-token"].(string))
}
