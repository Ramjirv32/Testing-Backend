package config

import (
	"context"
	"fmt"
	"log"
	"strings"

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

	var opt option.ClientOption

	if cfg.FirebasePrivateKey != "" && cfg.FirebaseClientEmail != "" {
		// Reconstruct service account JSON from env vars
		// We use standard Google URLs for the missing pieces
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
			strings.ReplaceAll(cfg.FirebasePrivateKey, "\n", "\\n"),
			cfg.FirebaseClientEmail,
			strings.ReplaceAll(cfg.FirebaseClientEmail, "@", "%%40")) // URL encoded for the cert URL

		opt = option.WithCredentialsJSON([]byte(jsonCreds))
	} else {
		log.Println("⚠️ Firebase environment variables missing. Firebase features will be disabled.")
		return
	}

	config := &firebase.Config{
		ProjectID:     cfg.FirebaseProjectID,
		StorageBucket: fmt.Sprintf("%s.firebasestorage.app", cfg.FirebaseProjectID),
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
	log.Println("✅ Firebase initialized successfully using environment variables")
}
