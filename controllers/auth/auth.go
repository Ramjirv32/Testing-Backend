package auth

import (
	"backend/config"
	"backend/models"
	"backend/repository"
	"backend/utils"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
)

type OTPData struct {
	OTP       string
	ExpiresAt time.Time
}

var userRepo = repository.NewUserRepository()
var emailOTPs = make(map[string]OTPData)
var emailOTPsMu sync.Mutex

// generateSecureOTP produces a cryptographically random 6-digit OTP string.
func generateSecureOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func SendEmailOTP(c fiber.Ctx) error {
	var req models.SendEmailOTPRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	otp, err := generateSecureOTP()
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to generate OTP")
	}

	emailOTPsMu.Lock()
	emailOTPs[req.Email] = OTPData{
		OTP:       otp,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	emailOTPsMu.Unlock()

	body := utils.GetOTPEmailTemplate(otp)

	err = utils.SendOTPEmail(utils.EmailAdmin, req.Email, "Email Verification - TicPin", body, otp)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to send email")
	}

	return utils.SuccessResponse(c, 200, "OTP sent successfully", nil)
}

func VerifyEmail(c fiber.Ctx) error {
	var req models.VerifyEmailRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	emailOTPsMu.Lock()
	storedOTP, ok := emailOTPs[req.Email]
	emailOTPsMu.Unlock()
	if !ok {
		return utils.ErrorResponse(c, 401, "No OTP sent for this email")
	}

	if time.Now().After(storedOTP.ExpiresAt) {
		emailOTPsMu.Lock()
		delete(emailOTPs, req.Email)
		emailOTPsMu.Unlock()
		return utils.ErrorResponse(c, 401, "OTP has expired")
	}

	if storedOTP.OTP != req.OTP {
		return utils.ErrorResponse(c, 401, "Invalid OTP code")
	}

	userID := c.Locals("uid").(string)
	user, err := userRepo.FindByID(c.Context(), userID)
	if err != nil || user == nil {
		return utils.ErrorResponse(c, 404, "User not found")
	}

	user.Email = req.Email
	user.IsEmailVerified = true
	user.UpdatedAt = time.Now()

	if err := userRepo.Update(c.Context(), user); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to update profile")
	}

	emailOTPsMu.Lock()
	delete(emailOTPs, req.Email)
	emailOTPsMu.Unlock()

	return utils.SuccessResponse(c, 200, "Email verified and profile updated", user)
}

// SendOTP is kept for route compatibility but phone OTP is now handled entirely
// by Firebase client-side (signInWithPhoneNumber). This endpoint is a no-op.
func SendOTP(c fiber.Ctx) error {
	return utils.SuccessResponse(c, 200, "Phone verification is handled via Firebase. Please use the app to authenticate.", nil)
}

func Login(c fiber.Ctx) error {
	var req models.LoginRequest

	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	var phoneNumber string
	var firebaseUID string

	if req.FirebaseToken != "" {
		// Verify Firebase ID token issued by Firebase phone authentication
		if config.FirebaseAuth == nil {
			return utils.ErrorResponse(c, 500, "Authentication service unavailable")
		}
		fbToken, err := config.FirebaseAuth.VerifyIDToken(c.Context(), req.FirebaseToken)
		if err != nil {
			fmt.Printf(" Firebase token verification failed: %v\n", err)
			return utils.ErrorResponse(c, 401, "OTP verification failed. Please request a new OTP and try again.")
		}
		// Firebase phone auth stores E.164 phone in the token claims
		fbPhone, _ := fbToken.Claims["phone_number"].(string)
		fbPhone = strings.TrimPrefix(fbPhone, "+91") // strip India country code
		fbPhone = strings.TrimPrefix(fbPhone, "+")
		if len(fbPhone) == 10 {
			phoneNumber = fbPhone
		} else if len(req.Phone) == 10 {
			phoneNumber = req.Phone
		} else {
			return utils.ErrorResponse(c, 400, "Could not resolve phone number from verification token")
		}
		firebaseUID = fbToken.UID
		fmt.Printf(" Firebase token verified for phone: %s\n", phoneNumber)
	} else {
		// No Firebase token — only allowed in development mode (for local testing)
		if os.Getenv("ENV") == "production" {
			return utils.ErrorResponse(c, 401, "Phone verification required. Please verify your phone number via OTP.")
		}
		if len(req.Phone) != 10 {
			return utils.ErrorResponse(c, 400, "Phone number must be exactly 10 digits")
		}
		phoneNumber = req.Phone
		fmt.Printf(" Login without Firebase token for phone: %s (dev mode)\n", phoneNumber)
	}

	// Read admin phones from env (comma-separated). Defaults to "0000000000" for local dev.
	adminPhones := os.Getenv("ADMIN_PHONES")
	if adminPhones == "" {
		adminPhones = "0000000000"
	}
	isAdmin := false
	for _, p := range strings.Split(adminPhones, ",") {
		if strings.TrimSpace(p) == phoneNumber {
			isAdmin = true
			break
		}
	}

	if config.FirestoreClient == nil {
		fmt.Println(" Firestore is not initialized. Cannot proceed with login.")
		return utils.ErrorResponse(c, 500, "Database service is not available. Please contact support.")
	}

	existingUser, err := userRepo.FindByPhone(c.Context(), phoneNumber)
	if err != nil {
		fmt.Printf(" Firestore query error for phone %s: %v\n", phoneNumber, err)
		return utils.ErrorResponse(c, 500, fmt.Sprintf("Database query failed: %v", err))
	}

	var user *models.User
	if existingUser != nil {
		user = existingUser

		if (user.FirebaseUID == "" && firebaseUID != "") || user.IsAdmin != isAdmin {
			user.FirebaseUID = firebaseUID
			user.IsAdmin = isAdmin
			_ = userRepo.Update(c.Context(), user)
		}
	} else {

		userID := firebaseUID
		if userID == "" {
			userID = utils.GenerateUUIDv7()
		}

		seqID := utils.GetNextSeqID()

		user = &models.User{
			ID:          userID,
			FirebaseUID: firebaseUID,
			SeqID:       seqID,
			Phone:       phoneNumber,
			IsAdmin:     isAdmin,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := userRepo.Create(c.Context(), user); err != nil {
			fmt.Printf("Database error creating user: %v\n", err)
			return utils.ErrorResponse(c, 500, "Failed to create user")
		}
	}

	token := utils.GenerateToken(user.ID, isAdmin)

	cookie := new(fiber.Cookie)
	cookie.Name = "authToken"
	cookie.Value = token
	cookie.Expires = time.Now().Add(30 * 24 * time.Hour)
	cookie.HTTPOnly = true
	cookie.SameSite = "Lax"
	c.Cookie(cookie)

	return utils.SuccessResponse(c, 200, "Login successful", fiber.Map{
		"token": token,
		"user":  user,
		"firebase_info": fiber.Map{
			"uid":         user.ID,
			"phone":       user.Phone,
			"is_new_user": existingUser == nil,
		},
	})
}

func Logout(c fiber.Ctx) error {
	cookie := new(fiber.Cookie)
	cookie.Name = "authToken"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-1 * time.Hour)
	cookie.HTTPOnly = true
	cookie.SameSite = "Lax"
	c.Cookie(cookie)

	return utils.SuccessResponse(c, 200, "Logged out successfully", nil)
}

func GetProfile(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)

	user, err := userRepo.FindByID(c.Context(), userID)
	if err != nil || user == nil {
		return utils.ErrorResponse(c, 404, "User not found")
	}

	return utils.SuccessResponse(c, 200, "Profile fetched", user)
}

func UpdateProfile(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)
	var req models.UpdateUserRequest

	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	user, err := userRepo.FindByID(c.Context(), userID)
	if err != nil || user == nil {
		return utils.ErrorResponse(c, 404, "User not found")
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		if user.Email != req.Email {
			user.IsEmailVerified = false
		}
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Address != "" {
		user.Address = req.Address
	}
	if req.State != "" {
		user.State = req.State
	}
	if req.District != "" {
		user.District = req.District
	}
	if req.Country != "" {
		user.Country = req.Country
	}
	user.UpdatedAt = time.Now()

	if err := userRepo.Update(c.Context(), user); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to update profile")
	}

	return utils.SuccessResponse(c, 200, "Profile updated successfully", user)
}
