package models

import "time"

type User struct {
	ID                string    `json:"id" firestore:"id"`
	FirebaseUID       string    `json:"firebase_uid" firestore:"firebase_uid"`
	SeqID             int64     `json:"seq_id" firestore:"seq_id"`
	Email             string    `json:"email" firestore:"email"`
	Password          string    `json:"-" firestore:"password"`
	Name              string    `json:"name" firestore:"name"`
	Phone             string    `json:"phone" firestore:"phone"`
	Avatar            string    `json:"avatar" firestore:"avatar"`
	IsEmailVerified   bool      `json:"is_email_verified" firestore:"is_email_verified"`
	IsOrganizer       bool      `json:"is_organizer" firestore:"is_organizer"`
	OrganizerCategory string    `json:"organizer_category" firestore:"organizer_category"` // 'play', 'dining', 'event'
	CreatedAt         time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" firestore:"updated_at"`
}

type CreateUserRequest struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required"`
	Phone string `json:"phone"`
}

type UpdateUserRequest struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Avatar string `json:"avatar"`
}

type SendEmailOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required"`
	Token string `json:"token"` // Optional token if already logged in via phone
}

type OrganizerRegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
}

type OrganizerLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type GoogleLoginRequest struct {
	Email   string `json:"email" validate:"required,email"`
	Name    string `json:"name"`
	IDToken string `json:"id_token"` // Firebase ID Token
}
