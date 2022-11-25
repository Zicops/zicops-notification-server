package global

import (
	"context"
	"log"
	"os"

	firestore "cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"

	lib "github.com/zicops/zicops-notification-server/lib"
	"google.golang.org/api/option"
)

var (
	App    *firebase.App
	Client *firestore.Client
	Ct     = context.Background()
)

func init() {
	//os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "firebase-key.json")
	serviceAccountZicops := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if serviceAccountZicops == "" {
		log.Printf("failed to get right credentials for course creator")
	}
	targetScopes := []string{
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/userinfo.email",
	}
	currentCreds, _, err := lib.ReadCredentialsFile(Ct, serviceAccountZicops, targetScopes)
	if err != nil {
		log.Println(err)
	}

	opt := option.WithCredentials(currentCreds)
	App, err = firebase.NewApp(Ct, nil, opt)
	if err != nil {
		log.Printf("error initializing app: %v", err)
	}

	Client, err = App.Firestore(Ct)
	if err != nil {
		log.Printf("Error while initialising firestore %v", err)
	}
}
