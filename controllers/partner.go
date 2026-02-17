package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"

	"backend/models"
	"backend/repository"
	"backend/tasks"
	"backend/utils"
)

var partnerRepo = repository.NewPartnerRepository()

func SubmitVerification(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)
	var profile models.PartnerProfile
	if err := c.Bind().Body(&profile); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	// Basic validation
	if profile.OrganizationDetails.PAN == "" {
		return utils.ErrorResponse(c, 400, "PAN number is required")
	}

	// Validate backup contact is not same as user contact
	user, err := userRepo.FindByID(c.Context(), userID)
	if err == nil && user != nil {
		if profile.BackupContact.Email != "" && user.Email != "" && profile.BackupContact.Email == user.Email {
			return utils.ErrorResponse(c, 400, "Backup email cannot be the same as your primary email")
		}
		if profile.BackupContact.Phone != "" && user.Phone != "" && profile.BackupContact.Phone == user.Phone {
			return utils.ErrorResponse(c, 400, "Backup phone cannot be the same as your primary phone")
		}
	}

	profile.ID = utils.GenerateUUIDv7()
	profile.UserID = userID
	profile.Status = "pending"
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()

	if err := partnerRepo.Create(c.Context(), &profile); err != nil {
		fmt.Printf("Error creating partner profile: %v\n", err)
		return utils.ErrorResponse(c, 500, "Failed to submit verification")
	}

	// Send verification email
	userEmail := ""
	if user != nil {
		userEmail = user.Email
	}
	if profile.BackupContact.Email != "" {
		userEmail = profile.BackupContact.Email
	}

	if userEmail != "" {
		body := fmt.Sprintf(`
			<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e0e0e0; border-radius: 10px;">
				<h2 style="color: #333; text-align: center;">Organizer Verification - TicPin</h2>
				<p>Hello %s,</p>
				<p>Thank you for showing interest in listing as a <b>%s</b> on TicPin. We have received your verification request.</p>
				<p>Our team will review your details (PAN: %s, Organization: %s) and respond to you within <b>24 hours</b>.</p>
				<p>Once approved, you will have access to your partner dashboard.</p>
				<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
				<p style="font-size: 12px; color: #777; text-align: center;">&copy; 2026 TicPin. All rights reserved.</p>
			</div>
		`, profile.OrganizationDetails.PANName, profile.OrganizationDetails.Category, profile.OrganizationDetails.PAN, profile.OrganizationDetails.PANName)

		tasks.EnqueueEmail(utils.EmailAdmin, userEmail, "Partner Verification Received - TicPin", body)
	}

	return utils.SuccessResponse(c, 200, "Verification request submitted successfully", profile)
}

func GetEventPosters(c fiber.Ctx) error {
	isAdmin, _ := c.Locals("isAdmin").(bool)
	if !isAdmin {
		return utils.ErrorResponse(c, 403, "Forbidden")
	}

	limitStr := c.Query("limit", "20")
	limit, _ := strconv.Atoi(limitStr)
	cursor := c.Query("cursor", "")
	category := c.Query("category", "")

	profiles, nextCursor, err := partnerRepo.GetPaginated(c.Context(), limit, cursor, category)
	if err != nil {
		fmt.Printf("Error fetching partner profiles: %v\n", err)
		return utils.ErrorResponse(c, 500, "Failed to fetch partner profiles")
	}

	fmt.Printf("Fetched %d profiles for category: %s\n", len(profiles), category)

	return utils.SuccessResponse(c, 200, "Partner profiles fetched successfully", fiber.Map{
		"items":  profiles,
		"cursor": nextCursor,
	})
}

func ApproveEventPoster(c fiber.Ctx) error {
	isAdmin, _ := c.Locals("isAdmin").(bool)
	if !isAdmin {
		return utils.ErrorResponse(c, 403, "Forbidden")
	}

	id := c.Params("id")
	status := c.Query("status", "approved")

	if err := partnerRepo.UpdateStatus(c.Context(), id, status); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to update status")
	}

	profile, _ := partnerRepo.FindByID(c.Context(), id)
	if profile != nil {
		user, _ := userRepo.FindByID(c.Context(), profile.UserID)

		if status == "approved" && user != nil {
			user.IsOrganizer = true
			user.OrganizerCategory = profile.OrganizationDetails.Category
			_ = userRepo.Update(c.Context(), user)
		}

		userEmail := ""
		if user != nil {
			userEmail = user.Email
		}
		if profile.BackupContact.Email != "" {
			userEmail = profile.BackupContact.Email
		}

		if userEmail != "" {
			subject := fmt.Sprintf("%s Access Approved - TicPin", profile.OrganizationDetails.Category)
			message := fmt.Sprintf("Congratulations! Your %s partner verification has been approved. You can now access your dashboard.", profile.OrganizationDetails.Category)
			if status == "rejected" {
				subject = "Partner Verification Update - TicPin"
				message = "We regret to inform you that your partner verification request was not approved at this time. Please contact support for more details."
			}

			body := fmt.Sprintf(`
				<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e0e0e0; border-radius: 10px;">
					<h2 style="color: #333; text-align: center;">%s</h2>
					<p>Hello %s,</p>
					<p>%s</p>
					<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
					<p style="font-size: 12px; color: #777; text-align: center;">&copy; 2026 TicPin. All rights reserved.</p>
				</div>
			`, subject, profile.OrganizationDetails.PANName, message)

			tasks.EnqueueEmail(utils.EmailAdmin, userEmail, subject, body)
		}
	}

	return utils.SuccessResponse(c, 200, fmt.Sprintf("Profile %s successfully", status), nil)
}

func UpdatePartnerProfile(c fiber.Ctx) error {
	isAdmin, _ := c.Locals("isAdmin").(bool)
	if !isAdmin {
		return utils.ErrorResponse(c, 403, "Forbidden")
	}

	id := c.Params("id")
	var profile models.PartnerProfile
	if err := c.Bind().Body(&profile); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	profile.ID = id
	if err := partnerRepo.Update(c.Context(), &profile); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to update profile")
	}

	return utils.SuccessResponse(c, 200, "Profile updated successfully", profile)
}

func GetMyVerificationStatus(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)
	isAdmin, _ := c.Locals("isAdmin").(bool)
	if isAdmin {
		return utils.SuccessResponse(c, 200, "Admin auto-approved", fiber.Map{
			"status":   "approved",
			"user_id":  userID,
			"category": "all",
		})
	}
	category := c.Query("category", "")

	profile, err := partnerRepo.FindByUserIDAndCategory(c.Context(), userID, category)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Database error")
	}

	if profile != nil {
		return utils.SuccessResponse(c, 200, "Verification status fetched", profile)
	}

	user, _ := userRepo.FindByID(c.Context(), userID)
	if user != nil && user.IsOrganizer {
		// If category is provided, check if it matches the current user's category (legacy or specific)
		if category == "" || user.OrganizerCategory == category {
			return utils.SuccessResponse(c, 200, "User is approved organizer", fiber.Map{
				"status":   "approved",
				"user_id":  userID,
				"category": user.OrganizerCategory,
			})
		}
	}

	return utils.SuccessResponse(c, 200, "No verification found", nil)
}
