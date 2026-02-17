package models

import "time"

type Offer struct {
	ID          string    `json:"id" firestore:"id"`
	UserID      string    `json:"user_id" firestore:"user_id"` // Mapped to a specific user
	Code        string    `json:"code" firestore:"code"`
	Discount    string    `json:"discount" firestore:"discount"`
	Description string    `json:"description" firestore:"description"`
	Terms       string    `json:"terms" firestore:"terms"`
	ExpiryDate  time.Time `json:"expiry_date" firestore:"expiry_date"`
	IsActive    bool      `json:"is_active" firestore:"is_active"`
	CreatedAt   time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" firestore:"updated_at"`
}

type CreateOfferRequest struct {
	UserID      string    `json:"user_id" validate:"required"`
	Code        string    `json:"code" validate:"required"`
	Discount    string    `json:"discount" validate:"required"`
	Description string    `json:"description"`
	ExpiryDate  time.Time `json:"expiry_date"`
}
