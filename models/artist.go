package models

import "time"

type Artist struct {
	ID              string    `json:"id" firestore:"id"`
	Name            string    `json:"name" firestore:"name"`
	Role            string    `json:"role" firestore:"role"`
	ImageURL        string    `json:"image_url" firestore:"image_url"`
	Description     string    `json:"description" firestore:"description"`
	Genre           string    `json:"genre" firestore:"genre"`
	IsVerified      bool      `json:"is_verified" firestore:"is_verified"`
	Rating          float64   `json:"rating" firestore:"rating"`
	ReviewCount     int       `json:"review_count" firestore:"review_count"`
	EventsHosted    int       `json:"events_hosted" firestore:"events_hosted"`
	Location        string    `json:"location" firestore:"location"`
	ContactEmail    string    `json:"contact_email" firestore:"contact_email"`
	ContactPhone    string    `json:"contact_phone" firestore:"contact_phone"`
	ExperienceYears int       `json:"experience_years" firestore:"experience_years"`
	Specialties     []string  `json:"specialties" firestore:"specialties"`
	FollowerCount   int       `json:"follower_count" firestore:"follower_count"`
	SocialLinks     []string  `json:"social_links" firestore:"social_links"`
	CreatedAt       time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" firestore:"updated_at"`
}
