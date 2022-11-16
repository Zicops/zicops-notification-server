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
func (r *mutationResolver) SendNotification(ctx context.Context, notification model.NotificationInput) (*model.Notification, error) {
	log.Println("Reached here")
	resp, err := handlers.SendNotification(ctx, notification)
	if err != nil {
		log.Printf("Error sending notification %v", err)
		return nil, err
	}
	return resp, err
}

// Dummy is the resolver for the dummy field.
func (r *queryResolver) Dummy(ctx context.Context, name *string) (*string, error) {
	return name, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
