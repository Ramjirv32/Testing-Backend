package models

import "time"

type User struct {
	ID                  string           `json:"id" firestore:"id"`
	FirebaseUID         string           `json:"firebase_uid" firestore:"firebase_uid"`
	SeqID               int64            `json:"seq_id" firestore:"seq_id"`
	Email               string           `json:"email" firestore:"email"`
	Password            string           `json:"-" firestore:"password"`
	Name                string           `json:"name" firestore:"name"`
	Phone               string           `json:"phone" firestore:"phone"`
	Avatar              string           `json:"avatar" firestore:"avatar"`
	IsEmailVerified     bool             `json:"is_email_verified" firestore:"is_email_verified"`
	IsOrganizer         bool             `json:"is_organizer" firestore:"is_organizer"`
	IsAdmin             bool             `json:"is_admin" firestore:"is_admin"`
	IsPanVerified       bool             `json:"is_pan_verified" firestore:"is_pan_verified"`
	OrganizerCategory   string           `json:"organizer_category" firestore:"organizer_category"`
	OrganizerCategories []string         `json:"organizer_categories" firestore:"organizer_categories"`
	PanNumber           string           `json:"pan_number,omitempty" firestore:"pan_number,omitempty"`
	PanDetails          *PANVerification `json:"pan_details,omitempty" firestore:"pan_details,omitempty"`
	Address             string           `json:"address,omitempty" firestore:"address"`
	State               string           `json:"state,omitempty" firestore:"state"`
	District            string           `json:"district,omitempty" firestore:"district"`
	Country             string           `json:"country,omitempty" firestore:"country"`
	CreatedAt           time.Time        `json:"created_at" firestore:"created_at"`
	UpdatedAt           time.Time        `json:"updated_at" firestore:"updated_at"`
}

// HasOrganizerCategory checks if the user is approved as an organizer for the given category.
// It checks both the OrganizerCategories slice (new) and the OrganizerCategory string (legacy).
func (u *User) HasOrganizerCategory(category string) bool {
	for _, c := range u.OrganizerCategories {
		if c == category {
			return true
		}
	}
	return u.OrganizerCategory == category
}

type CreateUserRequest struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required"`
	Phone string `json:"phone"`
}

type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Address  string `json:"address"`
	State    string `json:"state"`
	District string `json:"district"`
	Country  string `json:"country"`
}

type SendEmailOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required"`
	Token string `json:"token"`
}

type OrganizerRegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
	Phone    string `json:"phone"`
}

type OrganizerLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type GoogleLoginRequest struct {
	Email   string `json:"email" validate:"required,email"`
	Name    string `json:"name"`
	IDToken string `json:"id_token"`
}
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Email    string `json:"email" validate:"required,email"`
	OTP      string `json:"otp" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}
