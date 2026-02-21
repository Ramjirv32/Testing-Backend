package pass

import (
"context"
"fmt"
"log"
"time"

"backend/config"
"backend/utils"

"github.com/gofiber/fiber/v3"
)

type ActivatePassRequest struct {
Name          string `json:"name"`
Email         string `json:"email"`
Phone         string `json:"phone"`
Address       string `json:"address"`
State         string `json:"state"`
District      string `json:"district"`
Country       string `json:"country"`
PlanName      string `json:"planName"`
PurchaseDate  string `json:"purchaseDate"`
ExpiryDate    string `json:"expiryDate"`
Amount        int    `json:"amount"`
PaymentMethod string `json:"paymentMethod"`
}

// ActivatePass validates eligibility, creates the pass record in Firestore, and
// returns the new pass ID. Requires an authenticated user (middleware.Auth).
func ActivatePass(c fiber.Ctx) error {
userID, ok := c.Locals("uid").(string)
if !ok || userID == "" {
return utils.ErrorResponse(c, 401, "Unauthorized")
}

var req ActivatePassRequest
if err := c.Bind().Body(&req); err != nil {
return utils.ErrorResponse(c, 400, "Invalid request body")
}

if req.Name == "" || (req.Email == "" && req.Phone == "") {
return utils.ErrorResponse(c, 400, "name and at least one of email/phone are required")
}
if req.PurchaseDate == "" || req.ExpiryDate == "" {
return utils.ErrorResponse(c, 400, "purchaseDate and expiryDate are required")
}

client := config.FirestoreClient
if client == nil {
return utils.ErrorResponse(c, 500, "Database service unavailable")
}

// Check that no active pass already exists for this email or phone
result, err := utils.CheckPassEligibility(c.Context(), req.Email, req.Phone, client)
if err != nil {
log.Printf("[pass] eligibility check error for user %s: %v", userID, err)
return utils.ErrorResponse(c, 500, "Eligibility check failed")
}
if !result.Eligible {
return utils.ErrorResponse(c, 409, fmt.Sprintf(
"An active pass already exists (matched by %s, expires %s)",
result.MatchedBy, result.ExpiryDate,
))
}

// Fire-and-forget: soft-delete any stale expired duplicates
go cleanupExpiredPasses(req.Email, req.Phone)

passID := utils.GenerateUUIDv7()

doc := map[string]interface{}{
"id":            passID,
"userId":        userID,
"name":          req.Name,
"email":         req.Email,
"phone":         req.Phone,
"address":       req.Address,
"state":         req.State,
"district":      req.District,
"country":       req.Country,
"planName":      req.PlanName,
"purchaseDate":  req.PurchaseDate,
"expiryDate":    req.ExpiryDate,
"amount":        req.Amount,
"paymentMethod": req.PaymentMethod,
"paymentStatus": "confirmed",
"status":        "active",
"createdAt":     time.Now().UTC().Format(time.RFC3339),
}

if _, err := client.Collection("ticpin_pass_users").Doc(passID).Set(c.Context(), doc); err != nil {
log.Printf("[pass] create failed for user %s: %v", userID, err)
return utils.ErrorResponse(c, 500, "Failed to create pass")
}

return utils.SuccessResponse(c, 201, "Pass activated successfully", fiber.Map{
"passId":       passID,
"expiryDate":   req.ExpiryDate,
"purchaseDate": req.PurchaseDate,
"status":       "active",
})
}

// cleanupExpiredPasses deletes expired/stale pass records for the given identifiers.
// Runs as a goroutine — errors are only logged.
func cleanupExpiredPasses(email, phone string) {
client := config.FirestoreClient
if client == nil {
return
}

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

now := time.Now()
coll := client.Collection("ticpin_pass_users")

type candidate struct{ field, value string }
checks := []candidate{}
if email != "" {
checks = append(checks, candidate{"email", email})
}
if phone != "" {
checks = append(checks, candidate{"phone", phone})
}

seen := map[string]bool{}
for _, ch := range checks {
iter := coll.Where(ch.field, "==", ch.value).Documents(ctx)
for {
doc, err := iter.Next()
if err != nil {
break
}
if seen[doc.Ref.ID] {
continue
}
seen[doc.Ref.ID] = true

expStr, _ := doc.Data()["expiryDate"].(string)
if expStr == "" {
continue
}
exp, err := time.Parse(time.RFC3339, expStr)
if err != nil {
exp, err = time.Parse("2006-01-02", expStr)
if err != nil {
continue
}
}
if exp.Before(now) {
if _, delErr := doc.Ref.Delete(ctx); delErr != nil {
log.Printf("[pass] cleanup delete failed for %s: %v", doc.Ref.ID, delErr)
}
}
}
}
}
