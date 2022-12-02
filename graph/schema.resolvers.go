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

// SendNotification is the resolver for the sendNotification field.
func (r *mutationResolver) SendNotification(ctx context.Context, notification model.NotificationInput) ([]*model.Notification, error) {
	resp, err := handlers.SendNotification(ctx, notification)
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

// GetFCMToken is the resolver for the getFCMToken field.
func (r *mutationResolver) GetFCMToken(ctx context.Context) (string, error) {
	resp, err := handlers.GetFCMToken(ctx)
	if err != nil {
		log.Printf("Unable to map UserID with FCM token")
		return "", err
	}
	return resp, nil
}

// GetAll is the resolver for the getAll field.
func (r *queryResolver) GetAll(ctx context.Context, prevPageSnapShot string, pageSize int) (*model.PaginatedNotifications, error) {
	resp, err := handlers.GetAllNotifications(ctx, prevPageSnapShot, pageSize)
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
