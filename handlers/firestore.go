package handlers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/zicops/zicops-notification-server/global"
	//"github.com/zicops/zicops-notification-server/graph/model"
	"google.golang.org/api/iterator"
)

func AddToDatastore(m message) {

	//log.Println("Context in addToDatastore ", global.Ct)
	_, _, err := global.Client.Collection("notification").Add(global.Ct, map[string]interface{}{
		"title": m.Notification.Title,
		"body":  m.Notification.Body,
	})
	if err != nil {
		log.Fatalf("Failed adding value to cloud firestore: %v", err)
	}
}

func GetAllNotifications(ctx context.Context) (string, error) {

	//firestoreResp := &model.FirestoreMessage{}
	iter := global.Client.Collection("notification").Documents(ctx)
	var resp []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
			return "", err
		}
		resp = append(resp, doc.Data())
	}
	res, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Got error while converting response to JSON")
	}

	return string(res), nil
}
