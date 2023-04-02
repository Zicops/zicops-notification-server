package handlers

import (
	"context"

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
		"IsClassroomStarted":    input.IsClassroomStarted,
		"IsParticipantsPresent": input.IsParticipantsPresent,
		"IsAdDisplayed":         input.IsAdDisplayed,
		"IsBreak":               input.IsBreak,
		"IsModeratorJoined":     input.IsModeratorJoined,
		"IsTrainerJoined":       input.IsTrainerJoined,
		"AdVideoURL":            input.AdVideoURL,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
