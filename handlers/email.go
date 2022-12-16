package handlers

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// "d-bf691d7c93794afca36c326cd032ccbf"
func SendEmail(ctx context.Context, to string, user_name string, from string, sender_name string, body string, template_id string) (string, error) {

	fromMail := mail.NewEmail(sender_name, from)
	toMail := mail.NewEmail(user_name, to)
	mailSetup := mail.NewV3Mail()
	mailSetup.SetFrom(fromMail)
	mailSetup.SetTemplateID(template_id)
	p := mail.NewPersonalization()
	p.AddTos(toMail)
	// Now we will set the data from the body and put it in some interface, decode it and put it in p.SetDynamicTemplateData
	var bodyData map[string]string
	err := json.Unmarshal([]byte(body), &bodyData)
	if err != nil {
		log.Println(err)
		return "", nil
	}
	log.Println("Values for k and v are as given")
	for k, v := range bodyData {
		//log.Println(k, "    ", v)
		p.SetDynamicTemplateData(k, v)
	}
	mailSetup.AddPersonalizations(p)
	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = mail.GetRequestBody(mailSetup)
	request.Body = Body
	response, err := sendgrid.API(request)
	/*

			request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
			request.Method = "POST"
			request.Body = []byte(` {
			"personalizations": [
				{
					"to": [
						{
							"email": "joy@zicops.com"
						}
					],
					"subject": "Hello sir"
				}
			],
			"from": {
				"email": "no_reply@zicops.com"
			},
			"content": [
				{
					"type": "text/plain",
					"value": "Sir it is still giving error code 202 but sending mail, but we can look into this, I have made a few changes"
				}
			]
		}`)
			log.Println(os.Getenv("SENDGRID_API_KEY"))
			response, err := sendgrid.API(request)
	*/
	if err != nil {
		log.Println(err)
	} else {

		log.Println(response.Body)
		log.Println(response.StatusCode)

	}
	return strconv.Itoa(response.StatusCode), nil

}
