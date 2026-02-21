package dining

import (
	"backend/controllers/dining"
	"backend/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterDiningRoutes(router fiber.Router) {
	diningGroup := router.Group("/dining")
	diningGroup.Get("/", middleware.OptionalAuth, dining.GetAllDiningVenues)
	diningGroup.Get("/:slug", dining.GetDiningVenueBySlug)
	diningGroup.Get("/id/:id", dining.GetDiningVenueByID)

	// Direct access to match frontend expectations
	diningGroup.Post("/", middleware.Auth, middleware.OrganizerOnly("dining"), dining.CreateDiningVenue)
	diningGroup.Put("/:id", middleware.Auth, middleware.OrganizerOnly("dining"), dining.UpdateDiningVenue)
	diningGroup.Delete("/:id", middleware.Auth, middleware.OrganizerOnly("dining"), dining.DeleteDiningVenue)

	organizerGroup := diningGroup.Group("/organizer", middleware.Auth, middleware.OrganizerOnly("dining"))
	organizerGroup.Get("/my", dining.GetOrganizerDiningVenues)
	organizerGroup.Post("/", dining.CreateDiningVenue)
	organizerGroup.Put("/:id", dining.UpdateDiningVenue)
	organizerGroup.Delete("/:id", dining.DeleteDiningVenue)
	organizerGroup.Patch("/:id/resubmit", dining.ResubmitDiningVenue)

	diningGroup.Post("/seed", dining.SeedDiningVenues)
}
