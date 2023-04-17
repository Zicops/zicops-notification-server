package handlers

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/zicops/zicops-notification-server/global"
	"github.com/zicops/zicops-notification-server/graph/model"
)

func AddPoll(ctx context.Context, input *model.PollsInput) (*model.Polls, error) {
	_, err := GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if input == nil || input.MeetingID == nil || input.Status == nil || input.CourseID == nil || input.Question == nil {
		return nil, fmt.Errorf("please mention all the parameters")
	}
	id := uuid.New().String()

	var options []string
	var pollIds []string
	for k, vv := range input.Options {
		if vv == nil {
			continue
		}
		v := *vv
		options = append(options, v)
		pollinp := model.PollResponseInput{
			PollID:   &id,
			Response: input.Options[k],
		}
		poll_id, err := addPollResponse(ctx, pollinp)
		if err != nil {
			return nil, err
		}
		pollIds = append(pollIds, poll_id)
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
	res := model.Polls{
		ID:        &id,
		MeetingID: input.MeetingID,
		CourseID:  input.CourseID,
		TopicID:   input.TopicID,
		Question:  input.Question,
		Options:   input.Options,
		Status:    input.Status,
	}
	return &res, nil
}

func addPollResponse(ctx context.Context, input model.PollResponseInput) (string, error) {
	id := uuid.New().String()
	_, err := global.Client.Collection("polls_response").Doc(id).Set(ctx, map[string]interface{}{
		"poll_id":  *input.PollID,
		"response": *input.Response,
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func UpdatePoll(ctx context.Context, input *model.PollsInput) (*model.Polls, error) {
	_, err := GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	updates := []firestore.Update{}

	if input.Question != nil {
		updates = append(updates, firestore.Update{
			Path:  "question",
			Value: *input.Question,
		})
	}
	if input.Status != nil {
		updates = append(updates, firestore.Update{
			Path:  "status",
			Value: *input.Status,
		})
	}
	if input.Options != nil {
		var options []string
		for _, vv := range input.Options {
			if vv == nil {
				continue
			}
			v := *vv
			options = append(options, v)
		}

		updates = append(updates, firestore.Update{
			Path:  "options",
			Value: options,
		})
	}
	_, err = global.Client.Collection("polls").Doc(*input.ID).Update(ctx, updates)
	if err != nil {
		return nil, err
	}

	res := model.Polls{
		ID:        input.ID,
		MeetingID: input.MeetingID,
		CourseID:  input.CourseID,
		TopicID:   input.TopicID,
		Question:  input.Question,
		Options:   input.Options,
		Status:    input.Status,
	}
	return &res, nil
}

/*
classroom flag me

publishedQuiz: []
endedQuiz:[]

ye db me rakho

main tumko single quizId aur type: "publish" bhejunga tum usko publish me append karo
aur agar type: "end" bhejunga tho ended me append karo aur publish me se remove karo
*/
