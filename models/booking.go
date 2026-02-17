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
	ID                 string        `json:"id" bson:"_id,omitempty"`
	SeqID              int64         `json:"seq_id" bson:"seq_id"`
	UserID             string        `json:"user_id" bson:"user_id"`
	VenueID            string        `json:"venue_id" bson:"venue_id"`
	VenueName          string        `json:"venue_name" bson:"venue_name"`
	Sport              string        `json:"sport" bson:"sport"`
	Date               string        `json:"date" bson:"date"`
	TimeSlot           string        `json:"time_slot" bson:"time_slot"`
	PlayerName         string        `json:"player_name" bson:"player_name"`
	Price              float64       `json:"price" bson:"price"`
	BillingEmail       string        `json:"billing_email" bson:"billing_email"`
	BillingState       string        `json:"billing_state" bson:"billing_state"`
	BillingNationality string        `json:"billing_nationality" bson:"billing_nationality"`
	Status             BookingStatus `json:"status" bson:"status"`
	CreatedAt          time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at" bson:"updated_at"`
}

type DiningBooking struct {
	ID             string        `json:"id" bson:"_id,omitempty"`
	SeqID          int64         `json:"seq_id" bson:"seq_id"`
	UserID         string        `json:"user_id" bson:"user_id"`
	RestaurantID   string        `json:"restaurant_id" bson:"restaurant_id"`
	RestaurantName string        `json:"restaurant_name" bson:"restaurant_name"`
	Date           string        `json:"date" bson:"date"`
	TimeSlot       string        `json:"time_slot" bson:"time_slot"`
	GuestCount     int           `json:"guest_count" bson:"guest_count"`
	GuestName      string        `json:"guest_name" bson:"guest_name"`
	SpecialRequest string        `json:"special_request" bson:"special_request"`
	Status         BookingStatus `json:"status" bson:"status"`
	CreatedAt      time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at" bson:"updated_at"`
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
