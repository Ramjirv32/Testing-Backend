package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"fmt"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
)

func OrganizerRegister(c fiber.Ctx) error {
	var req models.OrganizerRegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	// Check if already exists
	existing, _ := userRepo.FindByEmail(c.Context(), req.Email)
	if existing != nil {
		return utils.ErrorResponse(c, 400, "This email is already registered. Please try logging in instead.")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to hash password")
	}

	user := &models.User{
		ID:              utils.GenerateUUIDv7(),
		SeqID:           utils.GetNextSeqID(),
		Email:           req.Email,
		Password:        string(hashedPassword),
		Name:            req.Name,
		Phone:           req.Phone,
		IsEmailVerified: false,
		IsOrganizer:     false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := userRepo.Create(c.Context(), user); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to create account")
	}

	// Send OTP
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	emailOTPs[req.Email] = OTPData{
		OTP:       otp,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	body := utils.GetOTPEmailTemplate(otp)

	_ = utils.SendEmail(utils.EmailAdmin, req.Email, "Verify Your Organizer Account - TicPin", body)

	return utils.SuccessResponse(c, 201, "Registration successful. Please verify OTP sent to your email.", nil)
}

func OrganizerLogin(c fiber.Ctx) error {
	var req models.OrganizerLoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	user, _ := userRepo.FindByEmail(c.Context(), req.Email)
	if user == nil {
		return utils.ErrorResponse(c, 401, "No account was found with this email. Please register first.")
	}

	if !user.IsEmailVerified {
		return utils.ErrorResponse(c, 403, "Email not verified. Please verify your account.")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return utils.ErrorResponse(c, 401, "The password you entered is incorrect. Please try again.")
	}

	token := utils.GenerateToken(user.ID, false) // Organizers are not admins by default

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
			"uid":   user.ID,
			"email": user.Email,
		},
	})
}

func OrganizerGoogleLogin(c fiber.Ctx) error {
	var req models.GoogleLoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	var email, name string

	// Verify Firebase Token
	if req.IDToken != "" {
		token, err := config.FirebaseAuth.VerifyIDToken(c.Context(), req.IDToken)
		if err != nil {
			return utils.ErrorResponse(c, 401, "Invalid Firebase token")
		}

		// Extract email and name from token claims
		emailRaw, ok := token.Claims["email"]
		if !ok || emailRaw == nil {
			return utils.ErrorResponse(c, 401, "Token does not contain email")
		}
		email = emailRaw.(string)

		nameRaw, ok := token.Claims["name"]
		if ok && nameRaw != nil {
			name = nameRaw.(string)
		} else {
			name = req.Name
		}
	} else {
		// Mock/Fallback (only for development if needed, but better to enforce token)
		if req.Email == "" {
			return utils.ErrorResponse(c, 400, "Email or Token is required")
		}
		email = req.Email
		name = req.Name
	}

	existingUser, _ := userRepo.FindByEmail(c.Context(), email)
	newUser := existingUser == nil
	user := existingUser

	if newUser {
		// New Google user, create profile
		user = &models.User{
			ID:              utils.GenerateUUIDv7(),
			SeqID:           utils.GetNextSeqID(),
			Email:           email,
			Name:            name,
			IsEmailVerified: true, // Google email is already verified
			IsOrganizer:     false,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		if err := userRepo.Create(c.Context(), user); err != nil {
			return utils.ErrorResponse(c, 500, "Failed to create account")
		}
	}

	token := utils.GenerateToken(user.ID, false)

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
			"uid":    user.ID,
			"email":  user.Email,
			"is_new": newUser,
		},
	})
}

func OrganizerVerifyOTP(c fiber.Ctx) error {
	var req models.VerifyEmailRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	storedOTP, ok := emailOTPs[req.Email]
	if !ok {
		return utils.ErrorResponse(c, 401, "No verification request found for this email. Please request a new code.")
	}

	if time.Now().After(storedOTP.ExpiresAt) {
		delete(emailOTPs, req.Email)
		return utils.ErrorResponse(c, 401, "The verification code has expired. Please request a new one.")
	}

	if storedOTP.OTP != req.OTP {
		return utils.ErrorResponse(c, 401, "The verification code you entered is incorrect. Please check and try again.")
	}

	user, _ := userRepo.FindByEmail(c.Context(), req.Email)
	if user == nil {
		return utils.ErrorResponse(c, 404, "User not found")
	}

	user.IsEmailVerified = true
	user.UpdatedAt = time.Now()

	if err := userRepo.Update(c.Context(), user); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to update verify email")
	}

	delete(emailOTPs, req.Email)

	token := utils.GenerateToken(user.ID, false)

	return utils.SuccessResponse(c, 200, "Verify successful", fiber.Map{
		"token": token,
		"user":  user,
	})
}

func ResendOrganizerOTP(c fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	if req.Email == "" {
		return utils.ErrorResponse(c, 400, "Email is required")
	}

	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	emailOTPs[req.Email] = OTPData{
		OTP:       otp,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	body := utils.GetOTPEmailTemplate(otp)
	err := utils.SendEmail(utils.EmailAdmin, req.Email, "Resend: Verify Your Organizer Account - TicPin", body)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to send email")
	}

	return utils.SuccessResponse(c, 200, "OTP resent successfully", nil)
}
