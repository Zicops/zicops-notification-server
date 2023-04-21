package handlers

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/zicops/zicops-notification-server/global"
	"github.com/zicops/zicops-notification-server/graph/model"
)

func AddMessagesMeet(ctx context.Context, message *model.Messages) (*bool, error) {
	if message.Body == nil || message.ChatType == nil || message.MeetingID == nil || message.UserID == nil {
		return nil, fmt.Errorf("please mention all the parameters")
	}
	_, err := GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if message == nil {
		return nil, err
	}
	id := uuid.New().String()

	_, err = global.Client.Collection("MeetMessages").Doc(id).Set(ctx, map[string]interface{}{
		"parent_id":  message.ParentID,
		"body":       message.Body,
		"user_id":    message.UserID,
		"time":       message.Time,
		"meeting_id": message.MeetingID,
		"chat_type":  message.ChatType,
		"responses":  message.Responses,
	})
	if err != nil {
		return nil, err
	}

	if message.ParentID != nil {
		_, err := global.Client.Collection("MeetMessages").Doc(*message.ParentID).Update(ctx, []firestore.Update{
			{Path: "responses",
				Value: firestore.Increment(1)},
		})
		if err != nil {
			return nil, err
		}
	}

	res := true
	return &res, nil
}
