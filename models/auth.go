package models

import "time"

type LoginRequest struct {
	Phone         string `json:"phone"`
	OTP           string `json:"otp"`
	FirebaseToken string `json:"firebase_token"`
}

type SendOTPRequest struct {
	Phone string `json:"phone" validate:"required"`
}

type LoginResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

type Session struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	UserID    string    `json:"user_id" bson:"user_id"`
	Token     string    `json:"token" bson:"token"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
