package event

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"

	"backend/models"
	"backend/repository"
	"backend/utils"
)

var eventRepo = repository.NewEventRepository()
var userRepo = repository.NewUserRepository()
var partnerRepo = repository.NewPartnerRepository()

func CreateEvent(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)

	user, err := userRepo.FindByID(c.Context(), userID)
	if err != nil || user == nil || !user.IsOrganizer {

		profile, _ := partnerRepo.FindByUserID(c.Context(), userID)
		if profile == nil || profile.Status != "approved" {
			return utils.ErrorResponse(c, 403, "Only approved organizers can create events")
		}
	}

	if user != nil && user.IsOrganizer && !user.HasOrganizerCategory("event") && !user.HasOrganizerCategory("individual") && !user.HasOrganizerCategory("creator") && !user.HasOrganizerCategory("company") && !user.HasOrganizerCategory("non-profit") && user.OrganizerCategory != "" {
		return utils.ErrorResponse(c, 403, "This account is not registered for events. Please use the correct dashboard.")
	}

	var event models.Event
	if err := c.Bind().Body(&event); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	event.ID = utils.GenerateUUIDv7()
	event.OrganizerID = userID
	event.Status = "pending" // Organizer events need admin approval
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	if err := eventRepo.Create(c.Context(), &event); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to create event")
	}

	return utils.SuccessResponse(c, 201, "Event created successfully", event)
}

func GetOrganizerEvents(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)

	user, err := userRepo.FindByID(c.Context(), userID)
	hasEventAccess := user != nil && (user.HasOrganizerCategory("event") || user.HasOrganizerCategory("individual") || user.HasOrganizerCategory("creator") || user.HasOrganizerCategory("company") || user.HasOrganizerCategory("non-profit") || user.OrganizerCategory == "")
	if err != nil || user == nil || !hasEventAccess {
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
		// Return empty list on error (e.g. no events created yet)
		return utils.SuccessResponse(c, 200, "Events fetched successfully", fiber.Map{
			"items":  []*models.Event{},
			"cursor": "",
		})
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

	events, nextCursor, err := eventRepo.GetPaginated(c.Context(), limit, cursor, category, city, searchQuery, status)
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
	updatedEvent.Status = existingEvent.Status // Preserve status - organizers can't change it, only admins can
	updatedEvent.CreatedAt = existingEvent.CreatedAt
	updatedEvent.UpdatedAt = time.Now()

	if err := eventRepo.Update(c.Context(), &updatedEvent); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to update event")
	}

	return utils.SuccessResponse(c, 200, "Event updated successfully", updatedEvent)
}

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

func AdminApproveEvent(c fiber.Ctx) error {
	id := c.Params("id")
	status := c.Query("status", "active")
	reason := c.Query("reason", "")

	existing, err := eventRepo.FindByID(c.Context(), id)
	if err != nil || existing == nil {
		return utils.ErrorResponse(c, 404, "Event not found")
	}

	existing.Status = status
	existing.RejectionReason = reason
	existing.UpdatedAt = time.Now()

	if err := eventRepo.Update(c.Context(), existing); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to update event status")
	}

	return utils.SuccessResponse(c, 200, "Event status updated successfully", nil)
}

func DeleteEvent(c fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("uid").(string)

	existing, err := eventRepo.FindByID(c.Context(), id)
	if err != nil || existing == nil {
		return utils.ErrorResponse(c, 404, "Event not found")
	}

	if existing.OrganizerID != userID {
		return utils.ErrorResponse(c, 403, "Access denied: You do not own this event")
	}

	if err := eventRepo.Delete(c.Context(), id); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to delete event")
	}

	return utils.SuccessResponse(c, 200, "Event deleted successfully", nil)
}

func AdminDeleteEvent(c fiber.Ctx) error {
	id := c.Params("id")

	existing, err := eventRepo.FindByID(c.Context(), id)
	if err != nil || existing == nil {
		return utils.ErrorResponse(c, 404, "Event not found")
	}

	if err := eventRepo.Delete(c.Context(), id); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to delete event")
	}

	return utils.SuccessResponse(c, 200, "Event deleted successfully (admin)", nil)
}

func ResubmitEvent(c fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("uid").(string)

	existing, err := eventRepo.FindByID(c.Context(), id)
	if err != nil || existing == nil {
		return utils.ErrorResponse(c, 404, "Event not found")
	}

	if existing.OrganizerID != userID {
		return utils.ErrorResponse(c, 403, "Access denied: You do not own this event")
	}

	if existing.Status != "rejected" {
		return utils.ErrorResponse(c, 400, "Only rejected events can be resubmitted")
	}

	existing.Status = "pending"
	existing.RejectionReason = ""

	if err := eventRepo.Update(c.Context(), existing); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to resubmit event")
	}

	return utils.SuccessResponse(c, 200, "Event resubmitted for review", nil)
}
