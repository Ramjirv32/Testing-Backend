package controllers

import (
	"github.com/gofiber/fiber/v3"

	"backend/models"
	"backend/repository"
	"backend/utils"
)

var artistRepo = repository.NewArtistRepository()

func GetAllArtists(c fiber.Ctx) error {
	artists, err := artistRepo.GetAll(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to fetch artists")
	}

	return utils.SuccessResponse(c, 200, "Artists fetched successfully", artists)
}

func GetArtistByID(c fiber.Ctx) error {
	id := c.Params("id")
	artist, err := artistRepo.FindByID(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, 404, "Artist not found")
	}
	return utils.SuccessResponse(c, 200, "Artist fetched successfully", artist)
}

func CreateArtist(c fiber.Ctx) error {
	var artist models.Artist
	if err := c.Bind().Body(&artist); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	if err := artistRepo.Create(c.Context(), &artist); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to create artist")
	}

	return utils.SuccessResponse(c, 201, "Artist created successfully", artist)
}

func SeedArtists(c fiber.Ctx) error {
	artists := []models.Artist{
		{Name: "DJ Aurora", Role: "Main Artist", ImageURL: "/events/artists/1.png"},
		{Name: "The Waves", Role: "Band", ImageURL: "/events/artists/2.png"},
		{Name: "Luna Smith", Role: "Solo Singer", ImageURL: "/events/artists/3.png"},
		{Name: "Echo Band", Role: "Band", ImageURL: "/events/artists/4.png"},
		{Name: "Sarah Jazz", Role: "Jazz Singer", ImageURL: "/events/artists/5.png"},
		{Name: "Max Power", Role: "Rock Artist", ImageURL: "/events/artists/6.png"},
	}

	for _, a := range artists {
		_ = artistRepo.Create(c.Context(), &a)
	}

	return utils.SuccessResponse(c, 200, "Artists seeded successfully", nil)
}
