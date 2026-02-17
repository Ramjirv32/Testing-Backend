package models

import "time"

type PlayLocation struct {
	VenueName string  `json:"venue_name" firestore:"venue_name"`
	Address   string  `json:"address" firestore:"address"`
	City      string  `json:"city" firestore:"city"`
	State     string  `json:"state" firestore:"state"`
	Latitude  float64 `json:"latitude" firestore:"latitude"`
	Longitude float64 `json:"longitude" firestore:"longitude"`
	MapURL    string  `json:"map_url" firestore:"map_url"`
}

type PlayImages struct {
	Hero    string   `json:"hero" firestore:"hero"`
	Promo   string   `json:"promo" firestore:"promo"`
	Gallery []string `json:"gallery" firestore:"gallery"`
}

type PlayOption struct {
	Sport        string  `json:"sport" firestore:"sport"`
	CourtType    string  `json:"court_type" firestore:"court_type"`
	Surface      string  `json:"surface" firestore:"surface"`
	PricePerSlot float64 `json:"price_per_slot" firestore:"price_per_slot"`
}

type SlotSettings struct {
	SlotDurationMinutes   int    `json:"slot_duration_minutes" firestore:"slot_duration_minutes"`
	OpenTime              string `json:"open_time" firestore:"open_time"`
	CloseTime             string `json:"close_time" firestore:"close_time"`
	MaxDaysAdvanceBooking int    `json:"max_days_advance_booking" firestore:"max_days_advance_booking"`
}

type PlayFAQ struct {
	Question string `json:"question" firestore:"question"`
	Answer   string `json:"answer" firestore:"answer"`
}

type PlayVenue struct {
	ID                     string       `json:"id" firestore:"id"`
	OrganizerID            string       `json:"organizer_id" firestore:"organizer_id"`
	Name                   string       `json:"name" firestore:"name"`
	Slug                   string       `json:"slug" firestore:"slug"`
	Status                 string       `json:"status" firestore:"status"` // active, inactive, draft
	About                  string       `json:"about" firestore:"about"`
	ShortAbout             string       `json:"short_about" firestore:"short_about"`
	DurationPerSlotMinutes int          `json:"duration_per_slot_minutes" firestore:"duration_per_slot_minutes"`
	Location               PlayLocation `json:"location" firestore:"location"`
	Images                 PlayImages   `json:"images" firestore:"images"`
	PlayOptions            []PlayOption `json:"play_options" firestore:"play_options"`
	SlotSettings           SlotSettings `json:"slot_settings" firestore:"slot_settings"`
	FAQs                   []PlayFAQ    `json:"faqs" firestore:"faqs"`
	TermsAndConditions     []string     `json:"terms_and_conditions" firestore:"terms_and_conditions"`
	CreatedAt              time.Time    `json:"created_at" firestore:"created_at"`
	UpdatedAt              time.Time    `json:"updated_at" firestore:"updated_at"`
}
