package models

import "time"

type Venue struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	SeqID       int64     `json:"seq_id" bson:"seq_id"`
	Name        string    `json:"name" bson:"name"`
	Location    string    `json:"location" bson:"location"`
	City        string    `json:"city" bson:"city"`
	Sports      []string  `json:"sports" bson:"sports"`
	Image       string    `json:"image" bson:"image"`
	Description string    `json:"description" bson:"description"`
	PricePerHr  float64   `json:"price_per_hr" bson:"price_per_hr"`
	Rating      float64   `json:"rating" bson:"rating"`
	IsActive    bool      `json:"is_active" bson:"is_active"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

type Restaurant struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	SeqID       int64     `json:"seq_id" bson:"seq_id"`
	Name        string    `json:"name" bson:"name"`
	Location    string    `json:"location" bson:"location"`
	City        string    `json:"city" bson:"city"`
	Cuisine     []string  `json:"cuisine" bson:"cuisine"`
	Image       string    `json:"image" bson:"image"`
	Description string    `json:"description" bson:"description"`
	Rating      float64   `json:"rating" bson:"rating"`
	PriceRange  string    `json:"price_range" bson:"price_range"`
	IsActive    bool      `json:"is_active" bson:"is_active"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}
