package service

import (
	"context"
	"log"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

type MessageService struct {
}

func NewMessageService() *MessageService {
	return &MessageService{}
}

func (ms *MessageService) SendMessage(senderId string, receiverId string, content string) error {

	client, ctx := createFirebaseDbClient()

	resourceRef := createMessage(senderId, receiverId, client)

	// if err := ref.Get(ctx, &data); err != nil {
	// 	log.Fatalln("Error reading from database:", err)
	// }
	// fmt.Println("data retrieved: ", data)

	msg := map[string]string{
		"from":      senderId,
		"to":        receiverId,
		"content":   content,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := resourceRef.Set(ctx, msg); err != nil {
		log.Fatalln("Error pushing message: ", err)
	}

	/*var data map[string]interface{}
	if err := resourceRef.Get(ctx, &data); err != nil {
		log.Fatalln("Error reading from database:", err)
	}

	log.Println("data retrieved after inserting into db: ", data)*/

	return nil
}

func createMessage(senderId string, receiverId string, client *db.Client) *db.Ref {
	firstUser, secondUser := func(a, b string) (string, string) {

		if a < b {
			return a, b
		}
		return b, a

	}(senderId, receiverId)

	newUUID := uuid.New().String()

	ref := client.NewRef("poc/msg-" + firstUser + "-" + secondUser + "/" + newUUID)
	return ref
}

func createFirebaseDbClient() (*db.Client, context.Context) {
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: "https://twitsnap-fab5c-default-rtdb.firebaseio.com/",
	}

	opt := option.WithCredentialsFile("twitsnap-fab5c-firebase-adminsdk-3qxha-c88972e6e9.json")

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("Error initializing firebase app:", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}
	return client, ctx
}

//
