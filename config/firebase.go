package config

import (
	"context"
	"log"

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

	// Check if we have hardcoded credentials (fallback) or environment variable
	var credentialsJSON string

	if cfg.FirebaseCredentials != "" {
		// Use environment variable (RECOMMENDED for production)
		credentialsJSON = cfg.FirebaseCredentials
		log.Println("🔍 Using Firebase credentials from FIREBASE_CREDENTIALS environment variable")
	} else {
		// Fallback to hardcoded credentials (for development/testing)
		log.Println("🔍 Using hardcoded Firebase credentials")
		credentialsJSON = `{
  "type": "service_account",
  "project_id": "ticpin-website",
  "private_key_id": "79b256dff03bdbbcfbbc950f59d6665471d420f9",
  "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDoH4MWLDk/ASUe\nBttI1ym+/z+51ffC+Gm/vEvAXGzTluZ2x+ZRGS9HQA0oqpZbz4zdC+SDfDmwaPO4\nAfEGeztIP34/oTZwZLpF8i8IKye4zk/3lLbWm+aNMzTxI8JKB4oxNT374jhQ2ybP\n2yFMLXALXMlm6zDDxvn/xahgdZdW7FklWBKbIIC9y9R02buLR1DHe69Vvl5fktPg\nLVIr/l7ilN+IaL3KZFDqtQkUJ+iNTkOMQOnT7lcT5zGYUmdOrVRSg7ygGmvySLho\nSPBoGafEIJsXQ7H4mnSNRxuPIgPm/NXz48N7WOXb5YSo58DKx7cau7BqniXOb224\nvqe8Q837AgMBAAECggEAAy6vauXnUQRQgHVim0CL63jvZDpZP7yNIppPxY7e1RXM\nChPahgEc41Ku+4A/OHoDeeJYWy8gUVlXAg5QwiB9YxOvxOqOZwMShLP5zhhdXozB\njujkitOvWP87OhUd7ErnK56Jv4LN99nRUec0sSksUJOQlU8jJ4P6WHXaxZvHG+Nx\n9CA3Ypp+0Y7cAgXqqZ7z6LC9Gc1tuGLEquOyzhTU+PIP5t5BeBX0AY6xLehdZZKJ\ncf7yuWOCET8teh6kD2l65N3lup9touA2Nx3ZhF7vdVF1m7rJpJDGqyi6HgMWYQ+t\n9v6Y5ovilUWn3dEHJ1i5jGP1Ut2//Oj4EpAPgt0iHQKBgQD3KTl4S4hgtW+0ZltT\n85uB/uBADcysrv9sL8thhBUjWn9gCQ1mFk+6kDQUpQdQ/RHS3DH5ar+fHSFVGAkx\n9KjjtKbwSWPpUnP47Ny/fkjf2bX4jYcCMCVi+Pu1tCOnGOm6E2sUEE9D1luM9OAX\nIGZXt1xFuvattv3KjPchBbUjrwKBgQDwbJ0/C6zN0XZRqoIyuq7u8E+0payYPgqg\n6bCFuz8e0XDbsFFqshuvI3U0p6vUsk2HwcDzcvIvjrUhOXMuyI33/zjMwqcNAivo\nRkRDp7pSK/iwKQldxKJN27JRRewon1TGfnEFtS1+7ojuTw5d6arE7ddy/7INfHXG\nEIgcUToxdQKBgC6TaC8RHMwMpNY8C63QVFe07hFkCFPqTlvWzd68gzc8UJCKZCn+\nvluL3SSezLgoWHmB4TD9OssDNErS0rjFQCZY3rSdP+SyEwSvrhGv/I+ieTYzhWOW\nKxVxkg11uto8SZ81FZKcWDOSa4IuiyQQiPiypwLE7sNhnoXS9qcUakQlAoGBAII/\nVTDCcmtN/ntflAlHeV2YcpW66zXO5pMmBqtsNVXMwQdDDdhvhO/slaJg84XW0omp\nPY6lxu5csWO+a9f8bmzbpznGehliA8dhybmdNCMwDxngIWLbE9J6IrBE4RtgtdyS\nw0gETxFkyGnSCkZ2QD1PXFjAjQUhV+xlKFeu6YfBAoGANMAEx7iTzh1IuneKp7O4\n6n7yIrlOzKZc+ViUY4GCvyQyKT/nZEFnn19tca5vXmptbGdLIbrosaMumLEy+DaL\n91t+/1MfX+yjrxLDET4mHdJlNbxapKz673/5oJPbXRV8vmDQVvKpOpjQg2HzaLaL\nyP/tei69eh95RW0EBXXhYzU=\n-----END PRIVATE KEY-----\n",
  "client_email": "firebase-adminsdk-fbsvc@ticpin-website.iam.gserviceaccount.com",
  "client_id": "102027578202972861068",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-fbsvc%40ticpin-website.iam.gserviceaccount.com",
  "universe_domain": "googleapis.com"
}`
	}

	// Pass credentials directly to Firebase - NO MANIPULATION
	opt := option.WithCredentialsJSON([]byte(credentialsJSON))

	// Initialize Firebase App
	app, err := firebase.NewApp(ctx, nil, opt)
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
