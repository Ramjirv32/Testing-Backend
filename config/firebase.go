package config

import (
	"context"
	"log"
	"os"

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

	// Check if file exists
	if _, err := os.Stat(cfg.FirebaseKey); os.IsNotExist(err) {
		log.Printf("⚠️ Firebase key file not found at %s. Firebase features will be disabled.\n", cfg.FirebaseKey)
		return
	}

	opt := option.WithCredentialsFile(cfg.FirebaseKey)
	config := &firebase.Config{
		ProjectID:     "ticpin-website",
		StorageBucket: "ticpin-website.firebasestorage.app",
	}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Printf("❌ Firebase init error: %v\n", err)
		return
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Printf("❌ Firebase auth error: %v\n", err)
		return
	}

	FirebaseApp = app
	FirebaseAuth = authClient
	log.Println("✅ Firebase initialized successfully")
}
