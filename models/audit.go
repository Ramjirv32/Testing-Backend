package models

import "time"

type AuditLog struct {
	ID        string    `json:"id" firestore:"id"`
	UserID    string    `json:"user_id" firestore:"user_id"`
	Action    string    `json:"action" firestore:"action"` // e.g., "PAN_VERIFICATION", "VERIFICATION_SUBMISSION"
	Category  string    `json:"category" firestore:"category"`
	Status    string    `json:"status" firestore:"status"` // e.g., "SUCCESS", "FAILED"
	Details   string    `json:"details" firestore:"details"`
	IPAddress string    `json:"ip_address" firestore:"ip_address"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
}
