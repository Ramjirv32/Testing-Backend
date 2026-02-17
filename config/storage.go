package config

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/storage"
)

var FirebaseStorage *storage.Client

func InitStorage(app *firebase.App) {
	ctx := context.Background()
	client, err := app.Storage(ctx)
	if err != nil {
		log.Fatalf("Firebase storage error: %v", err)
	}
	FirebaseStorage = client
}
