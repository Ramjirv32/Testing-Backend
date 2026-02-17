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

	// Hardcoded Firebase credentials
	const (
		projectID   = "ticpin-website"
		clientEmail = "firebase-adminsdk-fbsvc@ticpin-website.iam.gserviceaccount.com"
		privateKey  = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDoH4MWLDk/ASUe
BttI1ym+/z+51ffC+Gm/vEvAXGzTluZ2x+ZRGS9HQA0oqpZbz4zdC+SDfDmwaPO4
AfEGeztIP34/oTZwZLpF8i8IKye4zk/3lLbWm+aNMzTxI8JKB4oxNT374jhQ2ybP
2yFMLXALXMlm6zDDxvn/xahgdZdW7FklWBKbIIC9y9R02buLR1DHe69Vvl5fktPg
LVIr/l7ilN+IaL3KZFDqtQkUJ+iNTkOMQOnT7lcT5zGYUmdOrVRSg7ygGmvySLho
SPBoGafEIJsXQ7H4mnSNRxuPIgPm/NXz48N7WOXb5YSo58DKx7cau7BqniXOb224
vqe8Q837AgMBAAECggEAAy6vauXnUQRQgHVim0CL63jvZDpZP7yNIppPxY7e1RXM
ChPahgEc41Ku+4A/OHoDeeJYWy8gUVlXAg5QwiB9YxOvxOqOZwMShLP5zhhdXozB
jujkitOvWP87OhUd7ErnK56Jv4LN99nRUec0sSksUJOQlU8jJ4P6WHXaxZvHG+Nx
9CA3Ypp+0Y7cAgXqqZ7z6LC9Gc1tuGLEquOyzhTU+PIP5t5BeBX0AY6xLehdZZKJ
cf7yuWOCET8teh6kD2l65N3lup9touA2Nx3ZhF7vdVF1m7rJpJDGqyi6HgMWYQ+t
9v6Y5ovilUWn3dEHJ1i5jGP1Ut2//Oj4EpAPgt0iHQKBgQD3KTl4S4hgtW+0ZltT
85uB/uBADcysrv9sL8thhBUjWn9gCQ1mFk+6kDQUpQdQ/RHS3DH5ar+fHSFVGAkx
9KjjtKbwSWPpUnP47Ny/fkjf2bX4jYcCMCVi+Pu1tCOnGOm6E2sUEE9D1luM9OAX
IGZXt1xFuvattv3KjPchBbUjrwKBgQDwbJ0/C6zN0XZRqoIyuq7u8E+0payYPgqg
6bCFuz8e0XDbsFFqshuvI3U0p6vUsk2HwcDzcvIvjrUhOXMuyI33/zjMwqcNAivo
RkRDp7pSK/iwKQldxKJN27JRRewon1TGfnEFtS1+7ojuTw5d6arE7ddy/7INfHXG
EIgcUToxdQKBgC6TaC8RHMwMpNY8C63QVFe07hFkCFPqTlvWzd68gzc8UJCKZCn+
vluL3SSezLgoWHmB4TD9OssDNErS0rjFQCZY3rSdP+SyEwSvrhGv/I+ieTYzhWOW
KxVxkg11uto8SZ81FZKcWDOSa4IuiyQQiPiypwLE7sNhnoXS9qcUakQlAoGBAII/
VTDCcmtN/ntflAlHeV2YcpW66zXO5pMmBqtsNVXMwQdDDdhvhO/slaJg84XW0omp
PY6lxu5csWO+a9f8bmzbpznGehliA8dhybmdNCMwDxngIWLbE9J6IrBE4RtgtdyS
w0gETxFkyGnSCkZ2QD1PXFjAjQUhV+xlKFeu6YfBAoGANMAEx7iTzh1IuneKp7O4
6n7yIrlOzKZc+ViUY4GCvyQyKT/nZEFnn19tca5vXmptbGdLIbrosaMumLEy+DaL
91t+/1MfX+yjrxLDET4mHdJlNbxapKz673/5oJPbXRV8vmDQVvKpOpjQg2HzaLaL
yP/tei69eh95RW0EBXXhYzU=
-----END PRIVATE KEY-----
`
	)

	log.Println("🔍 Initializing Firebase with hardcoded credentials...")

	// Build the service account JSON
	jsonCreds := fmt.Sprintf(`{
		"type": "service_account",
		"project_id": "%s",
		"private_key": "%s",
		"client_email": "%s",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-fbsvc%%40ticpin-website.iam.gserviceaccount.com"
	}`,
		projectID,
		strings.ReplaceAll(privateKey, "\n", "\\n"),
		clientEmail)

	opt := option.WithCredentialsJSON([]byte(jsonCreds))

	firebaseCfg := &firebase.Config{
		ProjectID:     projectID,
		StorageBucket: fmt.Sprintf("%s.firebasestorage.app", projectID),
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

	log.Println("🚀 Firebase Services Initialized with Hardcoded Credentials")
}
