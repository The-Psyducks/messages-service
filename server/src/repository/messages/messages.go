// coverage:ignore
package repository

import (
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"os"
	"strings"
	"time"
)

type RealTimeDatabase struct {
}

func NewRealTimeDatabase() RealTimeDatabaseInterface {
	//if err := BuildFirebaseConfig(); err != nil {
	//	log.Fatalln("Error building firebase config:", err)
	//}

	return &RealTimeDatabase{}
}

type Message struct {
	Id        string `json:"id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

func (db *RealTimeDatabase) SendMessage(senderId string, receiverId string, content string) (string, error) {
	client, ctx := db.createFirebaseDbClient()

	resourceRef := db.createMessageRef(senderId, receiverId, client)

	msg := Message{
		Id:        "",
		From:      senderId,
		To:        receiverId,
		Content:   content,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	if _, err := resourceRef.Push(ctx, msg); err != nil {
		return "", err
	}

	return resourceRef.Path, nil

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
	} else {
		uri = "prod/" + uri
	}
	ref := client.NewRef(uri)
	return ref
}

func (db *RealTimeDatabase) createFirebaseDbClient() (*db.Client, context.Context) {
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: "https://twitsnap-fab5c-default-rtdb.firebaseio.com/",
	}
	//make opt with env vars insteaf of hardcoded path
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

func (db *RealTimeDatabase) GetChats(conversationId string) (*map[string]Message, error) {

	client, ctx := db.createFirebaseDbClient()
	uri := getUri() + conversationId

	ref := client.NewRef(uri)
	log.Println("Getting chat reference from uri: ", uri)
	chats := new(map[string]Message)
	if err := ref.Get(ctx, &chats); err != nil {
		return nil, err
	}

	return chats, nil

}

func getUri() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "HEROKU" {
		return "prod/"
	}
	return "test/"

}

func (db *RealTimeDatabase) GetConversations() ([]string, error) {
	client, ctx := db.createFirebaseDbClient()
	uri := getUri()

	ref := client.NewRef(uri)
	var data map[string]interface{}
	if err := ref.Get(ctx, &data); err != nil {
		return nil, err
	}
	conversations := make([]string, 0, len(data))

	// Iterate over the map and collect the keys
	for key := range data {
		conversations = append(conversations, key)
	}
	return conversations, nil
}

type FirebaseConfig struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
	UniverseDomain          string `json:"universe_domain"`
}

// BuildFirebaseConfig builds the Firebase configuration from environment variables
func BuildFirebaseConfig() error {
	fmt.Println("ENV VAR FOR CONFIG:", os.Getenv("SERVICE_ACCOUNT_PROJECT_ID"))
	privateKey := os.Getenv("SERVICE_ACCOUNT_PRIVATE_KEY")
	formatedPrivateKey := strings.ReplaceAll(privateKey, "\\n", "\n")
	configFile := &FirebaseConfig{
		Type:                    "service_account",
		ProjectID:               os.Getenv("SERVICE_ACCOUNT_PROJECT_ID"),
		PrivateKeyID:            os.Getenv("SERVICE_ACCOUNT_PRIVATE_KEY_ID"),
		PrivateKey:              formatedPrivateKey,
		ClientEmail:             os.Getenv("SERVICE_ACCOUNT_CLIENT_EMAIL"),
		ClientID:                os.Getenv("SERVICE_ACCOUNT_CLIENT_ID"),
		AuthURI:                 os.Getenv("SERVICE_ACCOUNT_AUTH_URI"),
		TokenURI:                os.Getenv("SERVICE_ACCOUNT_TOKEN_URI"),
		AuthProviderX509CertURL: os.Getenv("SERVICE_ACCOUNT_AUTH_PROVIDER_CERT_URL"),
		ClientX509CertURL:       os.Getenv("SERVICE_ACCOUNT_CLIENT_CERT_URL"),
		UniverseDomain:          os.Getenv("SERVICE_ACCOUNT_UNIVERSE_DOMAIN"),
	}

	jsonData, err := json.MarshalIndent(configFile, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}
	err = os.WriteFile("twitsnap-fab5c-firebase-adminsdk-3qxha-c88972e6e9.json", jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}
