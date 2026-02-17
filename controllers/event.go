package controllers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"

	"backend/models"
	"backend/repository"
	"backend/utils"
)

var eventRepo = repository.NewEventRepository()

func CreateEvent(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)

	// Check if user is an approved organizer
	user, err := userRepo.FindByID(c.Context(), userID)
	if err != nil || user == nil || !user.IsOrganizer {
		// Fallback check: see if there's an approved profile
		profile, _ := partnerRepo.FindByUserID(c.Context(), userID)
		if profile == nil || profile.Status != "approved" {
			return utils.ErrorResponse(c, 403, "Only approved organizers can create events")
		}
	}

	// Enterprise check: Only 'event' or 'individual' categories (default events) can create standard events
	if user != nil && user.OrganizerCategory != "" && user.OrganizerCategory != "event" && user.OrganizerCategory != "individual" && user.OrganizerCategory != "creator" && user.OrganizerCategory != "company" && user.OrganizerCategory != "non-profit" {
		// If it's play or dining, they should use those specific creation endpoints
		return utils.ErrorResponse(c, 403, "This account is registered for "+user.OrganizerCategory+". Please use the correct dashboard.")
	}

	var event models.Event
	if err := c.Bind().Body(&event); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	event.ID = utils.GenerateUUIDv7()
	event.OrganizerID = userID
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	if err := eventRepo.Create(c.Context(), &event); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to create event")
	}

	return utils.SuccessResponse(c, 201, "Event created successfully", event)
}

func GetOrganizerEvents(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)

	// Strict access check
	user, err := userRepo.FindByID(c.Context(), userID)
	if err != nil || user == nil || (user.OrganizerCategory != "event" && user.OrganizerCategory != "individual" && user.OrganizerCategory != "creator" && user.OrganizerCategory != "company" && user.OrganizerCategory != "non-profit" && user.OrganizerCategory != "") {
		// If category is empty, check if they have a pending/approved profile
		profile, _ := partnerRepo.FindByUserID(c.Context(), userID)
		if profile != nil && profile.OrganizationDetails.Category != "event" && profile.OrganizationDetails.Category != "individual" && profile.OrganizationDetails.Category != "non-profit" {
			return utils.ErrorResponse(c, 403, "Access denied: Incorrect organizer category")
		}
	}

	limitStr := c.Query("limit", "20")
	limit, _ := strconv.Atoi(limitStr)
	cursor := c.Query("cursor", "")

	events, nextCursor, err := eventRepo.FindPaginatedByOrganizerID(c.Context(), userID, limit, cursor)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to fetch events")
	}

	return utils.SuccessResponse(c, 200, "Events fetched successfully", fiber.Map{
		"items":  events,
		"cursor": nextCursor,
	})
}

func GetAllEvents(c fiber.Ctx) error {
	limitStr := c.Query("limit", "20")
	limit, _ := strconv.Atoi(limitStr)
	cursor := c.Query("cursor", "")
	category := c.Query("category", "")
	city := c.Query("city", "")
	searchQuery := c.Query("q", "")

	events, nextCursor, err := eventRepo.GetPaginated(c.Context(), limit, cursor, category, city, searchQuery)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to fetch events")
	}

	return utils.SuccessResponse(c, 200, "Events fetched successfully", fiber.Map{
		"items":  events,
		"cursor": nextCursor,
	})
}

func GetEventByID(c fiber.Ctx) error {
	id := c.Params("id")
	event, err := eventRepo.FindByID(c.Context(), id)
	if err != nil || event == nil {
		return utils.ErrorResponse(c, 404, "Event not found")
	}
	return utils.SuccessResponse(c, 200, "Event fetched successfully", event)
}

func UpdateEvent(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)
	eventID := c.Params("id")

	existingEvent, err := eventRepo.FindByID(c.Context(), eventID)
	if err != nil || existingEvent == nil {
		return utils.ErrorResponse(c, 404, "Event not found")
	}

	if existingEvent.OrganizerID != userID {
		return utils.ErrorResponse(c, 403, "You are not authorized to update this event")
	}

	var updatedEvent models.Event
	if err := c.Bind().Body(&updatedEvent); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	updatedEvent.ID = eventID
	updatedEvent.OrganizerID = userID
	updatedEvent.CreatedAt = existingEvent.CreatedAt
	updatedEvent.UpdatedAt = time.Now()

	if err := eventRepo.Update(c.Context(), &updatedEvent); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to update event")
	}

	return utils.SuccessResponse(c, 200, "Event updated successfully", updatedEvent)
}

// AdminCreateEvent - Admin can create events without organizer checks
func AdminCreateEvent(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)

	var event models.Event
	if err := c.Bind().Body(&event); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	event.ID = utils.GenerateUUIDv7()
	event.OrganizerID = userID
	event.Status = "active"
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	if err := eventRepo.Create(c.Context(), &event); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to create event")
	}

	return utils.SuccessResponse(c, 201, "Event created successfully (admin)", event)
}

// AdminUpdateEvent - Admin can update any event
func AdminUpdateEvent(c fiber.Ctx) error {
	id := c.Params("id")
	var event models.Event
	if err := c.Bind().Body(&event); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	existing, err := eventRepo.FindByID(c.Context(), id)
	if err != nil || existing == nil {
		return utils.ErrorResponse(c, 404, "Event not found")
	}

	event.ID = id
	event.OrganizerID = existing.OrganizerID
	event.CreatedAt = existing.CreatedAt
	event.UpdatedAt = time.Now()

	if err := eventRepo.Update(c.Context(), &event); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to update event")
	}

	return utils.SuccessResponse(c, 200, "Event updated successfully (admin)", event)
}
