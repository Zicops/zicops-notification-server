package handlers

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/zicops/zicops-notification-server/global"
	"github.com/zicops/zicops-notification-server/graph/model"
)

func AddClassroomFlags(ctx context.Context, input *model.ClassRoomFlagsInput) (*model.ClassRoomFlags, error) {
	_, err := GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if input == nil {
		return nil, nil
	}
	if input.ID == nil {
		return nil, err
	}

	_, err = global.Client.Collection("ClassroomFlags").Doc(*input.ID).Set(ctx, map[string]interface{}{
		"is_classroom_started":      input.IsClassroomStarted,
		"is_participants_present":   input.IsParticipantsPresent,
		"is_ad_displayed":           input.IsAdDisplayed,
		"is_break":                  input.IsBreak,
		"is_moderator_joined":       input.IsModeratorJoined,
		"is_trainer_joined":         input.IsTrainerJoined,
		"ad_video_url":              input.AdVideoURL,
		"is_microphone_enabled":     input.IsMicrophoneEnabled,
		"is_video_sharing_enabled":  input.IsVideoSharingEnabled,
		"is_screen_sharing_enabled": input.IsScreenSharingEnabled,
		"is_chat_enabled":           input.IsChatEnabled,
		"is_qa_enabled":             input.IsQaEnabled,
		"is_classroom_ended":        input.IsClassroomEnded,
	})
	if err != nil {
		return nil, err
	}

	res := model.ClassRoomFlags{
		ID:                     input.ID,
		IsClassroomStarted:     input.IsClassroomStarted,
		IsParticipantsPresent:  input.IsParticipantsPresent,
		IsAdDisplayed:          input.IsAdDisplayed,
		IsBreak:                input.IsBreak,
		IsModeratorJoined:      input.IsModeratorJoined,
		IsTrainerJoined:        input.IsTrainerJoined,
		AdVideoURL:             input.AdVideoURL,
		IsMicrophoneEnabled:    input.IsMicrophoneEnabled,
		IsVideoSharingEnabled:  input.IsVideoSharingEnabled,
		IsScreenSharingEnabled: input.IsScreenSharingEnabled,
		IsChatEnabled:          input.IsChatEnabled,
		IsQaEnabled:            input.IsQaEnabled,
		IsClassroomEnded:       input.IsClassroomEnded,
	}
	return &res, nil
}

func AddQuizToClassroomFlags(ctx context.Context, input *model.PublishedQuiz) (*bool, error) {
	_, err := GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if input == nil {
		return nil, nil
	}

	_, err = global.Client.Collection("ClassroomFlags").Doc(*input.ID).Update(ctx, []firestore.Update{
		{
			Path:  "quiz",
			Value: firestore.ArrayUnion(*input.QuizID),
		},
	})
	if err != nil {
		return nil, err
	}
	res := true
	return &res, nil
}
