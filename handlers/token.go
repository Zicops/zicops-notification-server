package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/zicops/zicops-notification-server/global"
	"google.golang.org/api/iterator"
)

func GetFCMToken(ctx context.Context) (string, error) {
	//log.Println("Reached here")
	global.Ct = ctx
	fcm_token := fmt.Sprintf("%s", ctx.Value("fcm-token"))
	iter := global.Client.Collection("tokens").Where("FCM-token", "==", fcm_token).Documents(global.Ct)
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			//log.Println("We have reached inside our so callled main function")
			claims, _ := GetClaimsFromContext(global.Ct)
			email_creator := claims["email"].(string)
			userId := base64.StdEncoding.EncodeToString([]byte(email_creator))

			//now we have both userID and fcm_token for a user, just fucking map them
			_, _, err := global.Client.Collection("tokens").Add(global.Ct, map[string]interface{}{
				"UserID":    userId,
				"FCM-token": fcm_token,
			})

			if err != nil {
				log.Println("Unable to add data to firestore database")
				return "", err
			}
			return "Tokens added successfully", nil
		} else if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
			return "", err
		} else {
			return "Tokens already present", nil
		}

	}

}
