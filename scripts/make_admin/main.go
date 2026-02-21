package main

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/v4/auth"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/iterator"
)

func main() {
	ctx := context.Background()

	_ = godotenv.Load()

	cfg := config.LoadConfig()
	config.InitFirebase(cfg)

	client := config.FirestoreClient
	authClient := config.FirebaseAuth

	if client == nil {
		log.Fatal("Firestore client is nil")
	}
	if authClient == nil {
		log.Fatal("Firebase Auth client is nil")
	}

	// ─── 1. Grant admin by phone ─────────────────────────────────────────
	grantAdminByPhone(ctx, client, "6383667872")

	// ─── 2. Ensure admin@ticpin.in exists with password and grant admin ──
	ensureAdminEmail(ctx, client, authClient, "admin@ticpin.in", "admin@ticpin")

	fmt.Println("\n✅ Done.")
}

// grantAdminByPhone finds a user by phone and sets is_admin=true.
func grantAdminByPhone(ctx context.Context, client *firestore.Client, phone string) {
	iter := client.Collection("users").Where("phone", "==", phone).Limit(1).Documents(ctx)
	doc, err := iter.Next()

	if err == iterator.Done {
		fmt.Printf("⚠️  No user found with phone %s\n", phone)
		return
	}
	if err != nil {
		log.Printf("❌ Error querying phone %s: %v\n", phone, err)
		return
	}

	var user models.User
	if err := doc.DataTo(&user); err != nil {
		log.Printf("❌ Failed to parse user: %v\n", err)
		return
	}

	if user.IsAdmin {
		fmt.Printf("ℹ️  Phone %s (%s) is already admin\n", phone, user.ID)
		return
	}

	_, err = client.Collection("users").Doc(user.ID).Update(ctx, []firestore.Update{
		{Path: "is_admin", Value: true},
		{Path: "updated_at", Value: time.Now()},
	})
	if err != nil {
		log.Printf("❌ Failed to set admin for phone %s: %v\n", phone, err)
		return
	}
	fmt.Printf("✅ Granted admin to phone %s (user ID: %s, name: %s)\n", phone, user.ID, user.Name)
}

// ensureAdminEmail creates or updates the Firebase Auth user for the given email/password,
// then ensures the Firestore user record has is_admin=true.
func ensureAdminEmail(ctx context.Context, client *firestore.Client, authClient *auth.Client, email, password string) {
	// ── Firebase Auth: create or update ────────────────────────────────
	var firebaseUID string

	existing, err := authClient.GetUserByEmail(ctx, email)
	if err != nil {
		if auth.IsUserNotFound(err) {
			// Create new Firebase Auth user
			params := (&auth.UserToCreate{}).
				Email(email).
				Password(password).
				EmailVerified(true).
				DisplayName("Admin")
			newUser, createErr := authClient.CreateUser(ctx, params)
			if createErr != nil {
				log.Printf("❌ Failed to create Firebase Auth user %s: %v\n", email, createErr)
				return
			}
			firebaseUID = newUser.UID
			fmt.Printf("✅ Created Firebase Auth user %s (UID: %s)\n", email, firebaseUID)
		} else {
			log.Printf("❌ Error looking up Firebase Auth user %s: %v\n", email, err)
			return
		}
	} else {
		// Update password
		firebaseUID = existing.UID
		params := (&auth.UserToUpdate{}).
			Password(password).
			EmailVerified(true)
		_, updateErr := authClient.UpdateUser(ctx, firebaseUID, params)
		if updateErr != nil {
			log.Printf("⚠️  Could not update password for %s: %v (continuing)\n", email, updateErr)
		} else {
			fmt.Printf("✅ Updated Firebase Auth password for %s (UID: %s)\n", email, firebaseUID)
		}
	}

	// Hash the password for organizer/email+password login
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("  Failed to hash password: %v\n", err)
		return
	}
	hashedPassword := string(hashed)

	// ── Firestore: find or create user, then set is_admin + password ────
	iter := client.Collection("users").Where("email", "==", email).Limit(1).Documents(ctx)
	doc, err := iter.Next()

	if err == iterator.Done {
		// No Firestore record — create one
		userID := utils.GenerateUUIDv7()
		newUser := models.User{
			ID:              userID,
			FirebaseUID:     firebaseUID,
			Email:           email,
			Name:            "Admin",
			Password:        hashedPassword,
			IsAdmin:         true,
			IsEmailVerified: true,
			IsOrganizer:     true,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		_, createErr := client.Collection("users").Doc(userID).Set(ctx, newUser)
		if createErr != nil {
			log.Printf("  Failed to create Firestore user for %s: %v\n", email, createErr)
			return
		}
		fmt.Printf("  Created Firestore user for %s (ID: %s) with admin=true, password set\n", email, userID)
		return
	}
	if err != nil {
		log.Printf("  Error querying Firestore for %s: %v\n", email, err)
		return
	}

	var user models.User
	_ = doc.DataTo(&user)

	updates := []firestore.Update{
		{Path: "is_admin", Value: true},
		{Path: "is_email_verified", Value: true},
		{Path: "is_organizer", Value: true},
		{Path: "password", Value: hashedPassword},
		{Path: "updated_at", Value: time.Now()},
	}
	if user.FirebaseUID == "" && firebaseUID != "" {
		updates = append(updates, firestore.Update{Path: "firebase_uid", Value: firebaseUID})
	}

	_, err = client.Collection("users").Doc(user.ID).Update(ctx, updates)
	if err != nil {
		log.Printf("  Failed to update Firestore user for %s: %v\n", email, err)
		return
	}
	fmt.Printf("  Patched %s (ID: %s): is_admin=true, is_organizer=true, is_email_verified=true, password=updated\n", email, user.ID)
}
