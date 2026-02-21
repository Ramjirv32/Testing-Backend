package admin

import (
	"backend/controllers/booking"
	"backend/controllers/common"
	"backend/controllers/dining"
	"backend/controllers/event"
	"backend/controllers/partner"
	"backend/controllers/play"
	"backend/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterAdminRoutes(router fiber.Router) {
	adminGroup := router.Group("/admin", middleware.Auth, middleware.AdminOnly)

	// Secure document signed-URL endpoint (admin-only)
	adminGroup.Get("/signed-url", common.GetSignedURL)

	adminGroup.Delete("/play-bookings/all", booking.DeleteAllPlayBookings)
	adminGroup.Delete("/dining-bookings/all", booking.DeleteAllDiningBookings)

	adminGroup.Delete("/users/all", booking.DeleteAllUsers)

	adminGroup.Get("/partners", partner.GetEventPosters)
	adminGroup.Patch("/partners/:id/approve", partner.ApproveEventPoster)
	adminGroup.Put("/partners/:id", partner.UpdatePartnerProfile)

	adminGroup.Post("/play", play.AdminCreatePlayVenue)
	adminGroup.Put("/play/:id", play.AdminUpdatePlayVenue)
	adminGroup.Delete("/play/:id", play.AdminDeletePlayVenue)
	adminGroup.Post("/dining", dining.AdminCreateDiningVenue)
	adminGroup.Put("/dining/:id", dining.AdminUpdateDiningVenue)
	adminGroup.Delete("/dining/:id", dining.AdminDeleteDiningVenue)
	adminGroup.Post("/events", event.AdminCreateEvent)
	adminGroup.Put("/events/:id", event.AdminUpdateEvent)
	adminGroup.Delete("/events/:id", event.AdminDeleteEvent)
	adminGroup.Patch("/events/:id/approve", event.AdminApproveEvent)
	adminGroup.Patch("/play/:id/approve", play.AdminApprovePlayVenue)
	adminGroup.Patch("/dining/:id/approve", dining.AdminApproveDiningVenue)
}
