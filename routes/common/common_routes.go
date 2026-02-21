package common

import (
	"backend/controllers/common"
	"backend/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterCommonRoutes(router fiber.Router) {
	router.Post("/upload", middleware.Auth, common.UploadFile)

	// Email endpoints
	email := router.Group("/emails")
	email.Post("/pass-purchase", common.SendPassPurchaseEmail)
	email.Post("/pass-renewal", common.SendPassRenewalEmail)
	email.Post("/pass-expiry-reminder", common.SendPassExpiryReminderEmail)
	email.Post("/booking-confirmation", common.SendBookingConfirmationEmail)

	// Pass endpoints
	pass := router.Group("/pass")
	pass.Post("/check-eligibility", common.CheckPassEligibility)
	pass.Post("/check-eligibility-by-phone", common.CheckPassEligibility) // phone only variant (same handler)

	health := router.Group("/health")
	health.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// State and district data (public)
	router.Get("/states", common.GetAllStates)
	router.Get("/states/:state/districts", common.GetDistrictsByState)
	router.Get("/locations/search", common.SearchLocations)
}
