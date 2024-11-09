package repository

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"messages/src/repository/messages"
)

func (db *messages.RealTimeDatabase) sendNotification(token, title, body string) error {

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

func (db *messages.RealTimeDatabase) SendNotificationToUserDevices(devicesTokens []string, title, body string) error {
	for _, token := range devicesTokens {
		if err := db.sendNotification(token, title, body); err != nil {
			return err
		}
	}
	return nil
}
