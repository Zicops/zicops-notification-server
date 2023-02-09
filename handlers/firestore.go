package handlers

import (
	"context"
	"encoding/base64"
	"log"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/zicops/zicops-notification-server/global"
	"github.com/zicops/zicops-notification-server/graph/model"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func AddToDatastore(ctx context.Context, m []*model.FirestoreDataInput) (string, error) {
	claims, err := GetClaimsFromContext(ctx)
	if err != nil {
		log.Printf("Error getting claims from headers: %v", err)
	}
	email_creator := claims["email"].(string)
	userId := base64.StdEncoding.EncodeToString([]byte(email_creator))

	//log.Println(email_creator)
	for _, message := range m {

		//if person has not yet seen the notification, i.e., notification as of now is just pushed to frontend
		if !message.IsRead {
			//we will add it to datastore
			_, err := global.Client.Collection("notification").Doc(message.MessageID).Set(global.Ct, model.FirestoreData{
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

type TokenSave struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}

func AddToDatastoreFCMToken(ctx context.Context, m TokenSave) (string, error) {
	_, _, err := global.Client.Collection("tokens").Add(global.Ct, m)
	if err != nil {
		log.Fatalf("Failed adding value to cloud firestore: %v", err)
	}
	return "Values added successfully", nil
}

func GetAllNotifications(ctx context.Context, prevPageSnapShot string, pageSize int, isRead *bool) (*model.PaginatedNotifications, error) {

	var firestoreResp []*model.FirestoreMessage
	claims, _ := GetClaimsFromContext(ctx)
	email_creator := claims["email"].(string)
	lspId := claims["lsp_id"].(string)
	userId := base64.StdEncoding.EncodeToString([]byte(email_creator))
	startAfter := prevPageSnapShot
	var iter *firestore.DocumentIterator
	if isRead != nil {
		if startAfter == "" {
			iter = global.Client.Collection("notification").Where("UserID", "==", userId).Where("IsRead", "==", isRead).Where("LspID", "==", lspId).OrderBy("CreatedAt", firestore.Desc).Limit(pageSize).Documents(ctx)

		} else {
			iter = global.Client.Collection("notification").Where("UserID", "==", userId).Where("IsRead", "==", isRead).Where("LspID", "==", lspId).OrderBy("CreatedAt", firestore.Desc).StartAfter(startAfter).Limit(pageSize).Documents(ctx)
		}
	} else {
		if startAfter == "" {
			iter = global.Client.Collection("notification").Where("UserID", "==", userId).Where("LspID", "==", lspId).OrderBy("CreatedAt", firestore.Desc).Limit(pageSize).Documents(ctx)

		} else {
			iter = global.Client.Collection("notification").Where("UserID", "==", userId).Where("LspID", "==", lspId).OrderBy("CreatedAt", firestore.Desc).StartAfter(startAfter).Limit(pageSize).Documents(ctx)
		}
	}

	var resp []map[string]interface{}
	var lastDoc *firestore.DocumentSnapshot
	for {
		doc, err := iter.Next()
		//see if iterator is done
		if err == iterator.Done {
			break
		}

		//see if the error is no more items in iterator
		if err != nil && err.Error() == "no more items in iterator" {
			break
			//return nil, nil
		}

		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
			return nil, err
		}
		lastDoc = doc
		resp = append(resp, doc.Data())

	}
	if resp == nil {
		return nil, nil
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
			IsRead:    v["IsRead"].(bool),
			Link:      v["Link"].(string),
			LspID:     v["LspID"].(string),
		}
		//log.Println(tmp.Body, "      ", tmp.Title)
		firestoreResp = append(firestoreResp, tmp)
	}
	return &model.PaginatedNotifications{
		Messages:         firestoreResp,
		NextPageSnapShot: &prevSeenData,
	}, nil
}

func GetAllPaginatedNotifications(ctx context.Context, prevPageSnapShot string, pageSize int, isRead *bool) (*model.PaginatedNotifications, error) {

	var firestoreResp []*model.FirestoreMessage
	claims, _ := GetClaimsFromContext(ctx)
	email_creator := claims["email"].(string)
	lspId := claims["lsp_id"].(string)
	userId := base64.StdEncoding.EncodeToString([]byte(email_creator))
	startAfter := prevPageSnapShot
	var iter *firestore.DocumentIterator
	if isRead != nil {
		if startAfter == "" {
			iter = global.Client.Collection("notification").Where("UserID", "==", userId).Where("IsRead", "==", isRead).Where("LspID", "==", lspId).OrderBy("CreatedAt", firestore.Desc).Limit(pageSize).Documents(ctx)

		} else {
			iter = global.Client.Collection("notification").Where("UserID", "==", userId).Where("IsRead", "==", isRead).Where("LspID", "==", lspId).OrderBy("CreatedAt", firestore.Desc).Limit(pageSize).StartAfter(startAfter).Documents(ctx)
		}
	} else {
		if startAfter == "" {
			iter = global.Client.Collection("notification").Where("UserID", "==", userId).Where("LspID", "==", lspId).OrderBy("CreatedAt", firestore.Desc).Limit(pageSize).Documents(ctx)

		} else {
			iter = global.Client.Collection("notification").Where("UserID", "==", userId).Where("LspID", "==", lspId).OrderBy("CreatedAt", firestore.Desc).Limit(pageSize).StartAfter(startAfter).Documents(ctx)
		}
	}
	tmp, err := iter.GetAll()
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	var resp []map[string]interface{}
	lstDoc := tmp[len(tmp)-1]
	for _, vv := range tmp {
		v := vv
		data := v.Data()
		resp = append(resp, data)
	}
	rs, err := lstDoc.DataAt("CreatedAt")
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	for _, v := range resp {
		createdAt, _ := v["CreatedAt"].(int64)
		tmp := &model.FirestoreMessage{
			Body:      v["Body"].(string),
			Title:     v["Title"].(string),
			CreatedAt: int(createdAt),
			UserID:    v["UserID"].(string),
			MessageID: v["MessageID"].(string),
			IsRead:    v["IsRead"].(bool),
			Link:      v["Link"].(string),
			LspID:     v["LspID"].(string),
		}
		//log.Println(tmp.Body, "      ", tmp.Title)
		firestoreResp = append(firestoreResp, tmp)
	}
	//r := lstDoc.Ref.ID
	ca := rs.(int64)
	res := strconv.Itoa(int(ca))
	return &model.PaginatedNotifications{
		Messages:         firestoreResp,
		NextPageSnapShot: &res,
	}, nil

}
