package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/zicops/zicops-notification-server/global"
	"github.com/zicops/zicops-notification-server/graph/model"
)

func AddMessagesMeet(ctx context.Context, message *model.Messages) (*bool, error) {
	_, err := GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if message == nil {
		return nil, err
	}
	id := uuid.New().String()

	_, err = global.Client.Collection("MeetMessages").Doc(id).Set(ctx, map[string]interface{}{
		"body":       message.Body,
		"user_id":    message.Body,
		"time":       message.Time,
		"meeting_id": message.MeetingID,
	})
	if err != nil {
		return nil, err
	}

	res := true
	return &res, nil
}
