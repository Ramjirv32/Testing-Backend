package config

import (
	"context"
	"log"
	"os"

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

	// Priority 1: Use FIREBASE_CREDENTIALS environment variable
	if cfg.FirebaseCredentials != "" {
		log.Println("🔍 Using Firebase credentials from FIREBASE_CREDENTIALS environment variable")
		opt = option.WithCredentialsJSON([]byte(cfg.FirebaseCredentials))
	} else {
		// Priority 2: Use service account JSON file
		credFile := "ticpin-website-firebase-adminsdk-fbsvc-fade32c947.json"
		if _, err := os.Stat(credFile); err == nil {
			log.Println("🔍 Using Firebase credentials from file:", credFile)
			opt = option.WithCredentialsFile(credFile)
		} else {
			log.Fatal("❌ No Firebase credentials found. Set FIREBASE_CREDENTIALS env var or place service account JSON file in project root.")
		}
	}

	// Initialize Firebase App with storage bucket config
	firebaseConfig := &firebase.Config{
		StorageBucket: "ticpin-website.firebasestorage.app",
	}
	app, err := firebase.NewApp(ctx, firebaseConfig, opt)
	if err != nil {
		log.Fatalf("❌ Firebase init error: %v", err)
	}
	FirebaseApp = app

	// Initialize Auth
	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("❌ Firebase auth error: %v", err)
	}
	FirebaseAuth = authClient
	log.Println("✅ Firebase Auth initialized")

	// Initialize Firestore
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("❌ Firestore init error: %v", err)
	}
	FirestoreClient = firestoreClient
	log.Println("✅ Firestore initialized")

	// Initialize Storage
	storageClient, err := app.Storage(ctx)
	if err != nil {
		log.Printf("⚠️ Firebase storage error: %v", err)
	} else {
		FirebaseStorage = storageClient
		log.Println("✅ Firebase Storage initialized")
	}

	log.Println("🚀 Firebase Services Initialized Successfully")
}
