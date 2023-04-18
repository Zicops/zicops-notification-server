package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/zicops/zicops-notification-server/global"
	"github.com/zicops/zicops-notification-server/graph/model"
	"google.golang.org/api/iterator"
)

func AddPoll(ctx context.Context, input *model.PollsInput) (*model.Polls, error) {
	claims, err := GetClaimsFromContext(ctx)
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
	createdBy := claims["email"].(string)
	createdAt := time.Now().String()
	_, err = global.Client.Collection("polls").Doc(id).Set(ctx, map[string]interface{}{
		"poll_name":  *input.PollName,
		"meeting_id": *input.MeetingID,
		"course_id":  *input.CourseID,
		"topic_id":   *input.TopicID,
		"question":   *input.Question,
		"options":    options,
		"created_at": createdAt,
		"created_by": createdBy,
		"updated_at": createdAt,
		"updated_by": createdBy,
		"status":     *input.Status,
	})
	if err != nil {
		return nil, err
	}
	var tmp []*string
	for _, vv := range pollIds {
		v := vv
		tmp = append(tmp, &v)
	}
	res := model.Polls{
		ID:            &id,
		MeetingID:     input.MeetingID,
		CourseID:      input.CourseID,
		TopicID:       input.TopicID,
		Question:      input.Question,
		Options:       input.Options,
		PollOptionIds: tmp,
		Status:        input.Status,
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
		PollName:  input.PollName,
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

func UpdatePollOptions(ctx context.Context, input *model.PollResponseInput) (*model.PollResponse, error) {
	_, err := GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if input.UserID != nil {
		if input.PollID == nil || input.Option == nil {
			return nil, fmt.Errorf("please enter poll id, option, and userId")
		}
		optionId, err := getIdOfPollOption(ctx, *input.PollID, *input.Option)
		if err != nil {
			return nil, err
		}
		iter := global.Client.Collection("polls_response").Where("user_ids", "array-contains", *input.UserID).Where("poll_id", "==", *input.PollID).Documents(ctx)
		for {
			doc, err := iter.Next()
			//see if iterator is done
			if err == iterator.Done {
				break
			}

			//see if the error is no more items in iterator
			if err != nil && err.Error() == "no more items in iterator" {
				break
			}

			if err != nil {
				log.Fatalf("Failed to iterate: %v", err)
				return nil, err
			}
			id := doc.Ref.ID
			_, err = global.Client.Collection("polls_response").Doc(id).Update(ctx, []firestore.Update{
				{
					Path:  "user_ids",
					Value: firestore.ArrayRemove(input.UserID),
				},
			})
			if err != nil {
				return nil, err
			}
		}

		_, err = global.Client.Collection("polls_response").Doc(optionId).Update(ctx, []firestore.Update{
			{
				Path:  "user_ids",
				Value: firestore.ArrayUnion(*input.UserID),
			},
		})
		if err != nil {
			return nil, err
		}

	}

	res := model.PollResponse{
		ID:       input.ID,
		PollID:   input.PollID,
		Response: input.Response,
		UserID:   input.UserID,
	}
	return &res, nil
}

func getIdOfPollOption(ctx context.Context, pollId string, option string) (string, error) {
	iter := global.Client.Collection("polls_response").Where("poll_id", "==", pollId).Where("response", "==", option).Documents(ctx)
	var res string
	for {
		doc, err := iter.Next()
		//see if iterator is done
		if err == iterator.Done {
			break
		}

		//see if the error is no more items in iterator
		if err != nil && err.Error() == "no more items in iterator" {
			break
		}

		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
			return "", err
		}
		res = doc.Ref.ID
	}
	return res, nil
}

func GetPollResults(ctx context.Context, pollID *string) (*model.PollResults, error) {
	if pollID == nil {
		return nil, fmt.Errorf("please enter poll id")
	}
	_, err := GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	dataS, err := global.Client.Collection("polls").Doc(*pollID).Get(ctx)
	if err != nil {
		return nil, err
	}

	data := dataS.Data()
	question := data["question"].(string)

	iter := global.Client.Collection("poll_response").Where("poll_id", "==", pollID).Documents(ctx)
	var pollData []map[string]interface{}
	var ids []string
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

		pollData = append(pollData, doc.Data())
		ids = append(ids, doc.Ref.ID)
	}
	if pollData == nil {
		return nil, nil
	}
	if len(pollData) == 0 {
		return nil, nil
	}

	var pollResponse []*model.PollResponse
	for k, vv := range pollData {
		v := vv
		response := v["response"].(string)
		users := v["user_ids"].([]interface{})
		for _, x := range users {
			userId := x.(string)
			tmp := model.PollResponse{
				ID:       &ids[k],
				PollID:   pollID,
				Response: &response,
				UserID:   &userId,
			}
			pollResponse = append(pollResponse, &tmp)
		}

	}
	res := model.PollResults{
		PollID:        pollID,
		Question:      &question,
		PollResponses: pollResponse,
	}
	return &res, nil
}
