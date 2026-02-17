package config

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var (
	FirebaseApp  *firebase.App
	FirebaseAuth *auth.Client
)

func InitFirebase(cfg *Config) {
	ctx := context.Background()

	opt := option.WithCredentialsFile(cfg.FirebaseKey)
	config := &firebase.Config{
		ProjectID:     "ticpin-website",
		StorageBucket: "ticpin-website.firebasestorage.app",
	}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalf("Firebase init error: %v", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("Firebase auth error: %v", err)
	}

	FirebaseApp = app
	FirebaseAuth = authClient
}
