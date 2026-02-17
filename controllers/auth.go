package controllers

import (
	"backend/config"
	"backend/models"
	"backend/repository"
	"backend/utils"
	"fmt"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v3"
)

type OTPData struct {
	OTP       string
	ExpiresAt time.Time
}

var userRepo = repository.NewUserRepository()
var emailOTPs = make(map[string]OTPData)
var phoneOTPs = make(map[string]OTPData)

func SendEmailOTP(c fiber.Ctx) error {
	var req models.SendEmailOTPRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	emailOTPs[req.Email] = OTPData{
		OTP:       otp,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	body := utils.GetOTPEmailTemplate(otp)

	err := utils.SendEmail(utils.EmailAdmin, req.Email, "Email Verification - TicPin", body)
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

	storedOTP, ok := emailOTPs[req.Email]
	if !ok {
		return utils.ErrorResponse(c, 401, "No OTP sent for this email")
	}

	if time.Now().After(storedOTP.ExpiresAt) {
		delete(emailOTPs, req.Email)
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

	delete(emailOTPs, req.Email)

	return utils.SuccessResponse(c, 200, "Email verified and profile updated", user)
}

func SendOTP(c fiber.Ctx) error {
	var req models.SendOTPRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	if req.Phone == "" {
		return utils.ErrorResponse(c, 400, "Phone number is required")
	}

	// Generate 6-digit OTP
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	// For testing/admin number, use static OTP
	if req.Phone == "6383667872" {
		otp = "123456"
	}

	phoneOTPs[req.Phone] = OTPData{
		OTP:       otp,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	// Send OTP via SMS
	message := utils.GetOTPSMSMessage(otp)
	err := utils.SendSMS(req.Phone, message)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to send SMS")
	}

	return utils.SuccessResponse(c, 200, "OTP sent successfully", nil)
}

func Login(c fiber.Ctx) error {
	var req models.LoginRequest

	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	var phoneNumber string
	var firebaseUID string

	// Handle Firebase Token Verification
	if req.FirebaseToken != "" {
		if config.FirebaseAuth == nil {
			return utils.ErrorResponse(c, 500, "Firebase authentication is not initialized on the server. Please check server logs.")
		}
		token, err := config.FirebaseAuth.VerifyIDToken(c.Context(), req.FirebaseToken)
		if err != nil {
			fmt.Printf("Firebase token verification failed: %v\n", err)
			return utils.ErrorResponse(c, 401, "Invalid or expired Firebase token")
		}

		firebaseUID = token.UID

		// Extract phone number from token claims
		phoneRaw, ok := token.Claims["phone_number"]
		if !ok || phoneRaw == nil {
			return utils.ErrorResponse(c, 401, "Token does not contain phone number")
		}
		phoneNumber = phoneRaw.(string)

		// Sanitize phone number (Firebase returns +91...)
		if len(phoneNumber) > 10 {
			phoneNumber = phoneNumber[len(phoneNumber)-10:]
		}
	} else {
		// Fallback to custom OTP system
		if req.Phone == "" {
			return utils.ErrorResponse(c, 400, "Phone number is required")
		}

		if req.OTP == "" {
			return utils.ErrorResponse(c, 400, "OTP is required")
		}

		// Verify OTP
		storedOTP, ok := phoneOTPs[req.Phone]
		if !ok {
			// Fallback for demo: allow 123456 for now if no OTP sent,
			// but only if it's the admin number or if we want to be strict.
			if req.OTP != "123456" || (req.Phone != "6383667872" && req.Phone != "9100000000") {
				return utils.ErrorResponse(c, 401, "No verification request found for this phone number. Please request a new code.")
			}
		} else {
			if time.Now().After(storedOTP.ExpiresAt) {
				delete(phoneOTPs, req.Phone)
				return utils.ErrorResponse(c, 401, "The verification code has expired. Please request a new one.")
			}

			if storedOTP.OTP != req.OTP {
				return utils.ErrorResponse(c, 401, "The verification code you entered is incorrect. Please check and try again.")
			}

			// Success - remove the OTP
			delete(phoneOTPs, req.Phone)
		}
		phoneNumber = req.Phone
	}

	isAdmin := phoneNumber == "6383667872"

	existingUser, err := userRepo.FindByPhone(c.Context(), phoneNumber)
	if err != nil {
		fmt.Printf("Database error finding user by phone (+%s): %v\n", phoneNumber, err)
		return utils.ErrorResponse(c, 500, "Database error")
	}

	var user *models.User
	if existingUser != nil {
		user = existingUser
		// Update Firebase UID if missing
		if user.FirebaseUID == "" && firebaseUID != "" {
			user.FirebaseUID = firebaseUID
			_ = userRepo.Update(c.Context(), user)
		}
	} else {
		// Create new user using Firebase UID or internal ID
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
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := userRepo.Create(c.Context(), user); err != nil {
			fmt.Printf("Database error creating user: %v\n", err)
			return utils.ErrorResponse(c, 500, "Failed to create user")
		}
	}

	token := utils.GenerateToken(user.ID, isAdmin)

	// Set cookie
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
	user.UpdatedAt = time.Now()

	if err := userRepo.Update(c.Context(), user); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to update profile")
	}

	return utils.SuccessResponse(c, 200, "Profile updated successfully", user)
}
