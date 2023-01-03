package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
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

// SendNotificationWith is the resolver for the sendNotificationWith field.
func (r *mutationResolver) SendNotificationWith(ctx context.Context, notification model.NotificationInput, link string) ([]*model.Notification, error) {
	panic(fmt.Errorf("not implemented: SendNotificationWith - sendNotificationWith"))
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

// GetAll is the resolver for the getAll field.
func (r *queryResolver) GetAll(ctx context.Context, prevPageSnapShot string, pageSize int, isRead *bool) (*model.PaginatedNotifications, error) {
	resp, err := handlers.GetAllNotifications(ctx, prevPageSnapShot, pageSize, isRead)
	if err != nil {
		log.Println("Error receiving notification list")
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

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *mutationResolver) SendNotification(ctx context.Context, notification model.NotificationInput) ([]*model.Notification, error) {
	return nil, nil
}
