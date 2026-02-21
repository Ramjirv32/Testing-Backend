package partner

import (
	"backend/controllers/partner"
	"backend/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterPartnerRoutes(router fiber.Router) {
	partnersGroup := router.Group("/partners", middleware.Auth)
	partnersGroup.Post("/verify", partner.SubmitVerification)
	partnersGroup.Post("/verify-pan", partner.VerifyPAN)
	partnersGroup.Post("/pan-gstin", partner.GetGSTINsFromPAN)
	partnersGroup.Get("/my-status", partner.GetMyVerificationStatus)
	partnersGroup.Get("/prefill", partner.GetPrefillData)

	// Admin routes
	adminGroup := router.Group("/admin/partners", middleware.Auth, middleware.AdminOnly)
	adminGroup.Get("/", partner.GetAllPartners)
}
