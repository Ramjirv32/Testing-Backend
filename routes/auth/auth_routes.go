package auth

import (
	"backend/controllers/auth"
	"backend/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterAuthRoutes(router fiber.Router) {
	
	orgAuth := router.Group("/organizer/auth")
	orgAuth.Post("/register", auth.OrganizerRegister)
	orgAuth.Post("/login", auth.OrganizerLogin)
	orgAuth.Post("/google", auth.OrganizerGoogleLogin)
	orgAuth.Post("/verify-otp", auth.OrganizerVerifyOTP)
	orgAuth.Post("/resend-otp", auth.ResendOrganizerOTP)
	orgAuth.Post("/forgot-password", auth.OrganizerForgotPassword)
	orgAuth.Post("/reset-password", auth.OrganizerResetPassword)

	userAuth := router.Group("/auth")
	userAuth.Post("/send-otp", auth.SendOTP)
	userAuth.Post("/login", auth.Login)
	userAuth.Post("/logout", auth.Logout)
	userAuth.Get("/profile", middleware.Auth, auth.GetProfile)
	userAuth.Put("/profile", middleware.Auth, auth.UpdateProfile)
	userAuth.Post("/email/send-otp", middleware.Auth, auth.SendEmailOTP)
	userAuth.Post("/email/verify", middleware.Auth, auth.VerifyEmail)
}
