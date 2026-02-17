package controllers

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

func CreatePlayBooking(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)
	var req models.CreatePlayBookingRequest

	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	if req.VenueID == "" || req.Sport == "" || req.Date == "" || req.TimeSlot == "" || req.PlayerName == "" || req.BillingEmail == "" || req.BillingState == "" {
		return utils.ErrorResponse(c, 400, "All required fields must be provided")
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
		Status:             models.BookingConfirmed,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := playBookingRepo.Create(c.Context(), booking); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to create booking")
	}

	emailBody := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e0e0e0; border-radius: 10px;">
			<h2 style="color: #333; text-align: center;">Booking Confirmed!</h2>
			<p>Hello %s,</p>
			<p>Your booking for <b>%s</b> has been confirmed.</p>
			<div style="background-color: #f9f9f9; padding: 15px; border-radius: 5px; margin: 20px 0;">
				<p><strong>Sport:</strong> %s</p>
				<p><strong>Date:</strong> %s</p>
				<p><strong>Time Slot:</strong> %s</p>
				<p><strong>Venue:</strong> %s</p>
			</div>
			<p>Thank you for choosing TicPin!</p>
			<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
			<p style="font-size: 12px; color: #777; text-align: center;">&copy; 2026 TicPin. All rights reserved.</p>
		</div>
	`, booking.PlayerName, booking.VenueName, booking.Sport, booking.Date, booking.TimeSlot, booking.VenueName)

	go utils.SendEmail(utils.EmailPlay, booking.BillingEmail, "Booking Confirmation - "+booking.VenueName, emailBody)

	return utils.SuccessResponse(c, 201, "Play booking created successfully", booking)
}

func CreateDiningBooking(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)
	var req models.CreateDiningBookingRequest

	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	if req.RestaurantID == "" || req.Date == "" || req.GuestCount < 1 || req.GuestName == "" {
		return utils.ErrorResponse(c, 400, "All required fields must be provided")
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
		SpecialRequest: req.SpecialRequest,
		Status:         models.BookingConfirmed,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := diningBookingRepo.Create(c.Context(), booking); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to create booking")
	}


	user, _ := userRepo.FindByID(c.Context(), userID)
	toEmail := ""
	if user != nil && user.Email != "" {
		toEmail = user.Email
	}

	if toEmail != "" {
		emailBody := fmt.Sprintf(`
			<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e0e0e0; border-radius: 10px;">
				<h2 style="color: #333; text-align: center;">Dining Booking Confirmed!</h2>
				<p>Hello %s,</p>
				<p>Your table at <b>%s</b> has been reserved.</p>
				<div style="background-color: #f9f9f9; padding: 15px; border-radius: 5px; margin: 20px 0;">
					<p><strong>Restaurant:</strong> %s</p>
					<p><strong>Date:</strong> %s</p>
					<p><strong>Time Slot:</strong> %s</p>
					<p><strong>Guests:</strong> %d</p>
				</div>
				<p>Enjoy your meal!</p>
				<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
				<p style="font-size: 12px; color: #777; text-align: center;">&copy; 2026 TicPin. All rights reserved.</p>
			</div>
		`, booking.GuestName, booking.RestaurantName, booking.RestaurantName, booking.Date, booking.TimeSlot, booking.GuestCount)

		go utils.SendEmail(utils.EmailDining, toEmail, "Dining Confirmation - "+booking.RestaurantName, emailBody)
	}

	return utils.SuccessResponse(c, 201, "Dining booking created successfully", booking)
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
	}

	playBookings, err1 := playBookingRepo.FindByUserID(c.Context(), userID)
	diningBookings, err2 := diningBookingRepo.FindByUserID(c.Context(), userID)

	if err1 != nil || err2 != nil {
		return utils.ErrorResponse(c, 500, "Failed to fetch bookings")
	}

	return utils.SuccessResponse(c, 200, "Bookings fetched", fiber.Map{
		"play_bookings":   playBookings,
		"dining_bookings": diningBookings,
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
	} else {
		booking, err := diningBookingRepo.FindByID(c.Context(), bookingID)
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
		if err := playBookingRepo.Update(c.Context(), booking); err != nil {
			return utils.ErrorResponse(c, 500, "Failed to cancel booking")
		}
		return utils.SuccessResponse(c, 200, "Booking cancelled successfully", booking)
	} else {
		booking, err := diningBookingRepo.FindByID(c.Context(), bookingID)
		if err != nil || booking == nil {
			return utils.ErrorResponse(c, 404, "Booking not found")
		}
		if booking.UserID != userID {
			return utils.ErrorResponse(c, 403, "Unauthorized")
		}
		booking.Status = models.BookingCancelled
		if err := diningBookingRepo.Update(c.Context(), booking); err != nil {
			return utils.ErrorResponse(c, 500, "Failed to cancel booking")
		}
		return utils.SuccessResponse(c, 200, "Booking cancelled successfully", booking)
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
