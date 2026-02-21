package offer

import (
	"backend/controllers/offer"
	"backend/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterOfferRoutes(router fiber.Router) {
	offersGroup := router.Group("/offers")
	offersGroup.Get("/", offer.GetAllOffers)
	offersGroup.Get("/user/:userId", offer.GetUserOffers)
	offersGroup.Get("/:id", offer.GetOfferByID)

	// Admin-only routes
	offersGroup.Post("/", middleware.Auth, middleware.AdminOnly, offer.CreateOffer)
	offersGroup.Put("/:id", middleware.Auth, middleware.AdminOnly, offer.UpdateOffer)
	offersGroup.Delete("/:id", middleware.Auth, middleware.AdminOnly, offer.DeleteOffer)
	offersGroup.Post("/seed", offer.SeedOffers)
}
