package partner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"

	"backend/config"
	"backend/models"
	"backend/repository"
	"backend/tasks"
	"backend/utils"
)

var partnerRepo = repository.NewPartnerRepository()
var userRepo = repository.NewUserRepository()
var auditRepo = repository.NewAuditRepository()

func GetAllPartners(c fiber.Ctx) error {
	partners, err := partnerRepo.GetAll(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to fetch partners")
	}
	return utils.SuccessResponse(c, 200, "Partners fetched successfully", partners)
}

func SubmitVerification(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)
	var profile models.PartnerProfile
	if err := c.Bind().Body(&profile); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	if profile.OrganizationDetails.Category == "" {
		return utils.ErrorResponse(c, 400, "Category is required")
	}

	existingProfile, _ := partnerRepo.FindByUserIDAndCategory(c.Context(), userID, profile.OrganizationDetails.Category)
	if existingProfile != nil {
		return utils.ErrorResponse(c, 400, "A verification profile already exists for this category. Duplicate creation is not allowed.")
	}

	// Get user to check PAN verification status
	user, err := userRepo.FindByID(c.Context(), userID)
	if err != nil || user == nil {
		return utils.ErrorResponse(c, 500, "Failed to find user")
	}

	// PAN must be verified via /verify-pan before submission — always pull from user record
	if !user.IsPanVerified || user.PanNumber == "" || user.PanDetails == nil {
		return utils.ErrorResponse(c, 400, "Please verify your PAN first before submitting verification.")
	}

	// Always use PAN details stored on user record (works for 1st, 2nd, 3rd category)
	profile.OrganizationDetails.PAN = user.PanNumber
	profile.OrganizationDetails.PANVerification = *user.PanDetails
	if profile.OrganizationDetails.PANName == "" {
		profile.OrganizationDetails.PANName = user.PanDetails.RegisteredName
	}

	// PAN uniqueness - different user cannot use same PAN
	existingPAN, _ := partnerRepo.FindByPAN(c.Context(), profile.OrganizationDetails.PAN)
	if existingPAN != nil && existingPAN.UserID != userID {
		return utils.ErrorResponse(c, 400, "This PAN is already registered under a different account. Each PAN can only be associated with one organizer account.")
	}

	// Bank account uniqueness - different user cannot use same bank account number
	if profile.BankDetails.AccountNumber != "" {
		existingBank, _ := partnerRepo.FindByBankAccountNumber(c.Context(), profile.BankDetails.AccountNumber)
		if existingBank != nil && existingBank.UserID != userID {
			return utils.ErrorResponse(c, 400, "This bank account number is already registered under a different account.")
		}
	}

	// Backup contact validation
	if profile.BackupContact.Email != "" && user.Email != "" && profile.BackupContact.Email == user.Email {
		return utils.ErrorResponse(c, 400, "Backup email cannot be the same as your primary email")
	}
	if profile.BackupContact.Phone != "" && user.Phone != "" && profile.BackupContact.Phone == user.Phone {
		return utils.ErrorResponse(c, 400, "Backup phone cannot be the same as your primary phone")
	}

	profile.ID = utils.GenerateUUIDv7()
	profile.UserID = userID
	profile.OrganizerEmail = user.Email
	profile.Status = "pending"
	// Set state from PAN verification if not provided
	if profile.State == "" && user.PanDetails != nil {
		profile.State = user.PanDetails.Address.State
	}
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()

	if err := partnerRepo.Create(c.Context(), &profile); err != nil {
		fmt.Printf("Error creating partner profile: %v\n", err)
		auditRepo.Log(c.Context(), &models.AuditLog{
			UserID:    userID,
			Action:    "VERIFICATION_SUBMISSION",
			Category:  profile.OrganizationDetails.Category,
			Status:    "FAILED",
			Details:   fmt.Sprintf("Failed to save profile: %v", err),
			IPAddress: c.IP(),
		})
		return utils.ErrorResponse(c, 500, "Failed to submit verification")
	}

	auditRepo.Log(c.Context(), &models.AuditLog{
		UserID:    userID,
		Action:    "VERIFICATION_SUBMISSION",
		Category:  profile.OrganizationDetails.Category,
		Status:    "SUCCESS",
		Details:   "Verification profile submitted pending review",
		IPAddress: c.IP(),
	})

	userEmail := user.Email
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

	limitStr := c.Query("limit", "20")
	limit, _ := strconv.Atoi(limitStr)
	cursor := c.Query("cursor", "")
	category := c.Query("category", "")
	status := c.Query("status", "")

	profiles, nextCursor, err := partnerRepo.GetPaginated(c.Context(), limit, cursor, category, status)
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

	id := c.Params("id")
	status := c.Query("status", "approved")

	if err := partnerRepo.UpdateStatus(c.Context(), id, status); err != nil {
		auditRepo.Log(c.Context(), &models.AuditLog{
			Action:    "ADMIN_PARTNER_VERIFICATION",
			Status:    "FAILED",
			Details:   fmt.Sprintf("Failed to update status for profile %s to %s: %v", id, status, err),
			IPAddress: c.IP(),
		})
		return utils.ErrorResponse(c, 500, "Failed to update status")
	}

	profile, _ := partnerRepo.FindByID(c.Context(), id)
	if profile != nil {
		user, _ := userRepo.FindByID(c.Context(), profile.UserID)

		if status == "approved" && user != nil {
			user.IsOrganizer = true
			user.IsPanVerified = true
			newCategory := profile.OrganizationDetails.Category
			// Keep OrganizerCategory as the first/primary category for backward compat
			if user.OrganizerCategory == "" {
				user.OrganizerCategory = newCategory
			}
			// Append to OrganizerCategories if not already present
			alreadyHas := false
			for _, cat := range user.OrganizerCategories {
				if cat == newCategory {
					alreadyHas = true
					break
				}
			}
			if !alreadyHas {
				user.OrganizerCategories = append(user.OrganizerCategories, newCategory)
			}
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

	auditRepo.Log(c.Context(), &models.AuditLog{
		Action:    "ADMIN_PARTNER_VERIFICATION",
		Status:    "SUCCESS",
		Details:   fmt.Sprintf("Admin %s profile %s", status, id),
		IPAddress: c.IP(),
	})

	return utils.SuccessResponse(c, 200, fmt.Sprintf("Profile %s successfully", status), nil)
}

func UpdatePartnerProfile(c fiber.Ctx) error {

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
			"status":               "approved",
			"user_id":              userID,
			"category":             "all",
			"organizer_categories": []string{"event", "dining", "play"},
		})
	}
	category := c.Query("category", "")

	user, _ := userRepo.FindByID(c.Context(), userID)

	// If specific category requested, return that profile
	if category != "" {
		profile, err := partnerRepo.FindByUserIDAndCategory(c.Context(), userID, category)
		if err != nil {
			return utils.ErrorResponse(c, 500, "Database error")
		}
		if profile != nil {
			return utils.SuccessResponse(c, 200, "Verification status fetched", profile)
		}
		// No profile for this category - check if user has PAN verified (for 2nd/3rd category)
		if user != nil && user.IsPanVerified && user.PanNumber != "" {
			return utils.SuccessResponse(c, 200, "PAN verified, new category submission needed", fiber.Map{
				"status":               "new_category",
				"user_id":              userID,
				"is_pan_verified":      true,
				"pan_number":           user.PanNumber,
				"organizer_categories": user.OrganizerCategories,
			})
		}
		return utils.SuccessResponse(c, 200, "No verification found", nil)
	}

	// No category - return all profiles for this user
	allProfiles, _ := partnerRepo.FindAllByUserID(c.Context(), userID)

	if user != nil && user.IsOrganizer {
		categories := user.OrganizerCategories
		if len(categories) == 0 && user.OrganizerCategory != "" {
			categories = []string{user.OrganizerCategory}
		}
		return utils.SuccessResponse(c, 200, "User is approved organizer", fiber.Map{
			"status":               "approved",
			"user_id":              userID,
			"category":             user.OrganizerCategory,
			"organizer_categories": categories,
			"is_pan_verified":      user.IsPanVerified,
			"profiles":             allProfiles,
		})
	}

	if len(allProfiles) > 0 {
		return utils.SuccessResponse(c, 200, "Verification status fetched", fiber.Map{
			"status":          allProfiles[0].Status,
			"user_id":         userID,
			"profiles":        allProfiles,
			"is_pan_verified": user != nil && user.IsPanVerified,
		})
	}

	return utils.SuccessResponse(c, 200, "No verification found", nil)
}

func VerifyPAN(c fiber.Ctx) error {
	type PANRequest struct {
		PAN  string `json:"pan"`
		Name string `json:"name"`
		DOB  string `json:"dob"`
	}
	var pr PANRequest
	if err := c.Bind().Body(&pr); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	if pr.PAN == "" || pr.Name == "" {
		return utils.ErrorResponse(c, 400, "PAN and Name are required")
	}

	userID := c.Locals("uid").(string)

	// Normalise PAN early so all comparisons use a consistent value
	panUpper := strings.ToUpper(strings.TrimSpace(pr.PAN))

	existing, _ := partnerRepo.FindByPAN(c.Context(), panUpper)
	if existing != nil && existing.UserID != userID {
		auditRepo.Log(c.Context(), &models.AuditLog{
			UserID:    userID,
			Action:    "PAN_VERIFICATION",
			Status:    "FAILED",
			Details:   fmt.Sprintf("PAN %s already linked to another account", panUpper),
			IPAddress: c.IP(),
		})
		return utils.ErrorResponse(c, 400, "This PAN is already registered under another account. Each PAN can only be associated with one account.")
	}

	// Check if user already verified this same PAN with same details - return stored data
	user, err := userRepo.FindByID(c.Context(), userID)
	if err != nil || user == nil {
		return utils.ErrorResponse(c, 500, "Failed to find user")
	}

	// Only shortcut if PAN, Name, AND DOB match what we have stored
	if user.IsPanVerified && user.PanNumber == panUpper && user.PanDetails != nil {
		storedName := strings.ToLower(strings.TrimSpace(user.PanDetails.RegisteredName))
		providedName := strings.ToLower(strings.TrimSpace(pr.Name))
		storedDOB := strings.TrimSpace(user.PanDetails.DOB)
		providedDOB := strings.TrimSpace(pr.DOB)

		if storedName == providedName && (providedDOB == "" || storedDOB == providedDOB) {
			return utils.SuccessResponse(c, 200, "PAN already verified", fiber.Map{
				"pan_verification": user.PanDetails,
				"pan":              user.PanNumber,
				"already_verified": true,
			})
		}
	}

	// Call real Cashfree PAN Lite API
	cfg := config.LoadConfig()
	if cfg.CashfreeClientID == "" || cfg.CashfreeSecret == "" {
		return utils.ErrorResponse(c, 500, "PAN verification service not configured. Please contact support.")
	}

	// PAN Lite API request — verification_id is mandatory
	type cashfreePANReq struct {
		VerificationID string `json:"verification_id"`
		PAN            string `json:"pan"`
		Name           string `json:"name"`
		DOB            string `json:"dob,omitempty"`
	}
	// PAN Lite response is flat (no nested data object)
	type cashfreePANResp struct {
		Status                   string      `json:"status"`
		Message                  interface{} `json:"message"`
		VerificationID           string      `json:"verification_id"`
		ReferenceID              int         `json:"reference_id"`
		PAN                      string      `json:"pan"`
		PANStatus                string      `json:"pan_status"`
		NameMatch                string      `json:"name_match"`
		DOBMatch                 string      `json:"dob_match"`
		DOB                      string      `json:"dob"`
		AadhaarSeedingStatus     string      `json:"aadhaar_seeding_status"`
		AadhaarSeedingStatusDesc string      `json:"aadhaar_seeding_status_desc"`
	}

	// panUpper already computed above; use UUID as verification_id
	cfReqBody, _ := json.Marshal(cashfreePANReq{
		VerificationID: utils.GenerateUUIDv7(),
		PAN:            panUpper,
		Name:           pr.Name,
		DOB:            pr.DOB,
	})

	httpReq, err := http.NewRequestWithContext(c.Context(), "POST", cfg.CashfreePANURL, bytes.NewReader(cfReqBody))
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to create PAN verification request")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-client-id", cfg.CashfreeClientID)
	httpReq.Header.Set("x-client-secret", cfg.CashfreeSecret)
	httpReq.Header.Set("x-api-version", "2023-08-01")

	httpClient := &http.Client{Timeout: 15 * time.Second}
	cfHTTPResp, err := httpClient.Do(httpReq)
	if err != nil {
		fmt.Printf(" Cashfree PAN API error: %v\n", err)
		return utils.ErrorResponse(c, 502, "PAN verification service temporarily unavailable. Please try again.")
	}
	defer cfHTTPResp.Body.Close()

	var cfResp cashfreePANResp
	if err := json.NewDecoder(cfHTTPResp.Body).Decode(&cfResp); err != nil {
		fmt.Printf(" Failed to decode Cashfree PAN response: %v\n", err)
		return utils.ErrorResponse(c, 502, "Invalid response from PAN verification service")
	}

	// Log raw response for debugging
	fmt.Printf(" Cashfree PAN Lite response: status=%s pan_status=%s name_match=%s dob_match=%s\n",
		cfResp.Status, cfResp.PANStatus, cfResp.NameMatch, cfResp.DOBMatch)

	// PAN Lite returns "VALID" on success (not "SUCCESS")
	if cfResp.Status != "VALID" {
		// Extract raw message for logging
		rawMsg := ""
		if cfResp.Message != nil {
			switch v := cfResp.Message.(type) {
			case string:
				rawMsg = v
			case map[string]interface{}:
				if b, err := json.Marshal(v); err == nil {
					rawMsg = string(b)
				}
			}
		}
		auditRepo.Log(c.Context(), &models.AuditLog{
			UserID:    userID,
			Action:    "PAN_VERIFICATION",
			Status:    "FAILED",
			Details:   fmt.Sprintf("Cashfree rejected PAN %s: status=%s raw=%s", panUpper, cfResp.Status, rawMsg),
			IPAddress: c.IP(),
		})
		// Give a clear, actionable message to the user
		return utils.ErrorResponse(c, 400, "PAN verification failed. Please check that your PAN number is correct and that your Name and Date of Birth exactly match your PAN card.")
	}

	// pan_status "E" = Existing (valid/active), "I" = Invalid
	if cfResp.PANStatus != "E" {
		auditRepo.Log(c.Context(), &models.AuditLog{
			UserID:    userID,
			Action:    "PAN_VERIFICATION",
			Status:    "FAILED",
			Details:   fmt.Sprintf("PAN %s invalid or inactive (pan_status: %s)", panUpper, cfResp.PANStatus),
			IPAddress: c.IP(),
		})
		return utils.ErrorResponse(c, 400, "Invalid PAN number. This PAN appears to be invalid or inactive. Please verify your PAN number and try again.")
	}

	// Use name_match returned by Cashfree directly
	nameMatch := cfResp.NameMatch
	if nameMatch == "" {
		nameMatch = "N"
	}

	// Reject if name does not match PAN records
	if nameMatch == "N" {
		auditRepo.Log(c.Context(), &models.AuditLog{
			UserID:    userID,
			Action:    "PAN_VERIFICATION",
			Status:    "FAILED",
			Details:   fmt.Sprintf("PAN %s name mismatch: provided=%s", panUpper, pr.Name),
			IPAddress: c.IP(),
		})
		return utils.ErrorResponse(c, 400, "Name mismatch. The name you entered does not match PAN records. Please enter your full name exactly as it appears on your PAN card.")
	}

	// Reject if DOB was provided and does not match PAN records
	dobMatch := cfResp.DOBMatch
	if pr.DOB != "" && dobMatch == "N" {
		auditRepo.Log(c.Context(), &models.AuditLog{
			UserID:    userID,
			Action:    "PAN_VERIFICATION",
			Status:    "FAILED",
			Details:   fmt.Sprintf("PAN %s dob mismatch: provided=%s", panUpper, pr.DOB),
			IPAddress: c.IP(),
		})
		return utils.ErrorResponse(c, 400, "Date of birth mismatch. The date of birth you entered does not match PAN records. Please enter your DOB exactly as registered on your PAN card.")
	}

	panVerification := models.PANVerification{
		Status:                   "VALID",
		ReferenceID:              cfResp.ReferenceID,
		VerificationID:           cfResp.VerificationID,
		RegisteredName:           pr.Name, // PAN Lite does not return registered name; store provided name
		NamePanCard:              "",
		NameProvided:             pr.Name,
		NameMatch:                nameMatch,
		PanStatus:                cfResp.PANStatus,
		DOB:                      cfResp.DOB,
		DOBMatch:                 cfResp.DOBMatch,
		Type:                     "",
		Gender:                   "",
		FirstName:                "",
		LastName:                 "",
		AadhaarSeedingStatus:     cfResp.AadhaarSeedingStatus,
		AadhaarSeedingStatusDesc: cfResp.AadhaarSeedingStatusDesc,
		VerifiedAt:               time.Now(),
	}

	// Store PAN details on user record
	user.PanNumber = panUpper
	user.PanDetails = &panVerification
	user.IsPanVerified = true
	user.UpdatedAt = time.Now()

	if err := userRepo.Update(c.Context(), user); err != nil {
		fmt.Printf(" Failed to store PAN details: %v\n", err)
		auditRepo.Log(c.Context(), &models.AuditLog{
			UserID:    userID,
			Action:    "PAN_VERIFICATION",
			Status:    "FAILED",
			Details:   fmt.Sprintf("Failed to update user record after successful PAN check: %v", err),
			IPAddress: c.IP(),
		})
		return utils.ErrorResponse(c, 500, "Failed to update user PAN status")
	}

	auditRepo.Log(c.Context(), &models.AuditLog{
		UserID:    userID,
		Action:    "PAN_VERIFICATION",
		Status:    "SUCCESS",
		Details:   fmt.Sprintf("PAN %s verified successfully (name_match=%s)", panUpper, nameMatch),
		IPAddress: c.IP(),
	})

	return utils.SuccessResponse(c, 200, "PAN verified successfully", fiber.Map{
		"pan_verification": panVerification,
		"pan":              panUpper,
		"already_verified": false,
	})
}

func GetGSTINsFromPAN(c fiber.Ctx) error {
	type PANRequest struct {
		PAN string `json:"pan"`
	}
	var pr PANRequest
	if err := c.Bind().Body(&pr); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	if pr.PAN == "" {
		return utils.ErrorResponse(c, 400, "PAN is required")
	}

	userID := c.Locals("uid").(string)

	// Call real Cashfree PAN-GSTIN API
	gstinCfg := config.LoadConfig()
	if gstinCfg.CashfreeClientID == "" || gstinCfg.CashfreeSecret == "" {
		return utils.ErrorResponse(c, 500, "GSTIN verification service not configured.")
	}

	type gstinReqBody struct {
		PAN string `json:"pan"`
	}
	type gstinEntry struct {
		GSTIN  string `json:"gstin"`
		Status string `json:"status"`
		State  string `json:"state"`
	}
	type gstinRespData struct {
		GSTINList []gstinEntry `json:"gstin_list"`
	}
	type cashfreeGSTINResp struct {
		Status         string        `json:"status"`
		Message        interface{}   `json:"message"`
		ReferenceID    int           `json:"reference_id"`
		VerificationID string        `json:"verification_id"`
		PAN            string        `json:"pan"`
		Data           gstinRespData `json:"data"`
	}

	panForGSTIN := strings.ToUpper(strings.TrimSpace(pr.PAN))
	gstinBodyBytes, _ := json.Marshal(gstinReqBody{PAN: panForGSTIN})

	gstinHTTPReq, err := http.NewRequestWithContext(c.Context(), "POST", gstinCfg.CashfreePANGSTINURL, bytes.NewReader(gstinBodyBytes))
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to create GSTIN request")
	}
	gstinHTTPReq.Header.Set("Content-Type", "application/json")
	gstinHTTPReq.Header.Set("x-client-id", gstinCfg.CashfreeClientID)
	gstinHTTPReq.Header.Set("x-client-secret", gstinCfg.CashfreeSecret)
	gstinHTTPReq.Header.Set("x-api-version", "2023-08-01")

	gstinClient := &http.Client{Timeout: 15 * time.Second}
	gstinHTTPResp, err := gstinClient.Do(gstinHTTPReq)
	if err != nil {
		fmt.Printf(" Cashfree GSTIN API error: %v\n", err)
		return utils.ErrorResponse(c, 502, "GSTIN verification service temporarily unavailable")
	}
	defer gstinHTTPResp.Body.Close()

	var gstinCFResp cashfreeGSTINResp
	if err := json.NewDecoder(gstinHTTPResp.Body).Decode(&gstinCFResp); err != nil {
		fmt.Printf(" Failed to decode Cashfree GSTIN response: %v\n", err)
		return utils.ErrorResponse(c, 502, "Invalid response from GSTIN verification service")
	}

	gstinList := []interface{}{}
	for _, g := range gstinCFResp.Data.GSTINList {
		gstinList = append(gstinList, map[string]interface{}{
			"gstin":  g.GSTIN,
			"status": g.Status,
			"state":  g.State,
		})
	}

	result := map[string]interface{}{
		"status":          gstinCFResp.Status,
		"reference_id":    gstinCFResp.ReferenceID,
		"verification_id": gstinCFResp.VerificationID,
		"pan":             panForGSTIN,
		"gstin_list":      gstinList,
	}

	auditRepo.Log(c.Context(), &models.AuditLog{
		UserID:    userID,
		Action:    "GSTIN_VERIFICATION",
		Status:    "SUCCESS",
		Details:   fmt.Sprintf("GSTIN lookup for PAN %s returned %d result(s)", panForGSTIN, len(gstinCFResp.Data.GSTINList)),
		IPAddress: c.IP(),
	})

	return utils.SuccessResponse(c, 200, "GSTIN mapping completed", result)
}

// GetPrefillData returns previously verified PAN details and profile data
// so when a user comes for a 2nd/3rd category (e.g., already verified for events, now wants play),
// the form is pre-filled with existing data. PAN is non-editable, other fields are editable.
func GetPrefillData(c fiber.Ctx) error {
	userID := c.Locals("uid").(string)

	user, err := userRepo.FindByID(c.Context(), userID)
	if err != nil || user == nil {
		return utils.ErrorResponse(c, 404, "User not found")
	}

	// Get all existing profiles for this user
	profiles, _ := partnerRepo.FindAllByUserID(c.Context(), userID)

	prefill := fiber.Map{
		"has_existing_verification": len(profiles) > 0,
		"is_pan_verified":           user.IsPanVerified,
		"pan_number":                user.PanNumber,
		"pan_details":               user.PanDetails,
		"existing_categories":       user.OrganizerCategories,
		"email":                     user.Email,
	}

	// If user has existing profiles, get the most recent one for bank/contact prefill
	if len(profiles) > 0 {
		latest := profiles[0]
		for _, p := range profiles {
			if p.UpdatedAt.After(latest.UpdatedAt) {
				latest = p
			}
		}
		prefill["bank_details"] = latest.BankDetails
		prefill["backup_contact"] = latest.BackupContact
		prefill["gst_details"] = latest.GSTDetails
		prefill["organization_name"] = latest.OrganizationDetails.PANName
		if user.PanDetails != nil {
			prefill["state"] = user.PanDetails.Address.State
		}
	}

	return utils.SuccessResponse(c, 200, "Prefill data fetched", prefill)
}
