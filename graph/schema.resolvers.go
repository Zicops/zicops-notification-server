package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"

	"github.com/zicops/zicops-notification-server/graph/generated"
	"github.com/zicops/zicops-notification-server/graph/model"
	"github.com/zicops/zicops-notification-server/handlers"
)

// SendNotificationWithLink is the resolver for the sendNotificationWithLink field.
func (r *mutationResolver) SendNotificationWithLink(ctx context.Context, notification model.NotificationInput, link string) ([]*model.Notification, error) {
	resp, err := handlers.SendNotificationWithLink(ctx, notification, link)
	if err != nil {
		log.Printf("Error sending notification %v", err)
		return nil, err
	}
	return resp, err
}

// AddToFirestore is the resolver for the addToFirestore field.
func (r *mutationResolver) AddToFirestore(ctx context.Context, message []*model.FirestoreDataInput) (string, error) {
	resp, err := handlers.AddToDatastore(ctx, message)

	if err != nil {
		log.Printf("Error adding data to firestore %v", err)
		return "", err
	}
	return resp, nil
}

// SendEmail is the resolver for the sendEmail field.
func (r *mutationResolver) SendEmail(ctx context.Context, to []*string, senderName string, userName []*string, body string, templateID string) ([]string, error) {
	resp, err := handlers.SendEmail(ctx, to, senderName, userName, body, templateID)
	if err != nil {
		var temp string
		return []string{temp}, err
	}
	return resp, nil
}

// GetFCMToken is the resolver for the getFCMToken field.
func (r *mutationResolver) GetFCMToken(ctx context.Context) (string, error) {
	resp, err := handlers.GetFCMToken(ctx)
	if err != nil {
		log.Printf("Unable to map UserID with FCM token")
		return "", err
	}
	return resp, nil
}

// AuthTokens is the resolver for the Auth_tokens field.
func (r *mutationResolver) AuthTokens(ctx context.Context) (string, error) {
	resp, err := handlers.Auth_tokens(ctx)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return resp, nil
}

// SendEmailUserID is the resolver for the sendEmail_UserId field.
func (r *mutationResolver) SendEmailUserID(ctx context.Context, userID []*string, senderName string, userName []*string, body string, templateID string) ([]string, error) {
	resp, err := handlers.SendEmailToUserIds(ctx, userID, senderName, userName, body, templateID)
	if err != nil {
		log.Println(err)
		var temp string
		return []string{temp}, nil
	}
	return resp, nil
}

// AddUserTags is the resolver for the addUserTags field.
func (r *mutationResolver) AddUserTags(ctx context.Context, ids []*model.UserDetails, tags []*string) (*bool, error) {
	resp, err := handlers.AddUserTags(ctx, ids, tags)
	if err != nil {
		log.Printf("Got error while setting uesr tags: %v", err)
		return nil, err
	}
	return resp, nil
}

// AddClassroomFlags is the resolver for the addClassroomFlags field.
func (r *mutationResolver) AddClassroomFlags(ctx context.Context, input *model.ClassRoomFlagsInput) (*model.ClassRoomFlags, error) {
	resp, err := handlers.AddClassroomFlags(ctx, input)
	if err != nil {
		log.Printf("Got error while setting topic classroom flags: %v", err)
		return nil, err
	}
	return resp, nil
}

// AddMessagesMeet is the resolver for the addMessagesMeet field.
func (r *mutationResolver) AddMessagesMeet(ctx context.Context, message *model.Messages) (*bool, error) {
	resp, err := handlers.AddMessagesMeet(ctx, message)
	if err != nil {
		log.Printf("Got error while sending messages: %v", err)
		return nil, err
	}
	return resp, nil
}

// AddPoll is the resolver for the addPoll field.
func (r *mutationResolver) AddPoll(ctx context.Context, input *model.PollsInput) (*model.Polls, error) {
	resp, err := handlers.AddPoll(ctx, input)
	if err != nil {
		log.Printf("Got error while adding polls: %v", err)
		return nil, err
	}
	return resp, nil
}

// UpdatePoll is the resolver for the updatePoll field.
func (r *mutationResolver) UpdatePoll(ctx context.Context, input *model.PollsInput) (*model.Polls, error) {
	resp, err := handlers.UpdatePoll(ctx, input)
	if err != nil {
		log.Printf("Got error while updating polls: %v", err)
		return nil, err
	}
	return resp, nil
}

// UpdatePollOptions is the resolver for the updatePollOptions field.
func (r *mutationResolver) UpdatePollOptions(ctx context.Context, input *model.PollResponseInput) (*model.PollResponse, error) {
	resp, err := handlers.UpdatePollOptions(ctx, input)
	if err != nil {
		log.Printf("Got error while updating poll options: %v", err)
		return nil, err
	}
	return resp, nil
}

// AddQuizToClassroomFlags is the resolver for the addQuizToClassroomFlags field.
func (r *mutationResolver) AddQuizToClassroomFlags(ctx context.Context, input *model.PublishedQuiz) (*bool, error) {
	resp, err := handlers.AddQuizToClassroomFlags(ctx, input)
	if err != nil {
		log.Printf("Got error while adding quiz to classroom flags: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetAll is the resolver for the getAll field.
func (r *queryResolver) GetAll(ctx context.Context, prevPageSnapShot string, pageSize int, isRead *bool) (*model.PaginatedNotifications, error) {
	resp, err := handlers.GetAllNotifications(ctx, prevPageSnapShot, pageSize, isRead)
	if err != nil {
		log.Printf("Got error while getting notifications: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetAllPaginatedNotifications is the resolver for the getAllPaginatedNotifications field.
func (r *queryResolver) GetAllPaginatedNotifications(ctx context.Context, pageIndex int, pageSize int, isRead *bool) ([]*model.FirestoreMessage, error) {
	resp, err := handlers.GetAllPaginatedNotifications(ctx, pageIndex, pageSize, isRead)
	if err != nil {
		log.Printf("Got error while getting notifications: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetUserLspIDTags is the resolver for the getUserLspIdTags field.
func (r *queryResolver) GetUserLspIDTags(ctx context.Context, userLspID []*string) ([]*model.TagsData, error) {
	resp, err := handlers.GetUserLspIDTags(ctx, userLspID)
	if err != nil {
		log.Printf("Got error while getting user lsp id tags: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetTagUsers is the resolver for the getTagUsers field.
func (r *queryResolver) GetTagUsers(ctx context.Context, prevPageSnapShot *string, pageSize *int, tags []*string) (*model.PaginatedTagsData, error) {
	resp, err := handlers.GetTagUsers(ctx, prevPageSnapShot, pageSize, tags)
	if err != nil {
		log.Printf("Got error while getting tags: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetPollResults is the resolver for the getPollResults field.
func (r *queryResolver) GetPollResults(ctx context.Context, pollID *string) (*model.PollResults, error) {
	resp, err := handlers.GetPollResults(ctx, pollID)
	if err != nil {
		log.Printf("Got error while getting poll results: %v", err)
		return nil, err
	}
	return resp, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
