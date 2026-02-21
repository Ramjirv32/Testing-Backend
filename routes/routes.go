package routes

import (
	"github.com/gofiber/fiber/v3"

	"backend/routes/admin"
	"backend/routes/ai"
	"backend/routes/auth"
	"backend/routes/booking"
	"backend/routes/common"
	"backend/routes/dining"
	"backend/routes/event"
	"backend/routes/offer"
	"backend/routes/partner"
	"backend/routes/pass"
	"backend/routes/play"
	"backend/routes/user"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	common.RegisterCommonRoutes(api)
	auth.RegisterAuthRoutes(api)
	booking.RegisterBookingRoutes(api)
	pass.RegisterPassRoutes(api)
	admin.RegisterAdminRoutes(api)
	partner.RegisterPartnerRoutes(api)
	event.RegisterEventRoutes(api)
	play.RegisterPlayRoutes(api)
	dining.RegisterDiningRoutes(api)
	user.RegisterUserRoutes(api)
	offer.RegisterOfferRoutes(api)
	ai.RegisterAIRoutes(api)
}
