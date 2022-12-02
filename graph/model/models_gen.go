// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type FirestoreData struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt int    `json:"created_at"`
	UserID    string `json:"user_id"`
	IsRead    bool   `json:"is_read"`
}

type FirestoreDataInput struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	UserID    string `json:"user_id"`
	IsRead    bool   `json:"is_read"`
	MessageID string `json:"message_id"`
}

type FirestoreMessage struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt int    `json:"created_at"`
	UserID    string `json:"user_id"`
}

type Notification struct {
	Statuscode string `json:"statuscode"`
}

type NotificationInput struct {
	Title  string    `json:"title"`
	Body   string    `json:"body"`
	Emails []*string `json:"emails"`
}

type PaginatedNotifications struct {
	Messages         []*FirestoreMessage `json:"messages"`
	NextPageSnapShot *string             `json:"nextPageSnapShot"`
}
