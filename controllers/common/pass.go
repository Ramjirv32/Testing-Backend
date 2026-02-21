package common

import (
	"backend/config"
	"backend/utils"

	"github.com/gofiber/fiber/v3"
)

// CheckPassEligibilityRequest checks eligibility by email AND phone for uniqueness.
type CheckPassEligibilityRequest struct {
	Email string `json:"email" validate:"omitempty,email"`
	Phone string `json:"phone"`
}

// CheckPassEligibility validates whether a user can purchase a new Ticpin Pass.
// A pass is rejected if EITHER the email OR the phone already has an active pass.
func CheckPassEligibility(c fiber.Ctx) error {
	var req CheckPassEligibilityRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	if req.Email == "" && req.Phone == "" {
		return utils.ErrorResponse(c, 400, "Email or phone is required")
	}

	result, err := utils.CheckPassEligibility(c.Context(), req.Email, req.Phone, config.FirestoreClient)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to check pass eligibility")
	}

	if !result.Eligible {
		return utils.SuccessResponse(c, 200, "Pass eligibility checked", fiber.Map{
			"eligible":         false,
			"has_active_pass":  true,
			"existing_pass_id": result.ExistingPassID,
			"expiry_date":      result.ExpiryDate,
			"matched_by":       result.MatchedBy,
			"reason":           result.Reason,
		})
	}

	return utils.SuccessResponse(c, 200, "User is eligible to purchase a Ticpin Pass", fiber.Map{
		"eligible":        true,
		"has_active_pass": false,
	})
}
