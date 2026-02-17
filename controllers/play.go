package controllers

import (
	"strconv"

	"backend/models"
	"backend/repository"
	"backend/utils"

	"github.com/gofiber/fiber/v3"
)

var playRepo = repository.NewPlayRepository()

func GetAllPlayVenues(c fiber.Ctx) error {
	limitStr := c.Query("limit", "20")
	limit, _ := strconv.Atoi(limitStr)
	cursor := c.Query("cursor", "")
	category := c.Query("category", "")
	city := c.Query("city", "")
	searchQuery := c.Query("q", "")

	venues, nextCursor, err := playRepo.GetPaginated(c.Context(), limit, cursor, category, city, searchQuery)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to fetch play venues")
	}

	return utils.SuccessResponse(c, 200, "Play venues fetched successfully", fiber.Map{
		"items":  venues,
		"cursor": nextCursor,
	})
}

func GetOrganizerPlayVenues(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)

	// Strict access check
	user, err := userRepo.FindByID(c.Context(), userID)
	if err != nil || user == nil || user.OrganizerCategory != "play" {
		// If category is empty, check if they have a pending/approved profile
		profile, _ := partnerRepo.FindByUserID(c.Context(), userID)
		if profile != nil && profile.OrganizationDetails.Category != "play" {
			return utils.ErrorResponse(c, 403, "Access denied: Incorrect organizer category")
		}
	}

	limitStr := c.Query("limit", "20")
	limit, _ := strconv.Atoi(limitStr)
	cursor := c.Query("cursor", "")

	venues, nextCursor, err := playRepo.FindPaginatedByOrganizerID(c.Context(), userID, limit, cursor)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to fetch venues")
	}

	return utils.SuccessResponse(c, 200, "Venues fetched successfully", fiber.Map{
		"items":  venues,
		"cursor": nextCursor,
	})
}

func CreatePlayVenue(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)

	// Enterprise check
	user, _ := userRepo.FindByID(c.Context(), userID)
	if user == nil || !user.IsOrganizer || user.OrganizerCategory != "play" {
		return utils.ErrorResponse(c, 403, "Only approved play organizers can create venues")
	}

	var venue models.PlayVenue
	if err := c.Bind().Body(&venue); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	venue.ID = utils.GenerateUUIDv7()
	venue.OrganizerID = userID
	venue.Status = "active"

	if err := playRepo.Create(c.Context(), &venue); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to create play venue")
	}

	return utils.SuccessResponse(c, 201, "Play venue created successfully", venue)
}

func GetPlayVenueBySlug(c fiber.Ctx) error {
	slug := c.Params("slug")
	venue, err := playRepo.GetBySlug(c.Context(), slug)
	if err != nil {
		return utils.ErrorResponse(c, 404, "Play venue not found")
	}
	return utils.SuccessResponse(c, 200, "Play venue fetched successfully", venue)
}

func GetPlayVenueByID(c fiber.Ctx) error {
	id := c.Params("id")
	venue, err := playRepo.FindByID(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, 404, "Play venue not found")
	}
	return utils.SuccessResponse(c, 200, "Play venue fetched successfully", venue)
}

func SeedPlayVenues(c fiber.Ctx) error {
	venue := models.PlayVenue{
		Name:                   "Turf Arena Chennai",
		Slug:                   "turf-arena-chennai",
		Status:                 "active",
		About:                  "Experience the finest sports infrastructure at our venue. Whether you're a professional athlete or a weekend player, our turf provides a high-quality playing surface with excellent lighting and facilities.",
		ShortAbout:             "Premium sports turf for football, cricket and more.",
		DurationPerSlotMinutes: 60,
		Location: models.PlayLocation{
			VenueName: "Turf Arena",
			Address:   "OMR Road, Thoraipakkam, Chennai, Tamil Nadu",
			City:      "Chennai",
			State:     "Tamil Nadu",
			Latitude:  12.9425,
			Longitude: 80.2363,
			MapURL:    "https://maps.google.com/?q=12.9425,80.2363",
		},
		Images: models.PlayImages{
			Hero:  "https://cdn.yoursite.com/play/turf-arena/hero.jpg",
			Promo: "https://cdn.yoursite.com/play/turf-arena/promo.jpg",
			Gallery: []string{
				"https://cdn.yoursite.com/play/turf-arena/gallery1.jpg",
				"https://cdn.yoursite.com/play/turf-arena/gallery2.jpg",
				"https://cdn.yoursite.com/play/turf-arena/gallery3.jpg",
				"https://cdn.yoursite.com/play/turf-arena/gallery4.jpg",
			},
		},
		PlayOptions: []models.PlayOption{
			{
				Sport:        "Football",
				CourtType:    "5v5",
				Surface:      "Artificial Grass",
				PricePerSlot: 1500,
			},
			{
				Sport:        "Cricket",
				CourtType:    "Box Cricket",
				Surface:      "Matting",
				PricePerSlot: 1200,
			},
			{
				Sport:        "Badminton",
				CourtType:    "Indoor Court",
				Surface:      "Wooden",
				PricePerSlot: 500,
			},
		},
		SlotSettings: models.SlotSettings{
			SlotDurationMinutes:   60,
			OpenTime:              "06:00",
			CloseTime:             "23:00",
			MaxDaysAdvanceBooking: 7,
		},
		FAQs: []models.PlayFAQ{
			{
				Question: "Are shoes available for rent?",
				Answer:   "No, players must bring their own sports shoes.",
			},
			{
				Question: "Is parking available?",
				Answer:   "Yes, free parking is available at the venue.",
			},
		},
		TermsAndConditions: []string{
			"Booking once confirmed cannot be cancelled or refunded.",
			"Players must arrive 10 minutes before slot time.",
			"Damage to property will be charged.",
			"Outside food is not allowed inside the turf.",
		},
	}

	if err := playRepo.Create(c.Context(), &venue); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to seed play venue")
	}

	return utils.SuccessResponse(c, 201, "Play venue seeded successfully", venue)
}

// AdminCreatePlayVenue - Admin can create play venues without organizer checks
func AdminCreatePlayVenue(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)

	var venue models.PlayVenue
	if err := c.Bind().Body(&venue); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	venue.ID = utils.GenerateUUIDv7()
	venue.OrganizerID = userID
	venue.Status = "active"

	if err := playRepo.Create(c.Context(), &venue); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to create play venue")
	}

	return utils.SuccessResponse(c, 201, "Play venue created successfully (admin)", venue)
}

// AdminUpdatePlayVenue - Admin can update play venues
func AdminUpdatePlayVenue(c fiber.Ctx) error {
	id := c.Params("id")
	var venue models.PlayVenue
	if err := c.Bind().Body(&venue); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	existing, err := playRepo.FindByID(c.Context(), id)
	if err != nil || existing == nil {
		return utils.ErrorResponse(c, 404, "Play venue not found")
	}

	venue.ID = id
	venue.OrganizerID = existing.OrganizerID
	venue.CreatedAt = existing.CreatedAt

	if err := playRepo.Update(c.Context(), &venue); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to update play venue")
	}

	return utils.SuccessResponse(c, 200, "Play venue updated successfully (admin)", venue)
}
