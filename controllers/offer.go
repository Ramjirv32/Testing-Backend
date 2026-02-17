package controllers

import (
	"github.com/gofiber/fiber/v3"

	"backend/models"
	"backend/repository"
	"backend/utils"
)

var offerRepo = repository.NewOfferRepository()

func CreateOffer(c fiber.Ctx) error {
	var req models.Offer
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	if err := offerRepo.Create(c.Context(), &req); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to create offer")
	}

	return utils.SuccessResponse(c, 201, "Offer created successfully", req)
}

func GetUserOffers(c fiber.Ctx) error {
	userID := c.Params("userId")
	if userID == "" {
		return utils.ErrorResponse(c, 400, "User ID is required")
	}

	offers, err := offerRepo.GetByUserID(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to fetch offers")
	}

	return utils.SuccessResponse(c, 200, "Offers fetched successfully", offers)
}

func GetAllOffers(c fiber.Ctx) error {
	offers, err := offerRepo.GetAll(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to fetch offers")
	}

	return utils.SuccessResponse(c, 200, "Offers fetched successfully", offers)
}

func SeedOffers(c fiber.Ctx) error {
	offers := []models.Offer{
		{
			Code:        "WELCOME50",
			Discount:    "50% OFF",
			Description: "Welcome offer for new users",
			IsActive:    true,
			UserID:      "global", // Using 'global' as a special ID for all users
		},
		{
			Code:        "OFFER20",
			Discount:    "20% OFF",
			Description: "Special discount",
			IsActive:    true,
			UserID:      "global",
		},
	}

	for _, o := range offers {
		_ = offerRepo.Create(c.Context(), &o)
	}

	return utils.SuccessResponse(c, 200, "Offers seeded successfully", nil)
}
