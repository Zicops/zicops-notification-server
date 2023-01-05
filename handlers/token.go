package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/zicops/zicops-notification-server/global"
)

func GetFCMToken(ctx context.Context) (string, error) {
	//log.Println("Reached here")
	global.Ct = ctx
	claims, _ := GetClaimsFromContext(ctx)
	lsp := claims["lsp_id"].(string)
	fcm_token := fmt.Sprintf("%s", ctx.Value("fcm-token"))
	iter := global.Client.Collection("tokens").Where("FCM-token", "==", fcm_token).Where("LspID", "==", lsp).Documents(global.Ct)
	for {
		_, err := iter.Next()

		//see if there's no token present for given fcm-token, by checking if iterator is empty
		if err != nil && err.Error() == "no more items in iterator" {

			email_creator := claims["email"].(string)
			userId := base64.StdEncoding.EncodeToString([]byte(email_creator))

			//now we have both userID and fcm_token for a user, just map them
			_, _, err := global.Client.Collection("tokens").Add(global.Ct, map[string]interface{}{
				"UserID":    userId,
				"FCM-token": fcm_token,
				"LspID":     lsp,
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
