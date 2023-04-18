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

	_, err = global.Client.Collection("ClassroomFlags").Doc(*input.ID).Set(ctx, map[string]interface{}{
		"IsClassroomStarted":        input.IsClassroomStarted,
		"IsParticipantsPresent":     input.IsParticipantsPresent,
		"IsAdDisplayed":             input.IsAdDisplayed,
		"IsBreak":                   input.IsBreak,
		"IsModeratorJoined":         input.IsModeratorJoined,
		"IsTrainerJoined":           input.IsTrainerJoined,
		"AdVideoURL":                input.AdVideoURL,
		"is_microphone_enabled":     input.IsMicrophoneEnabled,
		"is_video_sharing_enabled":  input.IsVideoSharingEnabled,
		"is_screen_sharing_enabled": input.IsScreenSharingEnabled,
		"is_chat_enabled":           input.IsChatEnabled,
		"is_qa_enabled":             input.IsQaEnabled,
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
