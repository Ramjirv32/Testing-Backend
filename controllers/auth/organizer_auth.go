package auth

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
)

func OrganizerRegister(c fiber.Ctx) error {
	var req models.OrganizerRegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	existing, _ := userRepo.FindByEmail(c.Context(), req.Email)
	if existing != nil {
		return utils.ErrorResponse(c, 400, "This email is already registered. Please try logging in instead.")
	}

	if req.Phone != "" {
		existingPhone, _ := userRepo.FindByPhone(c.Context(), req.Phone)
		if existingPhone != nil {
			return utils.ErrorResponse(c, 400, "This phone number is already registered. Please try logging in instead.")
		}
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
		IsOrganizer:     true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := userRepo.Create(c.Context(), user); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to create account")
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

	err = utils.SendOTPEmail(utils.EmailPlay, req.Email, "Verify Your Organizer Account - TicPin", body, otp)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to send OTP email")
	}

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

	token := utils.GenerateToken(user.ID, false)

	cookie := new(fiber.Cookie)
	cookie.Name = "authToken"
	cookie.Value = token
	cookie.Expires = time.Now().Add(30 * 24 * time.Hour)
	cookie.HTTPOnly = true
	cookie.SameSite = "Lax"
	c.Cookie(cookie)

	categories := user.OrganizerCategories
	if len(categories) == 0 && user.OrganizerCategory != "" {
		categories = []string{user.OrganizerCategory}
	}
	redirectTo := ""
	if user.IsOrganizer && user.IsPanVerified && len(categories) > 0 {
		redirectTo = "/organizer/" + categories[0] + "/dashboard"
	} else if user.IsOrganizer && user.IsPanVerified && user.OrganizerCategory != "" {
		redirectTo = "/organizer/" + user.OrganizerCategory + "/dashboard"
	}

	return utils.SuccessResponse(c, 200, "Login successful", fiber.Map{
		"token": token,
		"user":  user,
		"firebase_info": fiber.Map{
			"uid":   user.ID,
			"email": user.Email,
		},
		"is_pan_verified":      user.IsPanVerified,
		"organizer_categories": categories,
		"organizer_category":   user.OrganizerCategory,
		"redirect_to":          redirectTo,
	})
}

func OrganizerGoogleLogin(c fiber.Ctx) error {
	var req models.GoogleLoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	var email, name string

	if req.IDToken != "" {
		token, err := config.FirebaseAuth.VerifyIDToken(c.Context(), req.IDToken)
		if err != nil {
			fmt.Printf(" Firebase Token Verification Failed: %v\n", err)
			return utils.ErrorResponse(c, 401, "Invalid Firebase token")
		}

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

		user = &models.User{
			ID:              utils.GenerateUUIDv7(),
			SeqID:           utils.GetNextSeqID(),
			Email:           email,
			Name:            name,
			IsEmailVerified: true,
			IsOrganizer:     false,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		if err := userRepo.Create(c.Context(), user); err != nil {
			fmt.Printf(" OrganizerGoogleLogin: Failed to create user in Firestore: %v\n", err)
			return utils.ErrorResponse(c, 500, "Failed to create account")
		}
	}

	token := utils.GenerateToken(user.ID, false)

	cookie := new(fiber.Cookie)
	cookie.Name = "authToken"
	cookie.Value = token
	cookie.Expires = time.Now().Add(30 * 24 * time.Hour)
	cookie.HTTPOnly = true
	cookie.SameSite = "Lax"
	c.Cookie(cookie)

	gCategories := user.OrganizerCategories
	if len(gCategories) == 0 && user.OrganizerCategory != "" {
		gCategories = []string{user.OrganizerCategory}
	}
	gRedirectTo := ""
	if user.IsOrganizer && user.IsPanVerified && len(gCategories) > 0 {
		gRedirectTo = "/organizer/" + gCategories[0] + "/dashboard"
	} else if user.IsOrganizer && user.IsPanVerified && user.OrganizerCategory != "" {
		gRedirectTo = "/organizer/" + user.OrganizerCategory + "/dashboard"
	}

	return utils.SuccessResponse(c, 200, "Login successful", fiber.Map{
		"token": token,
		"user":  user,
		"firebase_info": fiber.Map{
			"uid":    user.ID,
			"email":  user.Email,
			"is_new": newUser,
		},
		"is_pan_verified":      user.IsPanVerified,
		"organizer_categories": gCategories,
		"organizer_category":   user.OrganizerCategory,
		"redirect_to":          gRedirectTo,
	})
}

func OrganizerVerifyOTP(c fiber.Ctx) error {
	var req models.VerifyEmailRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	emailOTPsMu.Lock()
	storedOTP, ok := emailOTPs[req.Email]
	emailOTPsMu.Unlock()
	if !ok {
		return utils.ErrorResponse(c, 401, "No verification request found for this email. Please request a new code.")
	}

	if time.Now().After(storedOTP.ExpiresAt) {
		emailOTPsMu.Lock()
		delete(emailOTPs, req.Email)
		emailOTPsMu.Unlock()
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

	emailOTPsMu.Lock()
	delete(emailOTPs, req.Email)
	emailOTPsMu.Unlock()

	token := utils.GenerateToken(user.ID, false)

	vCategories := user.OrganizerCategories
	if len(vCategories) == 0 && user.OrganizerCategory != "" {
		vCategories = []string{user.OrganizerCategory}
	}
	vRedirectTo := ""
	if user.IsOrganizer && user.IsPanVerified && len(vCategories) > 0 {
		vRedirectTo = "/organizer/" + vCategories[0] + "/dashboard"
	} else if user.IsOrganizer && user.IsPanVerified && user.OrganizerCategory != "" {
		vRedirectTo = "/organizer/" + user.OrganizerCategory + "/dashboard"
	}

	return utils.SuccessResponse(c, 200, "Verify successful", fiber.Map{
		"token":                token,
		"user":                 user,
		"is_pan_verified":      user.IsPanVerified,
		"organizer_categories": vCategories,
		"organizer_category":   user.OrganizerCategory,
		"redirect_to":          vRedirectTo,
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
	err = utils.SendOTPEmail(utils.EmailPlay, req.Email, "Resend: Verify Your Organizer Account - TicPin", body, otp)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to send email")
	}

	return utils.SuccessResponse(c, 200, "OTP resent successfully", nil)
}

func OrganizerForgotPassword(c fiber.Ctx) error {
	var req models.ForgotPasswordRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	user, _ := userRepo.FindByEmail(c.Context(), req.Email)
	if user == nil {
		return utils.ErrorResponse(c, 404, "Email does not exist. Please check your email or register.")
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
	err = utils.SendOTPEmail(utils.EmailPlay, req.Email, "Reset Your Organizer Password - TicPin", body, otp)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to send reset email")
	}

	return utils.SuccessResponse(c, 200, "Verification code sent to your email.", nil)
}

func OrganizerResetPassword(c fiber.Ctx) error {
	var req models.ResetPasswordRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	emailOTPsMu.Lock()
	storedOTP, ok := emailOTPs[req.Email]
	emailOTPsMu.Unlock()
	if !ok {
		return utils.ErrorResponse(c, 401, "No reset request found for this email. Please request a new code.")
	}

	if time.Now().After(storedOTP.ExpiresAt) {
		emailOTPsMu.Lock()
		delete(emailOTPs, req.Email)
		emailOTPsMu.Unlock()
		return utils.ErrorResponse(c, 401, "The reset code has expired. Please request a new one.")
	}

	if storedOTP.OTP != req.OTP {
		return utils.ErrorResponse(c, 401, "The reset code you entered is incorrect.")
	}

	user, _ := userRepo.FindByEmail(c.Context(), req.Email)
	if user == nil {
		return utils.ErrorResponse(c, 404, "User not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to hash password")
	}

	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	if err := userRepo.Update(c.Context(), user); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to reset password")
	}

	emailOTPsMu.Lock()
	delete(emailOTPs, req.Email)
	emailOTPsMu.Unlock()

	return utils.SuccessResponse(c, 200, "Password reset successfully. You can now login with your new password.", nil)
}
