package models

import "time"

type Artist struct {
	ID          string    `json:"id" firestore:"id"`
	Name        string    `json:"name" firestore:"name"`
	Role        string    `json:"role" firestore:"role"`
	ImageURL    string    `json:"image_url" firestore:"image_url"`
	Description string    `json:"description" firestore:"description"`
	SocialLinks []string  `json:"social_links" firestore:"social_links"`
	CreatedAt   time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" firestore:"updated_at"`
}
