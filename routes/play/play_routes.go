package play

import (
	"backend/controllers/play"
	"backend/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterPlayRoutes(router fiber.Router) {
	playGroup := router.Group("/play")
	playGroup.Get("/", middleware.OptionalAuth, play.GetAllPlayVenues)
	playGroup.Get("/:slug", play.GetPlayVenueBySlug)
	playGroup.Get("/id/:id", play.GetPlayVenueByID)

	// Direct access to match frontend expectations
	playGroup.Post("/", middleware.Auth, middleware.OrganizerOnly("play"), play.CreatePlayVenue)
	playGroup.Put("/:id", middleware.Auth, middleware.OrganizerOnly("play"), play.UpdatePlayVenue)
	playGroup.Delete("/:id", middleware.Auth, middleware.OrganizerOnly("play"), play.DeletePlayVenue)

	organizerGroup := playGroup.Group("/organizer", middleware.Auth, middleware.OrganizerOnly("play"))
	organizerGroup.Get("/my", play.GetOrganizerPlayVenues)
	organizerGroup.Post("/", play.CreatePlayVenue)
	organizerGroup.Put("/:id", play.UpdatePlayVenue)
	organizerGroup.Delete("/:id", play.DeletePlayVenue)
	organizerGroup.Patch("/:id/resubmit", play.ResubmitPlayVenue)

	playGroup.Post("/seed", play.SeedPlayVenues)
}
