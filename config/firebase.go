package config

import (
	"context"
	"fmt"
	"log"
	"strings"

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

	// Debug logging for credentials (sanitized)
	log.Printf("🔍 Firebase Config Check:")
	log.Printf("  - Project ID: %s", cfg.FirebaseProjectID)
	log.Printf("  - Client Email: %s", cfg.FirebaseClientEmail)
	log.Printf("  - Private Key Length: %d chars", len(cfg.FirebasePrivateKey))

	if cfg.FirebasePrivateKey == "" || cfg.FirebaseClientEmail == "" || cfg.FirebaseProjectID == "" {
		log.Println("⚠️ Firebase environment variables missing. Firebase features will be disabled.")
		log.Println("   Required: FIREBASE_PROJECT_ID, FIREBASE_PRIVATE_KEY, FIREBASE_CLIENT_EMAIL")
		return
	}

	// Handle private key - convert literal \n to actual newlines
	privateKey := cfg.FirebasePrivateKey

	// If the key contains literal \n (from environment variable), replace with actual newlines
	if strings.Contains(privateKey, "\\n") {
		privateKey = strings.ReplaceAll(privateKey, "\\n", "\n")
	}

	// Ensure the key starts and ends correctly
	if !strings.HasPrefix(privateKey, "-----BEGIN PRIVATE KEY-----") {
		log.Println("❌ Private key format is invalid - missing header")
		return
	}

	// Build the service account JSON
	jsonCreds := fmt.Sprintf(`{
		"type": "service_account",
		"project_id": "%s",
		"private_key": "%s",
		"client_email": "%s",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/%s"
	}`,
		cfg.FirebaseProjectID,
		strings.ReplaceAll(privateKey, "\n", "\\n"),
		cfg.FirebaseClientEmail,
		strings.ReplaceAll(cfg.FirebaseClientEmail, "@", "%40"))

	opt = option.WithCredentialsJSON([]byte(jsonCreds))

	firebaseCfg := &firebase.Config{
		ProjectID:     cfg.FirebaseProjectID,
		StorageBucket: fmt.Sprintf("%s.firebasestorage.app", cfg.FirebaseProjectID),
	}

	// 1. Initialize Firebase App
	app, err := firebase.NewApp(ctx, firebaseCfg, opt)
	if err != nil {
		log.Printf("❌ Firebase init error: %v\n", err)
		return
	}
	FirebaseApp = app

	// 2. Initialize Auth
	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Printf("❌ Firebase auth error: %v\n", err)
	} else {
		FirebaseAuth = authClient
		log.Println("✅ Firebase Auth initialized")
	}

	// 3. Initialize Firestore
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Printf("❌ Firestore init error: %v\n", err)
	} else {
		FirestoreClient = firestoreClient
		log.Println("✅ Firestore initialized")
	}

	// 4. Initialize Storage
	storageClient, err := app.Storage(ctx)
	if err != nil {
		log.Printf("❌ Firebase storage error: %v\n", err)
	} else {
		FirebaseStorage = storageClient
		log.Println("✅ Firebase Storage initialized")
	}

	log.Println("🚀 Firebase Services Sync Complete")
}
