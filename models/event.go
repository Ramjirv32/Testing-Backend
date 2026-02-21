package models

import "time"

type EventVenue struct {
	Name      string  `json:"name" firestore:"name"`
	Address   string  `json:"address" firestore:"address"`
	City      string  `json:"city" firestore:"city"`
	State     string  `json:"state" firestore:"state"`
	Latitude  float64 `json:"latitude" firestore:"latitude"`
	Longitude float64 `json:"longitude" firestore:"longitude"`
	MapURL    string  `json:"map_url" firestore:"map_url"`
}

type EventImages struct {
	Hero    string   `json:"hero" firestore:"hero"`
	Poster  string   `json:"poster" firestore:"poster"`
	Gallery []string `json:"gallery" firestore:"gallery"`
}

type EventArtist struct {
	Name            string   `json:"name" firestore:"name"`
	Role            string   `json:"role" firestore:"role"`
	ImageURL        string   `json:"image_url" firestore:"image_url"`
	Description     string   `json:"description" firestore:"description"`
	Genre           string   `json:"genre" firestore:"genre"`
	IsVerified      bool     `json:"is_verified" firestore:"is_verified"`
	Rating          float64  `json:"rating" firestore:"rating"`
	ReviewCount     int      `json:"review_count" firestore:"review_count"`
	EventsHosted    int      `json:"events_hosted" firestore:"events_hosted"`
	Location        string   `json:"location" firestore:"location"`
	ContactEmail    string   `json:"contact_email" firestore:"contact_email"`
	ContactPhone    string   `json:"contact_phone" firestore:"contact_phone"`
	ExperienceYears int      `json:"experience_years" firestore:"experience_years"`
	Specialties     []string `json:"specialties" firestore:"specialties"`
	FollowerCount   int      `json:"follower_count" firestore:"follower_count"`
	SocialLinks     []string `json:"social_links" firestore:"social_links"`
}

type EventTicket struct {
	TicketType        string  `json:"ticket_type" firestore:"ticket_type"`
	SeatType          string  `json:"seat_type" firestore:"seat_type"`
	Price             float64 `json:"price" firestore:"price"`
	TotalQuantity     int     `json:"total_quantity" firestore:"total_quantity"`
	AvailableQuantity int     `json:"available_quantity" firestore:"available_quantity"`
}

type EventFAQ struct {
	Question string `json:"question" firestore:"question"`
	Answer   string `json:"answer" firestore:"answer"`
}

type Event struct {
	ID                 string        `json:"id" firestore:"id"`
	OrganizerID        string        `json:"organizer_id" firestore:"organizer_id"`
	Title              string        `json:"title" firestore:"title"`
	Slug               string        `json:"slug" firestore:"slug"`
	Category           string        `json:"category" firestore:"category"`
	Language           string        `json:"language" firestore:"language"`
	DurationMinutes    int           `json:"duration_minutes" firestore:"duration_minutes"`
	AgeLimit           string        `json:"age_limit" firestore:"age_limit"`
	Description        string        `json:"description" firestore:"description"`
	ShortDescription   string        `json:"short_description" firestore:"short_description"`
	StartDatetime      time.Time     `json:"start_datetime" firestore:"start_datetime"`
	EndDatetime        time.Time     `json:"end_datetime" firestore:"end_datetime"`
	PriceStart         float64       `json:"price_start" firestore:"price_start"`
	Status             string        `json:"status" firestore:"status"`
	RejectionReason    string        `json:"rejection_reason" firestore:"rejection_reason"`
	Venue              EventVenue    `json:"venue" firestore:"venue"`
	Images             EventImages   `json:"images" firestore:"images"`
	Artists            []EventArtist `json:"artists" firestore:"artists"`
	Tickets            []EventTicket `json:"tickets" firestore:"tickets"`
	FAQs               []EventFAQ    `json:"faqs" firestore:"faqs"`
	TermsAndConditions []string      `json:"terms_and_conditions" firestore:"terms_and_conditions"`
	CreatedAt          time.Time     `json:"created_at" firestore:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at" firestore:"updated_at"`
}
