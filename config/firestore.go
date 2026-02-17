package config

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

var FirestoreClient *firestore.Client

func InitFirestore(cfg *Config) {
	ctx := context.Background()

	opt := option.WithCredentialsFile(cfg.FirebaseKey)
	client, err := firestore.NewClient(ctx, "ticpin-website", opt)
	if err != nil {
		log.Fatalf("Firestore init error: %v", err)
	}

	FirestoreClient = client
	log.Println("Firestore initialized successfully")
}
