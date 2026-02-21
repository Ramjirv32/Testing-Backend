package event

import (
	"backend/controllers/event"
	"backend/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterEventRoutes(router fiber.Router) {
	eventsGroup := router.Group("/events")
	eventsGroup.Get("/", middleware.OptionalAuth, event.GetAllEvents)
	eventsGroup.Get("/:id", event.GetEventByID)

	// Direct access to match frontend expectations
	eventsGroup.Post("/", middleware.Auth, middleware.OrganizerOnly("event"), event.CreateEvent)
	eventsGroup.Put("/:id", middleware.Auth, middleware.OrganizerOnly("event"), event.UpdateEvent)
	eventsGroup.Delete("/:id", middleware.Auth, middleware.OrganizerOnly("event"), event.DeleteEvent)

	organizerGroup := eventsGroup.Group("/organizer", middleware.Auth, middleware.OrganizerOnly("event"))
	organizerGroup.Get("/my", event.GetOrganizerEvents)
	organizerGroup.Post("/", event.CreateEvent)
	organizerGroup.Put("/:id", event.UpdateEvent)
	organizerGroup.Delete("/:id", event.DeleteEvent)
	organizerGroup.Patch("/:id/resubmit", event.ResubmitEvent)

	artistsGroup := router.Group("/artists")
	artistsGroup.Get("/", event.GetAllArtists)
	artistsGroup.Get("/:id", event.GetArtistByID)
	artistsGroup.Post("/", middleware.Auth, middleware.OrganizerOnly("event"), event.CreateArtist)
	artistsGroup.Post("/seed", event.SeedArtists)
}
