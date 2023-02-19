package global

import (
	"context"
	"log"
	"os"
	"time"

	firestore "cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"

	"github.com/robfig/cron/v3"
	lib "github.com/zicops/zicops-notification-server/lib"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var (
	App       *firebase.App
	Client    *firestore.Client
	Ct        = context.Background()
	Messanger *messaging.Client
)

func init() {
	//os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "zicops-cc.json")
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

	messanger, err := App.Messaging(Ct)
	if err != nil {
		log.Printf("Error while initialising messaging: %v", err)
	}
	Messanger = messanger
	//scheduler
	deleteNotifications()
	//sch()
}

func deleteNotifications() {

	c := cron.New()
	c.AddFunc("0 2 * * 5", sch)
	c.AddFunc("20 0 * * *", deleteNullTokens)

	c.Start()
}

func deleteNullTokens() {
	var resp []map[string]interface{}
	var ids []string
	iter := Client.Collection("tokens").Documents(Ct)
	for {

		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
		}
		ids = append(ids, doc.Ref.ID)
		resp = append(resp, doc.Data())
	}
	for k, v := range resp {
		if v["FCM-token"].(string) == "null" || v["LspID"].(string) == "null" {
			_, err := Client.Collection("tokens").Doc(ids[k]).Delete(Ct)
			if err != nil {
				log.Println("Got error while deleting the data ", err)
			}
		}
	}
}

func sch() {
	var resp []map[string]interface{}
	var ids []string
	iter := Client.Collection("notification").Documents(Ct)
	for {

		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
		}
		ids = append(ids, doc.Ref.ID)
		resp = append(resp, doc.Data())
	}
	//here we just need to delete older notifications
	thirtydays := int64((time.Hour * 24 * 30).Seconds())
	cur := int64(time.Now().Unix())
	for k, v := range resp {
		//log.Println("Reached here")
		t := v["CreatedAt"].(int64)
		dif := cur - t
		if dif > thirtydays {
			_, err := Client.Collection("notification").Doc(ids[k]).Delete(Ct)
			if err != nil {
				log.Println("Got error while deleting the data ", err)
			}
		}
	}
}
