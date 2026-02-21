package config

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/storage"
	"google.golang.org/api/option"
)

var (
	FirebaseApp     *firebase.App
	FirebaseAuth    *auth.Client
	FirestoreClient *firestore.Client
	FirebaseStorage *storage.Client
)

func InitFirebase(cfg *Config) {
	ctx := context.Background()

	var opt option.ClientOption

	if cfg.FirebaseCredentials != "" {
		log.Println(" Using Firebase credentials from FIREBASE_CREDENTIALS environment variable")
		opt = option.WithCredentialsJSON([]byte(cfg.FirebaseCredentials))
	} else {
		credFile := cfg.FirebaseKeyPath

		// Check local path first
		if _, err := os.Stat(credFile); err == nil {
			log.Println(" Using Firebase credentials from file:", credFile)
			opt = option.WithCredentialsFile(credFile)
		} else {
			// Fallback to Render's secret path if it's just a filename
			renderPath := "/etc/secrets/" + credFile
			if _, err := os.Stat(renderPath); err == nil {
				log.Println(" Using Firebase credentials from Render secrets:", renderPath)
				opt = option.WithCredentialsFile(renderPath)
			} else {
				log.Fatalf(" No Firebase credentials found at %s or %s. Set FIREBASE_CREDENTIALS or FIREBASE_KEY_PATH env var.", credFile, renderPath)
			}
		}
	}

	firebaseConfig := &firebase.Config{
		StorageBucket: "ticpin-fa6d2.firebasestorage.app",
	}
	app, err := firebase.NewApp(ctx, firebaseConfig, opt)
	if err != nil {
		log.Fatalf(" Firebase init error: %v", err)
	}
	FirebaseApp = app

	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf(" Firebase auth error: %v", err)
	}
	FirebaseAuth = authClient
	log.Println(" Firebase Auth initialized")

	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf(" Firestore init error: %v", err)
	}
	FirestoreClient = firestoreClient

	healthCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err = firestoreClient.Collection("health").Doc("check").Get(healthCtx)
	if err != nil && !strings.Contains(err.Error(), "NotFound") && !strings.Contains(err.Error(), "not found") {
		log.Printf(" Firestore connection warning (continuing anyway): %v", err)
	} else {
		log.Println(" Firestore initialized and verified")
	}

	storageClient, err := app.Storage(ctx)
	if err != nil {
		log.Printf(" Firebase storage error: %v", err)
	} else {
		FirebaseStorage = storageClient
		log.Println(" Firebase Storage initialized")
	}

	log.Println(" Firebase Services Initialized Successfully")
}
