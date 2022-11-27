package handlers

import (
	"context"
	"encoding/base64"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/zicops/zicops-notification-server/global"
	"github.com/zicops/zicops-notification-server/graph/model"
	"google.golang.org/api/iterator"
)

func AddToDatastore(m message, userId string) {

	//log.Println("Context in addToDatastore ", global.Ct)
	_, _, err := global.Client.Collection("notification").Add(global.Ct, model.FirestoreMessage{
		Title:     m.Notification.Title,
		Body:      m.Notification.Body,
		CreatedAt: int(m.CreatedAt),
		UserID:    userId,
	})
	if err != nil {
		log.Fatalf("Failed adding value to cloud firestore: %v", err)
	}
}

func GetAllNotifications(ctx context.Context, prevPageSnapShot string, pageSize int) (*model.PaginatedNotifications, error) {

	var firestoreResp []*model.FirestoreMessage
	claims, _ := GetClaimsFromContext(ctx)
	email_creator := claims["email"].(string)
	userId := base64.StdEncoding.EncodeToString([]byte(email_creator))
	startAfter := prevPageSnapShot
	var iter *firestore.DocumentIterator
	if startAfter == "" {
		iter = global.Client.Collection("notification").Where("UserID", "==", userId).OrderBy("CreatedAt", firestore.Desc).Limit(pageSize).Documents(ctx)

	} else {
		iter = global.Client.Collection("notification").Where("UserID", "==", userId).OrderBy("CreatedAt", firestore.Desc).StartAfter(startAfter).Limit(pageSize).Documents(ctx)
	}
	var resp []map[string]interface{}
	var lastDoc *firestore.DocumentSnapshot
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
			return nil, err
		}
		lastDoc = doc
		resp = append(resp, doc.Data())
	}
	prevSeenData := lastDoc.Ref.ID
	for _, v := range resp {
		createdAt, _ := v["CreatedAt"].(int64)
		tmp := &model.FirestoreMessage{
			Body:      v["Body"].(string),
			Title:     v["Title"].(string),
			CreatedAt: int(createdAt),
			UserID:    v["UserID"].(string),
		}
		//log.Println(tmp.Body, "      ", tmp.Title)
		firestoreResp = append(firestoreResp, tmp)
	}
	return &model.PaginatedNotifications{
		Messages:         firestoreResp,
		NextPageSnapShot: &prevSeenData,
	}, nil
}
