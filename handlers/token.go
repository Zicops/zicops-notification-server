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
	fcm_token := fmt.Sprintf("%s", ctx.Value("fcm-token"))

	claims, _ := GetClaimsFromContext(ctx)
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
}
