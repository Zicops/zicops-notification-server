package handlers

import (
	"context"
	"encoding/base64"
	"log"

	"github.com/zicops/zicops-notification-server/global"
	"github.com/zicops/zicops-notification-server/graph/model"
	"google.golang.org/api/iterator"
)

func AddToDatastore(m message, userId string) {

	//log.Println("Context in addToDatastore ", global.Ct)
	_, _, err := global.Client.Collection("notification/"+userId).Add(global.Ct, map[string]interface{}{
		"title": m.Notification.Title,
		"body":  m.Notification.Body,
	})
	if err != nil {
		log.Fatalf("Failed adding value to cloud firestore: %v", err)
	}
}

func GetAllNotifications(ctx context.Context) ([]*model.FirestoreMessage, error) {

	var firestoreResp []*model.FirestoreMessage
	claims, _ := GetClaimsFromContext(ctx)
	email_creator := claims["email"].(string)
	userId := base64.StdEncoding.EncodeToString([]byte(email_creator))
	iter := global.Client.Collection("notification/" + userId).Documents(ctx)
	var resp []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
			return firestoreResp, err
		}
		resp = append(resp, doc.Data())
	}

	for _, v := range resp {
		tmp := &model.FirestoreMessage{
			Body:  v["body"].(string),
			Title: v["title"].(string),
		}
		//log.Println(tmp.Body, "      ", tmp.Title)
		firestoreResp = append(firestoreResp, tmp)
	}

	return firestoreResp, nil
}
