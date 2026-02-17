package controllers

import (
	"strconv"

	"backend/models"
	"backend/repository"
	"backend/utils"

	"github.com/gofiber/fiber/v3"
)

var diningRepo = repository.NewDiningRepository()

func GetAllDiningVenues(c fiber.Ctx) error {
	limitStr := c.Query("limit", "20")
	limit, _ := strconv.Atoi(limitStr)
	cursor := c.Query("cursor", "")
	category := c.Query("category", "")
	city := c.Query("city", "")
	searchQuery := c.Query("q", "")

	venues, nextCursor, err := diningRepo.GetPaginated(c.Context(), limit, cursor, category, city, searchQuery)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to fetch dining venues")
	}

	return utils.SuccessResponse(c, 200, "Dining venues fetched successfully", fiber.Map{
		"items":  venues,
		"cursor": nextCursor,
	})
}

func GetOrganizerDiningVenues(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)

	// Strict access check
	user, err := userRepo.FindByID(c.Context(), userID)
	if err != nil || user == nil || user.OrganizerCategory != "dining" {
		// If category is empty, check if they have a pending/approved profile
		profile, _ := partnerRepo.FindByUserID(c.Context(), userID)
		if profile != nil && profile.OrganizationDetails.Category != "dining" {
			return utils.ErrorResponse(c, 403, "Access denied: Incorrect organizer category")
		}
	}

	limitStr := c.Query("limit", "20")
	limit, _ := strconv.Atoi(limitStr)
	cursor := c.Query("cursor", "")

	venues, nextCursor, err := diningRepo.FindPaginatedByOrganizerID(c.Context(), userID, limit, cursor)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to fetch venues")
	}

	return utils.SuccessResponse(c, 200, "Venues fetched successfully", fiber.Map{
		"items":  venues,
		"cursor": nextCursor,
	})
}

func CreateDiningVenue(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)

	// Enterprise check
	user, _ := userRepo.FindByID(c.Context(), userID)
	if user == nil || !user.IsOrganizer || user.OrganizerCategory != "dining" {
		return utils.ErrorResponse(c, 403, "Only approved dining organizers can create outlets")
	}

	var venue models.DiningVenue
	if err := c.Bind().Body(&venue); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	venue.ID = utils.GenerateUUIDv7()
	venue.OrganizerID = userID
	venue.Status = "active"

	if err := diningRepo.Create(c.Context(), &venue); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to create dining venue")
	}

	return utils.SuccessResponse(c, 201, "Dining outlet created successfully", venue)
}

func GetDiningVenueBySlug(c fiber.Ctx) error {
	slug := c.Params("slug")
	venue, err := diningRepo.GetBySlug(c.Context(), slug)
	if err != nil {
		return utils.ErrorResponse(c, 404, "Dining venue not found")
	}
	return utils.SuccessResponse(c, 200, "Dining venue fetched successfully", venue)
}

func GetDiningVenueByID(c fiber.Ctx) error {
	id := c.Params("id")
	venue, err := diningRepo.FindByID(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, 404, "Dining venue not found")
	}
	return utils.SuccessResponse(c, 200, "Dining venue fetched successfully", venue)
}

func SeedDiningVenues(c fiber.Ctx) error {
	venue := models.DiningVenue{
		Name:             "The Grand Buffet",
		Slug:             "the-grand-buffet",
		Status:           "active",
		Description:      "The Grand Buffet offers a premium multi-cuisine dining experience with a luxurious ambience, live counters, and a wide variety of vegetarian and non-vegetarian dishes. Perfect for family dinners, corporate meetings, and celebrations.",
		ShortDescription: "Premium multi-cuisine buffet with luxury ambience.",
		Rating:           4.8,
		IsOpen:           true,
		OpeningTime:      "11:00",
		ClosingTime:      "23:00",
		Contact: models.DiningContact{
			Phone: "+916383667872",
			Email: "info@grandbuffet.com",
		},
		Location: models.DiningLocation{
			VenueName: "The Grand Buffet",
			Address:   "Anna Nagar Main Road, Chennai, Tamil Nadu",
			City:      "Chennai",
			State:     "Tamil Nadu",
			Latitude:  13.0850,
			Longitude: 80.2101,
			MapURL:    "https://maps.google.com/?q=13.0850,80.2101",
		},
		Images: models.DiningImages{
			Hero: "https://cdn.yoursite.com/dining/grand-buffet/hero.jpg",
			Gallery: []string{
				"https://cdn.yoursite.com/dining/grand-buffet/gallery1.jpg",
				"https://cdn.yoursite.com/dining/grand-buffet/gallery2.jpg",
				"https://cdn.yoursite.com/dining/grand-buffet/gallery3.jpg",
				"https://cdn.yoursite.com/dining/grand-buffet/gallery4.jpg",
			},
		},
		MenuImages: []string{
			"https://cdn.yoursite.com/dining/grand-buffet/menu1.jpg",
			"https://cdn.yoursite.com/dining/grand-buffet/menu2.jpg",
			"https://cdn.yoursite.com/dining/grand-buffet/menu3.jpg",
			"https://cdn.yoursite.com/dining/grand-buffet/menu4.jpg",
		},
		Offers: []models.DiningOffer{
			{
				Title:       "Flat 30% Off",
				Code:        "DINE30",
				Description: "Get 30% off on buffet bookings.",
			},
			{
				Title:       "Buy 2 Get 1 Free",
				Code:        "GROUP3",
				Description: "Valid for group dining.",
			},
		},
		Facilities: []string{
			"Air Conditioning",
			"Parking Available",
			"Free WiFi",
			"Live Music",
			"Family Friendly",
			"Card Payment",
			"Wheelchair Access",
			"Outdoor Seating",
		},
		SeatingTypes: []models.SeatingType{
			{
				Type:             "Indoor",
				TotalTables:      40,
				AvailableTables:  18,
				CapacityPerTable: 4,
			},
			{
				Type:             "Outdoor",
				TotalTables:      15,
				AvailableTables:  6,
				CapacityPerTable: 4,
			},
			{
				Type:             "Private Dining",
				TotalTables:      5,
				AvailableTables:  2,
				CapacityPerTable: 8,
			},
		},
		BookingSettings: models.BookingSettings{
			AdvanceBookingDays: 7,
			TimeSlots: []string{
				"11:00",
				"12:00",
				"13:00",
				"14:00",
				"18:00",
				"19:00",
				"20:00",
				"21:00",
			},
			AverageDiningDurationMinutes: 90,
		},
		FAQs: []models.DiningFAQ{
			{
				Question: "Is buffet available daily?",
				Answer:   "Yes, buffet is available for lunch and dinner every day.",
			},
			{
				Question: "Do you allow birthday decorations?",
				Answer:   "Yes, prior booking is required for special arrangements.",
			},
		},
		TermsAndConditions: []string{
			"Table will be held for 15 minutes from booking time.",
			"Cancellation allowed up to 2 hours before slot.",
			"Outside food and drinks are not permitted.",
			"Management reserves the right of admission.",
		},
	}

	if err := diningRepo.Create(c.Context(), &venue); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to seed dining venue")
	}

	return utils.SuccessResponse(c, 219, "Dining venue seeded successfully", venue)
}

// AdminCreateDiningVenue - Admin can create dining venues without organizer checks
func AdminCreateDiningVenue(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)

	var venue models.DiningVenue
	if err := c.Bind().Body(&venue); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	venue.ID = utils.GenerateUUIDv7()
	venue.OrganizerID = userID
	venue.Status = "active"

	if err := diningRepo.Create(c.Context(), &venue); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to create dining venue")
	}

	return utils.SuccessResponse(c, 201, "Dining outlet created successfully (admin)", venue)
}

// AdminUpdateDiningVenue - Admin can update dining venues
func AdminUpdateDiningVenue(c fiber.Ctx) error {
	id := c.Params("id")
	var venue models.DiningVenue
	if err := c.Bind().Body(&venue); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	existing, err := diningRepo.FindByID(c.Context(), id)
	if err != nil || existing == nil {
		return utils.ErrorResponse(c, 404, "Dining venue not found")
	}

	venue.ID = id
	venue.OrganizerID = existing.OrganizerID
	venue.CreatedAt = existing.CreatedAt

	if err := diningRepo.Update(c.Context(), &venue); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to update dining venue")
	}

	return utils.SuccessResponse(c, 200, "Dining venue updated successfully (admin)", venue)
}
