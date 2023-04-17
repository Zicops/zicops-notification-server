package handlers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/zicops/zicops-notification-server/global"
	"github.com/zicops/zicops-notification-server/graph/model"
)

func AddPoll(ctx context.Context, input *model.PollsInput) (*model.Polls, error) {
	_, err := GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if input == nil || input.MeetingID == nil || input.Status == nil {
		return nil, fmt.Errorf("please mention meeting id as well as status")
	}
	id := uuid.New().String()

	var options []string
	for k, vv := range input.Options {
		v := *vv
		options = append(options, v)
		pollinp := model.PollResponseInput{
			PollID:   &id,
			Response: input.Options[k],
		}
		err := addPollResponse(ctx, pollinp)
		if err != nil {
			return nil, err
		}
	}
	_, err = global.Client.Collection("polls").Doc(id).Set(ctx, map[string]interface{}{
		"meeting_id": *input.MeetingID,
		"course_id":  *input.CourseID,
		"topic_id":   *input.TopicID,
		"question":   *input.Question,
		"options":    options,
		"status":     *input.Status,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func addPollResponse(ctx context.Context, input model.PollResponseInput) error {
	id := uuid.New().String()
	_, err := global.Client.Collection("polls_response").Doc(id).Set(ctx, map[string]interface{}{
		"poll_id":  *input.PollID,
		"response": *input.Response,
	})
	if err != nil {
		return err
	}
	return nil
}
