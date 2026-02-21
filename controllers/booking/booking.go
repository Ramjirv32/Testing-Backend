package booking

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"

	"backend/models"
	"backend/repository"
	"backend/utils"
)

var playBookingRepo = repository.NewPlayBookingRepository()
var diningBookingRepo = repository.NewDiningBookingRepository()
var eventBookingRepo = repository.NewEventBookingRepository()

var playRepo = repository.NewPlayRepository()
var diningRepo = repository.NewDiningRepository()
var eventRepo = repository.NewEventRepository()
var userRepo = repository.NewUserRepository()

func CreatePlayBooking(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)
	var req models.CreatePlayBookingRequest

	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	if req.VenueID == "" || req.Sport == "" || req.Date == "" || req.TimeSlot == "" || req.PlayerName == "" || req.BillingEmail == "" {
		return utils.ErrorResponse(c, 400, "All required fields must be provided")
	}

	// Fetch Venue to get Organizer info
	venue, err := playRepo.FindByID(c.Context(), req.VenueID)
	if err != nil {
		return utils.ErrorResponse(c, 404, "Venue not found")
	}

	// Fetch Organizer Email
	organizer, _ := userRepo.FindByID(c.Context(), venue.OrganizerID)
	organizerEmail := ""
	if organizer != nil {
		organizerEmail = organizer.Email
	}

	// Check slot availability: allow up to venue.SlotSettings.TotalCourts simultaneous bookings
	available, err := playBookingRepo.CheckPlayAvailability(c.Context(), req.VenueID, req.Date, req.TimeSlot, venue.SlotSettings.TotalCourts)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to check slot availability")
	}
	if !available {
		courts := venue.SlotSettings.TotalCourts
		if courts <= 0 {
			courts = 1
		}
		return utils.ErrorResponse(c, 409, fmt.Sprintf("This time slot is fully booked (%d/%d courts taken). Please choose a different slot.", courts, courts))
	}

	bookingID := utils.GenerateUUIDv7()
	seqID := utils.GetNextSeqID()

	booking := &models.PlayBooking{
		ID:                 bookingID,
		SeqID:              seqID,
		UserID:             userID,
		VenueID:            req.VenueID,
		VenueName:          req.VenueName,
		Sport:              req.Sport,
		Date:               req.Date,
		TimeSlot:           req.TimeSlot,
		PlayerName:         req.PlayerName,
		Price:              req.Price,
		BillingEmail:       req.BillingEmail,
		BillingState:       req.BillingState,
		BillingNationality: req.BillingNationality,
		PaymentID:          req.PaymentID,
		PaymentGateway:     req.PaymentGateway,
		PaymentAmount:      req.PaymentAmount,
		OrganizerID:        venue.OrganizerID,
		OrganizerEmail:     organizerEmail,
		Status:             models.BookingConfirmed,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := playBookingRepo.Create(c.Context(), booking); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to create booking")
	}

	// Send confirmation email
	emailBody := utils.GetPlayBookingEmailTemplate(
		booking.PlayerName,
		booking.VenueName,
		booking.Sport,
		booking.Date,
		booking.TimeSlot,
		booking.ID,
	)
	go func() {
		utils.SendEmail(utils.EmailPlay, booking.BillingEmail, "Booking Confirmed - "+booking.VenueName, emailBody)
	}()

	return utils.SuccessResponse(c, 201, "Play booking created successfully", booking)
}

func CreateDiningBooking(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)
	var req models.CreateDiningBookingRequest

	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	if req.RestaurantID == "" || req.Date == "" || req.GuestName == "" {
		return utils.ErrorResponse(c, 400, "All required fields must be provided")
	}

	// Fetch Restaurant to get Organizer info
	restaurant, err := diningRepo.FindByID(c.Context(), req.RestaurantID)
	if err != nil {
		return utils.ErrorResponse(c, 404, "Restaurant not found")
	}

	organizer, _ := userRepo.FindByID(c.Context(), restaurant.OrganizerID)
	organizerEmail := ""
	if organizer != nil {
		organizerEmail = organizer.Email
	}

	bookingID := utils.GenerateUUIDv7()
	seqID := utils.GetNextSeqID()

	booking := &models.DiningBooking{
		ID:             bookingID,
		SeqID:          seqID,
		UserID:         userID,
		RestaurantID:   req.RestaurantID,
		RestaurantName: req.RestaurantName,
		Date:           req.Date,
		TimeSlot:       req.TimeSlot,
		GuestCount:     req.GuestCount,
		GuestName:      req.GuestName,
		BillingEmail:   req.BillingEmail,
		SpecialRequest: req.SpecialRequest,
		PaymentID:      req.PaymentID,
		PaymentGateway: req.PaymentGateway,
		PaymentAmount:  req.PaymentAmount,
		OrganizerID:    restaurant.OrganizerID,
		OrganizerEmail: organizerEmail,
		Status:         models.BookingConfirmed,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := diningBookingRepo.Create(c.Context(), booking); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to create booking")
	}

	// Send confirmation email — prefer billing email, fallback to user account email
	emailTo := booking.BillingEmail
	if emailTo == "" {
		user, _ := userRepo.FindByID(c.Context(), userID)
		if user != nil {
			emailTo = user.Email
		}
	}
	if emailTo != "" {
		emailBody := utils.GetDiningBookingEmailTemplate(
			booking.GuestName,
			booking.RestaurantName,
			booking.Date,
			booking.TimeSlot,
			booking.ID,
			booking.GuestCount,
			booking.SpecialRequest,
		)
		go func() {
			utils.SendEmail(utils.EmailDining, emailTo, "Table Reserved - "+booking.RestaurantName, emailBody)
		}()
	}

	return utils.SuccessResponse(c, 201, "Dining booking created successfully", booking)
}

func CreateEventBooking(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)
	var req models.CreateEventBookingRequest

	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	if req.EventID == "" || req.EventTitle == "" || req.Quantity < 1 || req.GuestName == "" || req.BillingEmail == "" {
		return utils.ErrorResponse(c, 400, "All required fields must be provided")
	}

	// Fetch Event to get Organizer info
	event, err := eventRepo.FindByID(c.Context(), req.EventID)
	if err != nil {
		return utils.ErrorResponse(c, 404, "Event not found")
	}

	// Fetch Organizer Email
	organizer, _ := userRepo.FindByID(c.Context(), event.OrganizerID)
	organizerEmail := ""
	if organizer != nil {
		organizerEmail = organizer.Email
	}

	bookingID := utils.GenerateUUIDv7()
	seqID := utils.GetNextSeqID()

	booking := &models.EventBooking{
		ID:             bookingID,
		SeqID:          seqID,
		UserID:         userID,
		EventID:        req.EventID,
		EventTitle:     req.EventTitle,
		TicketType:     req.TicketType,
		SeatType:       req.SeatType,
		Quantity:       req.Quantity,
		UnitPrice:      req.UnitPrice,
		TotalPrice:     req.UnitPrice * float64(req.Quantity),
		GuestName:      req.GuestName,
		BillingEmail:   req.BillingEmail,
		BillingState:   req.BillingState,
		PaymentID:      req.PaymentID,
		PaymentGateway: req.PaymentGateway,
		PaymentAmount:  req.PaymentAmount,
		OrganizerID:    event.OrganizerID,
		OrganizerEmail: organizerEmail,
		Status:         models.BookingConfirmed,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := eventBookingRepo.Create(c.Context(), booking); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to create event booking")
	}

	// Send confirmation email
	// Template needs: guestName, eventName, venue, date, time string, ticketCount int, qrImageURL string
	eventDate := event.StartDatetime.Format("02 Jan 2006")
	eventTime := event.StartDatetime.Format("03:04 PM")
	emailBody := utils.GetEventBookingEmailTemplate(
		booking.GuestName,
		booking.EventTitle,
		event.Venue.Name,
		eventDate,
		eventTime,
		booking.Quantity,
		"https://api.qrserver.com/v1/create-qr-code/?size=150x150&data="+booking.ID,
		booking.ID,
	)
	go func() {
		utils.SendEmail(utils.EmailEvents, booking.BillingEmail, "Tickets Confirmed - "+booking.EventTitle, emailBody)
	}()

	return utils.SuccessResponse(c, 201, "Event booking created successfully", booking)
}

func GetUserBookings(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)
	bookingType := c.Query("type")

	if bookingType == "play" {
		bookings, err := playBookingRepo.FindByUserID(c.Context(), userID)
		if err != nil {
			return utils.ErrorResponse(c, 500, "Failed to fetch bookings")
		}
		return utils.SuccessResponse(c, 200, "Bookings fetched", fiber.Map{"bookings": bookings})
	} else if bookingType == "dining" {
		bookings, err := diningBookingRepo.FindByUserID(c.Context(), userID)
		if err != nil {
			return utils.ErrorResponse(c, 500, "Failed to fetch bookings")
		}
		return utils.SuccessResponse(c, 200, "Bookings fetched", fiber.Map{"bookings": bookings})
	} else if bookingType == "event" {
		bookings, err := eventBookingRepo.FindByUserID(c.Context(), userID)
		if err != nil {
			return utils.ErrorResponse(c, 500, "Failed to fetch bookings")
		}
		return utils.SuccessResponse(c, 200, "Bookings fetched", fiber.Map{"bookings": bookings})
	}

	pBk, _ := playBookingRepo.FindByUserID(c.Context(), userID)
	dBk, _ := diningBookingRepo.FindByUserID(c.Context(), userID)
	eBk, _ := eventBookingRepo.FindByUserID(c.Context(), userID)

	return utils.SuccessResponse(c, 200, "Bookings fetched", fiber.Map{
		"play_bookings":   pBk,
		"dining_bookings": dBk,
		"event_bookings":  eBk,
	})
}

func GetBookingByID(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)
	bookingID := c.Params("id")
	bookingType := c.Query("type", "play")

	if bookingType == "play" {
		booking, err := playBookingRepo.FindByID(c.Context(), bookingID)
		if err != nil || booking == nil {
			return utils.ErrorResponse(c, 404, "Booking not found")
		}
		if booking.UserID != userID {
			return utils.ErrorResponse(c, 403, "Unauthorized")
		}
		return utils.SuccessResponse(c, 200, "Booking fetched", booking)
	} else if bookingType == "dining" {
		booking, err := diningBookingRepo.FindByID(c.Context(), bookingID)
		if err != nil || booking == nil {
			return utils.ErrorResponse(c, 404, "Booking not found")
		}
		if booking.UserID != userID {
			return utils.ErrorResponse(c, 403, "Unauthorized")
		}
		return utils.SuccessResponse(c, 200, "Booking fetched", booking)
	} else {
		booking, err := eventBookingRepo.FindByID(c.Context(), bookingID)
		if err != nil || booking == nil {
			return utils.ErrorResponse(c, 404, "Booking not found")
		}
		if booking.UserID != userID {
			return utils.ErrorResponse(c, 403, "Unauthorized")
		}
		return utils.SuccessResponse(c, 200, "Booking fetched", booking)
	}
}

func CancelBooking(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)
	bookingID := c.Params("id")
	bookingType := c.Query("type", "play")

	if bookingType == "play" {
		booking, err := playBookingRepo.FindByID(c.Context(), bookingID)
		if err != nil || booking == nil {
			return utils.ErrorResponse(c, 404, "Booking not found")
		}
		if booking.UserID != userID {
			return utils.ErrorResponse(c, 403, "Unauthorized")
		}
		booking.Status = models.BookingCancelled
		booking.UpdatedAt = time.Now()
		if err := playBookingRepo.Update(c.Context(), booking); err != nil {
			return utils.ErrorResponse(c, 500, "Failed to cancel booking")
		}
		return utils.SuccessResponse(c, 200, "Booking cancelled", booking)
	} else if bookingType == "dining" {
		booking, err := diningBookingRepo.FindByID(c.Context(), bookingID)
		if err != nil || booking == nil {
			return utils.ErrorResponse(c, 404, "Booking not found")
		}
		if booking.UserID != userID {
			return utils.ErrorResponse(c, 403, "Unauthorized")
		}
		booking.Status = models.BookingCancelled
		booking.UpdatedAt = time.Now()
		if err := diningBookingRepo.Update(c.Context(), booking); err != nil {
			return utils.ErrorResponse(c, 500, "Failed to cancel booking")
		}
		return utils.SuccessResponse(c, 200, "Booking cancelled", booking)
	} else {
		booking, err := eventBookingRepo.FindByID(c.Context(), bookingID)
		if err != nil || booking == nil {
			return utils.ErrorResponse(c, 404, "Booking not found")
		}
		if booking.UserID != userID {
			return utils.ErrorResponse(c, 403, "Unauthorized")
		}
		booking.Status = models.BookingCancelled
		booking.UpdatedAt = time.Now()
		if err := eventBookingRepo.Update(c.Context(), booking); err != nil {
			return utils.ErrorResponse(c, 500, "Failed to cancel booking")
		}
		return utils.SuccessResponse(c, 200, "Booking cancelled", booking)
	}
}

func DeleteAllPlayBookings(c fiber.Ctx) error {
	if err := playBookingRepo.DeleteAll(c.Context()); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to delete bookings")
	}
	return utils.SuccessResponse(c, 200, "All play bookings deleted", fiber.Map{})
}

func DeleteAllDiningBookings(c fiber.Ctx) error {
	if err := diningBookingRepo.DeleteAll(c.Context()); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to delete bookings")
	}
	return utils.SuccessResponse(c, 200, "All dining bookings deleted", fiber.Map{})
}

func DeleteAllUsers(c fiber.Ctx) error {
	if err := userRepo.DeleteAll(c.Context()); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to delete users")
	}
	return utils.SuccessResponse(c, 200, "All users deleted", fiber.Map{})
}
