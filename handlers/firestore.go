package handlers

import (
	"context"
	"encoding/base64"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/zicops/zicops-notification-server/global"
	"github.com/zicops/zicops-notification-server/graph/model"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func AddToDatastore(ctx context.Context, m []*model.FirestoreDataInput) (string, error) {
	claims, _ := GetClaimsFromContext(ctx)
	email_creator := claims["email"].(string)
	userId := base64.StdEncoding.EncodeToString([]byte(email_creator))

	log.Println(email_creator)
	for _, message := range m {

		//if person has not yet seen the notification, i.e., notification as of now is just pushed to frontend
		if !message.IsRead {
			//we will add it to datastore
			_, _, err := global.Client.Collection("notification").Add(ctx, model.FirestoreData{
				Title:     message.Title,
				Body:      message.Body,
				CreatedAt: int(time.Now().Unix()),
				MessageID: message.MessageID,
				UserID:    userId,
				IsRead:    false,
			})
			if err != nil {
				log.Fatalf("Failed adding value to cloud firestore: %v", err)
			}

		} else if message.IsRead {
			//means person has clicked on the notification and we want to update the value
			//if value does not exist then give error
			_, err := global.Client.Collection("notification").Doc(message.MessageID).Get(ctx)
			if status.Code(err) == codes.NotFound {
				return "Value not found", err
			}

			//else update
			_, err = global.Client.Collection("notification").Doc(message.MessageID).Update(ctx, []firestore.Update{
				{
					Path:  "IsRead",
					Value: true,
				},
			})
			if err != nil {
				return "Unable to update the notification", err
			}
			return "Values updated successfully", nil
		}
	}

	return "Values added successfully", nil
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
			MessageID: v["MessageID"].(string),
		}
		//log.Println(tmp.Body, "      ", tmp.Title)
		firestoreResp = append(firestoreResp, tmp)
	}
	return &model.PaginatedNotifications{
		Messages:         firestoreResp,
		NextPageSnapShot: &prevSeenData,
	}, nil
}
