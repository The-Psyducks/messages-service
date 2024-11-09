package firebases_connector

import (
	"context"
	"firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"google.golang.org/api/option"
	"log"
)

type FirebaseConnector struct{}

func (fc *FirebaseConnector) sendNotification(token, title, body string) error {

	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: "https://twitsnap-fab5c-default-rtdb.firebaseio.com/",
	}

	opt := option.WithCredentialsFile("twitsnap-fab5c-firebase-adminsdk-3qxha-c88972e6e9.json")

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("Error initializing firebase app:", err)
	}
	client, err := app.Messaging(ctx)

	if err != nil {
		log.Fatalln("Error initializing messaging client:", err)
	}

	message := &messaging.Message{
		Data: map[string]string{
			"deeplink": "dale juancito mandame el deep",
		},
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}
	response, err := client.Send(ctx, message)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println("firebase response", response)

	return nil

}

func (fc *FirebaseConnector) SendNotificationToUserDevices(devicesTokens []string, title, body string) error {
	for _, token := range devicesTokens {
		if err := fc.sendNotification(token, title, body); err != nil {
			return err
		}
	}
	return nil
}
