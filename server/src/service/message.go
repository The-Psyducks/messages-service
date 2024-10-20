package service

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type MessageService struct {
}

func NewMessageService() (*MessageService, error) {

	return &MessageService{}, nil
}

func (ms *MessageService) SendMessage() error {

	ctx := context.Background() //firebase context
	conf := &firebase.Config{
		DatabaseURL: "https://twitsnap-fab5c-default-rtdb.firebaseio.com/",
	}

	// Fetch the service account key JSON file contents
	opt := option.WithCredentialsFile("twitsnap-fab5c-firebase-adminsdk-3qxha-c88972e6e9.json")

	// Initialize the app with a service account, granting admin privileges
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("Error initializing firebase app:", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}

	// As an admin, the app has access to read and write all data, regradless of Security Rules
	ref := client.NewRef("ejemplo")
	var data map[string]interface{}
	if err := ref.Get(ctx, &data); err != nil {
		log.Fatalln("Error reading from database:", err)
	}
	fmt.Println("data retrieved: ", data)

	dict := map[string]string{"clave": "valor"}

	if _, err := ref.Push(ctx, dict); err != nil {
		log.Fatalln("Error pushing message: ", err)
	}

	if err := ref.Get(ctx, &data); err != nil {
		log.Fatalln("Error reading from database:", err)
	}
	fmt.Println("data retrieved after inserting into db: ", data)

	return nil

}

//
