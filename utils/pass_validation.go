package utils

import (
"context"
"fmt"
"log"
"sort"
"time"

"cloud.google.com/go/firestore"
"google.golang.org/api/iterator"
)

// PassRecord represents a Ticpin Pass record from Firestore
type PassRecord struct {
ID           string
Email        string
Phone        string
UserID       string
PurchaseDate time.Time
ExpiryDate   time.Time
// "active" or "expired"
Status string
}

// PassEligibilityResult is the result of a uniqueness/eligibility check
type PassEligibilityResult struct {
Eligible       bool   `json:"eligible"`
HasActivePass  bool   `json:"has_active_pass"`
ExistingPassID string `json:"existing_pass_id,omitempty"`
ExpiryDate     string `json:"expiry_date,omitempty"`
MatchedBy      string `json:"matched_by,omitempty"` // "email" | "phone" | ""
Reason         string `json:"reason,omitempty"`
}

// CheckPassEligibility checks uniqueness by BOTH email and phone.
// Returns whether the user is eligible (no active pass) and details about any existing pass.
func CheckPassEligibility(ctx context.Context, email, phone string, client *firestore.Client) (*PassEligibilityResult, error) {
now := time.Now()

// Check by email first
if email != "" {
pass, err := queryActivePass(ctx, client, "email", email, now)
if err != nil {
return nil, err
}
if pass != nil {
return &PassEligibilityResult{
Eligible:       false,
HasActivePass:  true,
ExistingPassID: pass.ID,
ExpiryDate:     pass.ExpiryDate.Format("02/01/2006"),
MatchedBy:      "email",
Reason: fmt.Sprintf(
"You already have an active Ticpin Pass valid until %s. You can renew it after expiry.",
pass.ExpiryDate.Format("02/01/2006"),
),
}, nil
}
}

// Check by phone number
if phone != "" {
pass, err := queryActivePass(ctx, client, "phone", phone, now)
if err != nil {
return nil, err
}
if pass != nil {
return &PassEligibilityResult{
Eligible:       false,
HasActivePass:  true,
ExistingPassID: pass.ID,
ExpiryDate:     pass.ExpiryDate.Format("02/01/2006"),
MatchedBy:      "phone",
Reason: fmt.Sprintf(
"This phone number already has an active Ticpin Pass valid until %s. You can renew it after expiry.",
pass.ExpiryDate.Format("02/01/2006"),
),
}, nil
}
}

return &PassEligibilityResult{Eligible: true}, nil
}

// CheckUserHasActivePass checks by email only. Kept for backward compatibility.
func CheckUserHasActivePass(ctx context.Context, email string, client *firestore.Client) (bool, string, error) {
if email == "" {
return false, "", fmt.Errorf("email is required")
}
pass, err := queryActivePass(ctx, client, "email", email, time.Now())
if err != nil {
return false, "", err
}
if pass != nil {
return true, pass.ID, nil
}
return false, "", nil
}

// GetUserActivePass retrieves the user's active pass by email (or most recent expired).
func GetUserActivePass(ctx context.Context, email string, client *firestore.Client) (*PassRecord, error) {
if email == "" {
return nil, fmt.Errorf("email is required")
}

iter := client.Collection("ticpin_pass_users").Where("email", "==", email).Documents(ctx)
defer iter.Stop()

now := time.Now()
var mostRecent *PassRecord

for {
doc, err := iter.Next()
if err == iterator.Done {
break
}
if err != nil {
return nil, err
}

pass, ok := parsePassDoc(doc)
if !ok {
continue
}

if pass.ExpiryDate.After(now) {
pass.Status = "active"
return pass, nil
}
if mostRecent == nil || pass.PurchaseDate.After(mostRecent.PurchaseDate) {
pass.Status = "expired"
mostRecent = pass
}
}

return mostRecent, nil
}

// CleanupDuplicatePasses removes excess duplicate records.
// Keeps 1 active + 1 most recent expired; deletes the rest.
func CleanupDuplicatePasses(ctx context.Context, email string, client *firestore.Client) (int, error) {
if email == "" {
return 0, fmt.Errorf("email is required")
}

iter := client.Collection("ticpin_pass_users").Where("email", "==", email).Documents(ctx)
defer iter.Stop()

now := time.Now()
var active []*PassRecord
var expired []*PassRecord

for {
doc, err := iter.Next()
if err == iterator.Done {
break
}
if err != nil {
return 0, err
}

pass, ok := parsePassDoc(doc)
if !ok {
continue
}

if pass.ExpiryDate.After(now) {
pass.Status = "active"
active = append(active, pass)
} else {
pass.Status = "expired"
expired = append(expired, pass)
}
}

// Newest first
sort.Slice(active, func(i, j int) bool {
return active[i].PurchaseDate.After(active[j].PurchaseDate)
})
sort.Slice(expired, func(i, j int) bool {
return expired[i].PurchaseDate.After(expired[j].PurchaseDate)
})

deleted := 0
for i := 1; i < len(active); i++ {
if _, err := client.Collection("ticpin_pass_users").Doc(active[i].ID).Delete(ctx); err != nil {
log.Printf("Error deleting duplicate active pass %s: %v", active[i].ID, err)
continue
}
deleted++
}
for i := 1; i < len(expired); i++ {
if _, err := client.Collection("ticpin_pass_users").Doc(expired[i].ID).Delete(ctx); err != nil {
log.Printf("Error deleting old expired pass %s: %v", expired[i].ID, err)
continue
}
deleted++
}

if deleted > 0 {
log.Printf("Cleaned up %d duplicate pass records for %s", deleted, email)
}
return deleted, nil
}

// --- internal helpers ---

func queryActivePass(ctx context.Context, client *firestore.Client, field, value string, now time.Time) (*PassRecord, error) {
iter := client.Collection("ticpin_pass_users").Where(field, "==", value).Documents(ctx)
defer iter.Stop()

for {
doc, err := iter.Next()
if err == iterator.Done {
break
}
if err != nil {
log.Printf("Error querying passes by %s=%s: %v", field, value, err)
return nil, err
}

pass, ok := parsePassDoc(doc)
if !ok {
continue
}
if pass.ExpiryDate.After(now) {
pass.Status = "active"
return pass, nil
}
}
return nil, nil
}

func parsePassDoc(doc *firestore.DocumentSnapshot) (*PassRecord, bool) {
data := doc.Data()

expiryDate, expiryOk := data["expiryDate"].(time.Time)
purchaseDate, purchaseOk := data["purchaseDate"].(time.Time)
if !expiryOk || !purchaseOk {
log.Printf("Warning: Could not parse dates for pass doc %s", doc.Ref.ID)
return nil, false
}

userID, _ := data["userId"].(string)
email, _ := data["email"].(string)
phone, _ := data["phone"].(string)

return &PassRecord{
ID:           doc.Ref.ID,
Email:        email,
Phone:        phone,
UserID:       userID,
PurchaseDate: purchaseDate,
ExpiryDate:   expiryDate,
}, true
}
