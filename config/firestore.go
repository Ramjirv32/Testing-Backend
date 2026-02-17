package config

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
)

var FirestoreClient *firestore.Client

func InitFirestore(app *firebase.App) {
	if app == nil {
		log.Println("⚠️ Firebase app is nil, skipping Firestore initialization")
		return
	}

	ctx := context.Background()
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Firestore init error: %v", err)
	}

	FirestoreClient = client
	log.Println("✅ Firestore initialized successfully")
}
