package models

import "time"

type BookingStatus string

const (
	BookingPending   BookingStatus = "pending"
	BookingConfirmed BookingStatus = "confirmed"
	BookingCancelled BookingStatus = "cancelled"
	BookingCompleted BookingStatus = "completed"
)

type BookingType string

const (
	BookingTypePlay   BookingType = "play"
	BookingTypeDining BookingType = "dining"
	BookingTypeEvent  BookingType = "event"
)

type PlayBooking struct {
	ID                 string        `json:"id" firestore:"id"`
	SeqID              int64         `json:"seq_id" firestore:"seq_id"`
	UserID             string        `json:"user_id" firestore:"user_id"`
	VenueID            string        `json:"venue_id" firestore:"venue_id"`
	VenueName          string        `json:"venue_name" firestore:"venue_name"`
	Sport              string        `json:"sport" firestore:"sport"`
	Date               string        `json:"date" firestore:"date"`
	TimeSlot           string        `json:"time_slot" firestore:"time_slot"`
	PlayerName         string        `json:"player_name" firestore:"player_name"`
	Price              float64       `json:"price" firestore:"price"`
	BillingEmail       string        `json:"billing_email" firestore:"billing_email"`
	BillingState       string        `json:"billing_state" firestore:"billing_state"`
	BillingNationality string        `json:"billing_nationality" firestore:"billing_nationality"`
	OrganizerID        string        `json:"organizer_id" firestore:"organizer_id"`
	OrganizerEmail     string        `json:"organizer_email" firestore:"organizer_email"`
	Status             BookingStatus `json:"status" firestore:"status"`
	CreatedAt          time.Time     `json:"created_at" firestore:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at" firestore:"updated_at"`
}

type DiningBooking struct {
	ID             string        `json:"id" firestore:"id"`
	SeqID          int64         `json:"seq_id" firestore:"seq_id"`
	UserID         string        `json:"user_id" firestore:"user_id"`
	RestaurantID   string        `json:"restaurant_id" firestore:"restaurant_id"`
	RestaurantName string        `json:"restaurant_name" firestore:"restaurant_name"`
	Date           string        `json:"date" firestore:"date"`
	TimeSlot       string        `json:"time_slot" firestore:"time_slot"`
	GuestCount     int           `json:"guest_count" firestore:"guest_count"`
	GuestName      string        `json:"guest_name" firestore:"guest_name"`
	SpecialRequest string        `json:"special_request" firestore:"special_request"`
	OrganizerID    string        `json:"organizer_id" firestore:"organizer_id"`
	OrganizerEmail string        `json:"organizer_email" firestore:"organizer_email"`
	Status         BookingStatus `json:"status" firestore:"status"`
	CreatedAt      time.Time     `json:"created_at" firestore:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at" firestore:"updated_at"`
}

type CreatePlayBookingRequest struct {
	VenueID            string  `json:"venue_id" validate:"required"`
	VenueName          string  `json:"venue_name" validate:"required"`
	Sport              string  `json:"sport" validate:"required"`
	Date               string  `json:"date" validate:"required"`
	TimeSlot           string  `json:"time_slot" validate:"required"`
	PlayerName         string  `json:"player_name" validate:"required"`
	Price              float64 `json:"price"`
	BillingEmail       string  `json:"billing_email" validate:"required,email"`
	BillingState       string  `json:"billing_state" validate:"required"`
	BillingNationality string  `json:"billing_nationality" validate:"required"`
}

type CreateDiningBookingRequest struct {
	RestaurantID   string `json:"restaurant_id" validate:"required"`
	RestaurantName string `json:"restaurant_name" validate:"required"`
	Date           string `json:"date" validate:"required"`
	TimeSlot       string `json:"time_slot"`
	GuestCount     int    `json:"guest_count" validate:"required,min=1"`
	GuestName      string `json:"guest_name" validate:"required"`
	SpecialRequest string `json:"special_request"`
}

type EventBooking struct {
	ID             string        `json:"id" firestore:"id"`
	SeqID          int64         `json:"seq_id" firestore:"seq_id"`
	UserID         string        `json:"user_id" firestore:"user_id"`
	EventID        string        `json:"event_id" firestore:"event_id"`
	EventTitle     string        `json:"event_title" firestore:"event_title"`
	TicketType     string        `json:"ticket_type" firestore:"ticket_type"`
	SeatType       string        `json:"seat_type" firestore:"seat_type"`
	Quantity       int           `json:"quantity" firestore:"quantity"`
	UnitPrice      float64       `json:"unit_price" firestore:"unit_price"`
	TotalPrice     float64       `json:"total_price" firestore:"total_price"`
	GuestName      string        `json:"guest_name" firestore:"guest_name"`
	BillingEmail   string        `json:"billing_email" firestore:"billing_email"`
	BillingState   string        `json:"billing_state" firestore:"billing_state"`
	OrganizerID    string        `json:"organizer_id" firestore:"organizer_id"`
	OrganizerEmail string        `json:"organizer_email" firestore:"organizer_email"`
	Status         BookingStatus `json:"status" firestore:"status"`
	CreatedAt      time.Time     `json:"created_at" firestore:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at" firestore:"updated_at"`
}

type CreateEventBookingRequest struct {
	EventID      string  `json:"event_id" validate:"required"`
	EventTitle   string  `json:"event_title" validate:"required"`
	TicketType   string  `json:"ticket_type" validate:"required"`
	SeatType     string  `json:"seat_type"`
	Quantity     int     `json:"quantity" validate:"required,min=1"`
	UnitPrice    float64 `json:"unit_price" validate:"required"`
	GuestName    string  `json:"guest_name" validate:"required"`
	BillingEmail string  `json:"billing_email" validate:"required,email"`
	BillingState string  `json:"billing_state" validate:"required"`
}
