package repository

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"log"
	"os"
	"time"
)

type RealTimeDatabase struct {
}

func NewRealTimeDatabase() *RealTimeDatabase {
	return &RealTimeDatabase{}
}

func (db *RealTimeDatabase) SendMessage(senderId string, receiverId string, content string) error {
	client, ctx := db.createFirebaseDbClient()

	resourceRef := db.createMessageRef(senderId, receiverId, client)

	// if err := ref.Get(ctx, &data); err != nil {
	// 	log.Fatalln("Error reading from database:", err)
	// }
	// fmt.Println("data retrieved: ", data)

	msg := map[string]string{
		"id":        uuid.New().String(),
		"from":      senderId,
		"to":        receiverId,
		"content":   content,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}

	if _, err := resourceRef.Push(ctx, msg); err != nil {
		return err
	}

	/*var data map[string]interface{}
	if err := resourceRef.Get(ctx, &data); err != nil {
		log.Fatalln("Error reading from database:", err)
	}

	log.Println("data retrieved after inserting into db: ", data)*/

	return nil
}

func (db *RealTimeDatabase) createMessageRef(senderId string, receiverId string, client *db.Client) *db.Ref {
	firstUser, secondUser := func(a, b string) (string, string) {
		if a < b {
			return a, b
		}
		return b, a
	}(senderId, receiverId)
	uri := "dm-" + firstUser + "-" + secondUser
	if os.Getenv("ENVIRONMENT") == "test" {
		uri = "test/" + uri
	}
	ref := client.NewRef(uri)
	return ref
}

func (db *RealTimeDatabase) createFirebaseDbClient() (*db.Client, context.Context) {
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

type RealTimeDatabaseInterface interface {
	SendMessage(senderId string, receiverId string, content string) error
}
