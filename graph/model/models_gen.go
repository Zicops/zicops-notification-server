// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Notification struct {
	Statuscode string `json:"statuscode"`
}

type NotificationInput struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Token string `json:"token"`
}
