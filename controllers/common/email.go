package common

import (
	"fmt"

	"backend/utils"

	"github.com/gofiber/fiber/v3"
)

// ----- Request types -----

// SendPassPurchaseEmailRequest is the payload for a pass-purchase confirmation email.
type SendPassPurchaseEmailRequest struct {
	Email        string `json:"email" validate:"required,email"`
	Name         string `json:"name" validate:"required"`
	PassID       string `json:"passId" validate:"required"`
	Amount       int    `json:"amount" validate:"required,min=1"`
	PurchaseDate string `json:"purchaseDate" validate:"required"`
	ExpiryDate   string `json:"expiryDate" validate:"required"`
}

// SendPassRenewalEmailRequest is the payload for a pass-renewal confirmation email.
type SendPassRenewalEmailRequest struct {
	Email         string `json:"email" validate:"required,email"`
	Name          string `json:"name" validate:"required"`
	PassID        string `json:"passId" validate:"required"`
	RenewalDate   string `json:"renewalDate" validate:"required"`
	NewExpiryDate string `json:"newExpiryDate" validate:"required"`
	HTML          string `json:"html"`
}

// SendPassExpiryReminderEmailRequest is the payload for a pass expiry reminder email.
type SendPassExpiryReminderEmailRequest struct {
	Email         string `json:"email" validate:"required,email"`
	Name          string `json:"name" validate:"required"`
	PassID        string `json:"passId" validate:"required"`
	ExpiryDate    string `json:"expiryDate" validate:"required"`
	DaysRemaining int    `json:"daysRemaining" validate:"required,min=1"`
	HTML          string `json:"html"`
}

// SendBookingConfirmationEmailRequest is the payload for a booking confirmation email.
type SendBookingConfirmationEmailRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Name      string `json:"name" validate:"required"`
	BookingID string `json:"bookingId" validate:"required"`
	// "event" | "play" | "dining"
	BookingType    string `json:"bookingType" validate:"required"`
	VenueName      string `json:"venueName" validate:"required"`
	BookingDate    string `json:"bookingDate" validate:"required"`
	TotalAmount    int    `json:"totalAmount" validate:"required,min=0"`
	OriginalAmount int    `json:"originalAmount" validate:"required,min=1"`
	// "discount" | "free-booking" | ""
	PassBenefitApplied string `json:"passBenefitApplied"`
	SavingsAmount      int    `json:"savingsAmount" validate:"required,min=0"`
	HTML               string `json:"html"`
}

// ----- Handlers -----

// SendPassPurchaseEmail sends a purchase confirmation email after a pass is bought.
func SendPassPurchaseEmail(c fiber.Ctx) error {
	var req SendPassPurchaseEmailRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	html := utils.GetPassPurchaseEmailTemplate(req.Name, req.PassID, req.PurchaseDate, req.ExpiryDate, req.Amount)
	if err := utils.SendEmail(utils.EmailPlay, req.Email, fmt.Sprintf("🎉 Ticpin Pass Purchased! Pass ID: %s", req.PassID), html); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to send email")
	}

	return utils.SuccessResponse(c, 200, "Purchase confirmation email sent", nil)
}

// SendPassRenewalEmail sends a renewal confirmation email.
func SendPassRenewalEmail(c fiber.Ctx) error {
	var req SendPassRenewalEmailRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	if err := utils.SendEmail(utils.EmailPlay, req.Email, fmt.Sprintf("✅ Ticpin Pass Renewed! Valid until %s", req.NewExpiryDate), req.HTML); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to send email")
	}

	return utils.SuccessResponse(c, 200, "Renewal confirmation email sent", nil)
}

// SendPassExpiryReminderEmail sends a reminder before a pass expires.
func SendPassExpiryReminderEmail(c fiber.Ctx) error {
	var req SendPassExpiryReminderEmailRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	if err := utils.SendEmail(utils.EmailPlay, req.Email, fmt.Sprintf("⏰ Your Ticpin Pass expires in %d days!", req.DaysRemaining), req.HTML); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to send email")
	}

	return utils.SuccessResponse(c, 200, "Expiry reminder email sent", nil)
}

// SendBookingConfirmationEmail sends a booking confirmation email.
func SendBookingConfirmationEmail(c fiber.Ctx) error {
	var req SendBookingConfirmationEmailRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	if err := utils.SendEmail(utils.EmailPlay, req.Email, fmt.Sprintf("✅ Booking Confirmed - ID: %s", req.BookingID), req.HTML); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to send email")
	}

	return utils.SuccessResponse(c, 200, "Booking confirmation email sent", nil)
}
