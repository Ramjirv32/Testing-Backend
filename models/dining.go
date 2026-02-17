package models

import "time"

type DiningContact struct {
	Phone string `json:"phone" firestore:"phone"`
	Email string `json:"email" firestore:"email"`
}

type DiningLocation struct {
	VenueName string  `json:"venue_name" firestore:"venue_name"`
	Address   string  `json:"address" firestore:"address"`
	City      string  `json:"city" firestore:"city"`
	State     string  `json:"state" firestore:"state"`
	Latitude  float64 `json:"latitude" firestore:"latitude"`
	Longitude float64 `json:"longitude" firestore:"longitude"`
	MapURL    string  `json:"map_url" firestore:"map_url"`
}

type DiningImages struct {
	Hero    string   `json:"hero" firestore:"hero"`
	Gallery []string `json:"gallery" firestore:"gallery"`
}

type DiningOffer struct {
	Title       string `json:"title" firestore:"title"`
	Code        string `json:"code" firestore:"code"`
	Description string `json:"description" firestore:"description"`
}

type SeatingType struct {
	Type             string `json:"type" firestore:"type"`
	TotalTables      int    `json:"total_tables" firestore:"total_tables"`
	AvailableTables  int    `json:"available_tables" firestore:"available_tables"`
	CapacityPerTable int    `json:"capacity_per_table" firestore:"capacity_per_table"`
}

type BookingSettings struct {
	AdvanceBookingDays           int      `json:"advance_booking_days" firestore:"advance_booking_days"`
	TimeSlots                    []string `json:"time_slots" firestore:"time_slots"`
	AverageDiningDurationMinutes int      `json:"average_dining_duration_minutes" firestore:"average_dining_duration_minutes"`
}

type DiningFAQ struct {
	Question string `json:"question" firestore:"question"`
	Answer   string `json:"answer" firestore:"answer"`
}

type DiningVenue struct {
	ID                 string          `json:"id" firestore:"id"`
	OrganizerID        string          `json:"organizer_id" firestore:"organizer_id"`
	Name               string          `json:"name" firestore:"name"`
	Slug               string          `json:"slug" firestore:"slug"`
	Status             string          `json:"status" firestore:"status"` // active, inactive, draft
	Description        string          `json:"description" firestore:"description"`
	ShortDescription   string          `json:"short_description" firestore:"short_description"`
	Rating             float64         `json:"rating" firestore:"rating"`
	IsOpen             bool            `json:"is_open" firestore:"is_open"`
	OpeningTime        string          `json:"opening_time" firestore:"opening_time"`
	ClosingTime        string          `json:"closing_time" firestore:"closing_time"`
	Contact            DiningContact   `json:"contact" firestore:"contact"`
	Location           DiningLocation  `json:"location" firestore:"location"`
	Images             DiningImages    `json:"images" firestore:"images"`
	MenuImages         []string        `json:"menu_images" firestore:"menu_images"`
	Offers             []DiningOffer   `json:"offers" firestore:"offers"`
	Facilities         []string        `json:"facilities" firestore:"facilities"`
	SeatingTypes       []SeatingType   `json:"seating_types" firestore:"seating_types"`
	BookingSettings    BookingSettings `json:"booking_settings" firestore:"booking_settings"`
	FAQs               []DiningFAQ     `json:"faqs" firestore:"faqs"`
	TermsAndConditions []string        `json:"terms_and_conditions" firestore:"terms_and_conditions"`
	CreatedAt          time.Time       `json:"created_at" firestore:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at" firestore:"updated_at"`
}
