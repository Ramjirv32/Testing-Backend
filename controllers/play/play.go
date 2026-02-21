package play

import (
	"strconv"

	"backend/models"
	"backend/repository"
	"backend/utils"

	"github.com/gofiber/fiber/v3"
)

var playRepo = repository.NewPlayRepository()
var userRepo = repository.NewUserRepository()
var partnerRepo = repository.NewPartnerRepository()

func GetAllPlayVenues(c fiber.Ctx) error {
	limitStr := c.Query("limit", "20")
	limit, _ := strconv.Atoi(limitStr)
	cursor := c.Query("cursor", "")
	category := c.Query("category", "")
	city := c.Query("city", "")
	searchQuery := c.Query("q", "")
	statusParam := c.Query("status", "")

	isAdmin, _ := c.Locals("isAdmin").(bool)

	status := "active"
	showAll := c.Query("all", "") == "true"

	// If a specific status is requested and user is admin, use that
	if statusParam != "" && isAdmin {
		status = statusParam
	} else if showAll && isAdmin {
		status = "" // Show all statuses
	}

	venues, nextCursor, err := playRepo.GetPaginated(c.Context(), limit, cursor, category, city, searchQuery, status)
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

	user, err := userRepo.FindByID(c.Context(), userID)
	if err != nil || user == nil || !user.HasOrganizerCategory("play") {
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
		// Return empty list on error (e.g. no venues created yet)
		return utils.SuccessResponse(c, 200, "Venues fetched successfully", fiber.Map{
			"items":  []*models.PlayVenue{},
			"cursor": "",
		})
	}

	return utils.SuccessResponse(c, 200, "Venues fetched successfully", fiber.Map{
		"items":  venues,
		"cursor": nextCursor,
	})
}

// extractSports returns a deduplicated list of sports from play_options
func extractSports(options []models.PlayOption) []string {
	seen := map[string]bool{}
	var sports []string
	for _, opt := range options {
		if opt.Sport != "" && !seen[opt.Sport] {
			seen[opt.Sport] = true
			sports = append(sports, opt.Sport)
		}
	}
	return sports
}

func CreatePlayVenue(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)

	user, _ := userRepo.FindByID(c.Context(), userID)
	if user == nil || !user.IsOrganizer || !user.HasOrganizerCategory("play") {
		return utils.ErrorResponse(c, 403, "Only approved play organizers can create venues")
	}

	var venue models.PlayVenue
	if err := c.Bind().Body(&venue); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	venue.ID = utils.GenerateUUIDv7()
	venue.OrganizerID = userID
	venue.Status = "pending" // Organizer submissions need admin approval
	venue.Sports = extractSports(venue.PlayOptions)

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
	venue.Sports = extractSports(venue.PlayOptions)

	if err := playRepo.Create(c.Context(), &venue); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to seed play venue")
	}

	return utils.SuccessResponse(c, 201, "Play venue seeded successfully", venue)
}

func AdminCreatePlayVenue(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)

	var venue models.PlayVenue
	if err := c.Bind().Body(&venue); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	venue.ID = utils.GenerateUUIDv7()
	venue.OrganizerID = userID
	venue.Status = "active"
	venue.Sports = extractSports(venue.PlayOptions)

	if err := playRepo.Create(c.Context(), &venue); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to create play venue")
	}

	return utils.SuccessResponse(c, 201, "Play venue created successfully (admin)", venue)
}

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
	venue.Sports = extractSports(venue.PlayOptions)

	if err := playRepo.Update(c.Context(), &venue); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to update play venue")
	}

	return utils.SuccessResponse(c, 200, "Play venue updated successfully (admin)", venue)
}

func UpdatePlayVenue(c fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("uid").(string)

	existing, err := playRepo.FindByID(c.Context(), id)
	if err != nil || existing == nil {
		return utils.ErrorResponse(c, 404, "Play venue not found")
	}

	if existing.OrganizerID != userID {
		return utils.ErrorResponse(c, 403, "Access denied: You do not own this venue")
	}

	var venue models.PlayVenue
	if err := c.Bind().Body(&venue); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	venue.ID = id
	venue.OrganizerID = userID
	venue.Status = existing.Status // Preserve status - organizers can't change it
	venue.CreatedAt = existing.CreatedAt
	venue.Sports = extractSports(venue.PlayOptions)

	if err := playRepo.Update(c.Context(), &venue); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to update play venue")
	}

	return utils.SuccessResponse(c, 200, "Play venue updated successfully", venue)
}

func AdminApprovePlayVenue(c fiber.Ctx) error {
	id := c.Params("id")
	status := c.Query("status", "active")
	reason := c.Query("reason", "")

	existing, err := playRepo.FindByID(c.Context(), id)
	if err != nil || existing == nil {
		return utils.ErrorResponse(c, 404, "Play venue not found")
	}

	existing.Status = status
	existing.RejectionReason = reason

	if err := playRepo.Update(c.Context(), existing); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to update play venue status")
	}

	return utils.SuccessResponse(c, 200, "Play venue status updated successfully", nil)
}

func DeletePlayVenue(c fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("uid").(string)

	existing, err := playRepo.FindByID(c.Context(), id)
	if err != nil || existing == nil {
		return utils.ErrorResponse(c, 404, "Play venue not found")
	}

	if existing.OrganizerID != userID {
		return utils.ErrorResponse(c, 403, "Access denied: You do not own this venue")
	}

	if err := playRepo.Delete(c.Context(), id); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to delete play venue")
	}

	return utils.SuccessResponse(c, 200, "Play venue deleted successfully", nil)
}

func AdminDeletePlayVenue(c fiber.Ctx) error {
	id := c.Params("id")

	existing, err := playRepo.FindByID(c.Context(), id)
	if err != nil || existing == nil {
		return utils.ErrorResponse(c, 404, "Play venue not found")
	}

	if err := playRepo.Delete(c.Context(), id); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to delete play venue")
	}

	return utils.SuccessResponse(c, 200, "Play venue deleted successfully (admin)", nil)
}

func ResubmitPlayVenue(c fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("uid").(string)

	existing, err := playRepo.FindByID(c.Context(), id)
	if err != nil || existing == nil {
		return utils.ErrorResponse(c, 404, "Play venue not found")
	}

	if existing.OrganizerID != userID {
		return utils.ErrorResponse(c, 403, "Access denied: You do not own this venue")
	}

	if existing.Status != "rejected" {
		return utils.ErrorResponse(c, 400, "Only rejected venues can be resubmitted")
	}

	existing.Status = "pending"
	existing.RejectionReason = ""

	if err := playRepo.Update(c.Context(), existing); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to resubmit play venue")
	}

	return utils.SuccessResponse(c, 200, "Play venue resubmitted for review", nil)
}
