package service

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type MessageService struct {
	App *firebase.App
}

func NewMessageService() (*MessageService, error) {
	opt := option.WithCredentialsFile("path/to/serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}

	return &MessageService{App: app}, nil
}

func (ms *MessageService) SendMessage() {

	log.Println("Sending Message: Not yet implemented")

}
