package main

import (
	"backend/config"
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
)

// Seed data structures
type PlayVenueData struct {
	ID                     string
	OrganizerID            string
	Name                   string
	Slug                   string
	Status                 string
	About                  string
	ShortAbout             string
	DurationPerSlotMinutes int
	Location               map[string]interface{}
	Images                 map[string]interface{}
	PlayOptions            []map[string]interface{}
	SlotSettings           map[string]interface{}
	FAQs                   []map[string]interface{}
	TermsAndConditions     []string
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

type DiningVenueData struct {
	ID                 string
	OrganizerID        string
	Name               string
	Slug               string
	Status             string
	Description        string
	ShortDescription   string
	Rating             float64
	IsOpen             bool
	OpeningTime        string
	ClosingTime        string
	Contact            map[string]interface{}
	Location           map[string]interface{}
	Images             map[string]interface{}
	MenuImages         []string
	PopularDishes      []string
	Offers             []map[string]interface{}
	Facilities         []string
	SeatingTypes       []map[string]interface{}
	BookingSettings    map[string]interface{}
	FAQs               []map[string]interface{}
	TermsAndConditions []string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type EventData struct {
	ID                 string
	OrganizerID        string
	Title              string
	Slug               string
	Category           string
	Language           string
	DurationMinutes    int
	AgeLimit           string
	Description        string
	ShortDescription   string
	StartDatetime      time.Time
	EndDatetime        time.Time
	PriceStart         float64
	Status             string
	Venue              map[string]interface{}
	Images             map[string]interface{}
	Artists            []map[string]interface{}
	Tickets            []map[string]interface{}
	FAQs               []map[string]interface{}
	TermsAndConditions []string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func main() {
	ctx := context.Background()

	cfg := config.LoadConfig()

	fmt.Println(" Initializing Firebase Services...")
	config.InitFirebase(cfg)

	client := config.FirestoreClient
	if client == nil {
		log.Fatal(" Firestore client is nil")
	}

	fmt.Println(" Starting database seeding with real data...")

	// Create organizer user first
	organizerID := uuid.New().String()
	if err := seedOrganizer(ctx, client, organizerID); err != nil {
		log.Printf(" Failed to seed organizer: %v\n", err)
	} else {
		fmt.Printf(" Created organizer with ID: %s\n", organizerID)
	}

	// Seed Play Venues
	fmt.Println("\n Seeding Play Venues...")
	if err := seedPlayVenues(ctx, client, organizerID); err != nil {
		log.Printf(" Failed to seed play venues: %v\n", err)
	} else {
		fmt.Println(" Play venues seeded successfully")
	}

	// Seed Dining Venues
	fmt.Println("\n Seeding Dining Venues...")
	if err := seedDiningVenues(ctx, client, organizerID); err != nil {
		log.Printf(" Failed to seed dining venues: %v\n", err)
	} else {
		fmt.Println(" Dining venues seeded successfully")
	}

	// Seed Events
	fmt.Println("\n Seeding Events...")
	if err := seedEvents(ctx, client, organizerID); err != nil {
		log.Printf(" Failed to seed events: %v\n", err)
	} else {
		fmt.Println(" Events seeded successfully")
	}

	fmt.Println("\n Database seeding complete!")
}

func seedOrganizer(ctx context.Context, client *firestore.Client, organizerID string) error {
	organizer := map[string]interface{}{
		"id":                  organizerID,
		"email":               "organizer@ticpin.com",
		"password":            "hashed_password_here",
		"fullName":            "TicPin Organizer",
		"phone":               "+91-9999999999",
		"status":              "active",
		"isDeletable":         false,
		"userType":            "organizer",
		"organizerCategories": []string{"play", "dining", "event"},
		"createdAt":           time.Now(),
		"updatedAt":           time.Now(),
	}

	_, err := client.Collection("users").Doc(organizerID).Set(ctx, organizer)
	return err
}

func seedPlayVenues(ctx context.Context, client *firestore.Client, organizerID string) error {
	playVenues := []PlayVenueData{
		{
			ID:                     uuid.New().String(),
			OrganizerID:            organizerID,
			Name:                   "Elite Badminton Club",
			Slug:                   "elite-badminton-club",
			Status:                 "active",
			About:                  "Premium badminton facility with 8 international standard courts. Professional coaching available.",
			ShortAbout:             "Premier badminton courts with professional setup",
			DurationPerSlotMinutes: 60,
			Location: map[string]interface{}{
				"venue_name": "Elite Badminton Club",
				"address":    "Plot 123, Sports Avenue, Tech Park",
				"city":       "Bangalore",
				"state":      "Karnataka",
				"latitude":   12.9716,
				"longitude":  77.5946,
				"map_url":    "https://maps.google.com/?q=12.9716,77.5946",
			},
			Images: map[string]interface{}{
				"hero":    "https://images.unsplash.com/photo-1617622614882-09fdf9106556?w=1200",
				"promo":   "https://images.unsplash.com/photo-1617622614882-09fdf9106556?w=600",
				"gallery": []string{"https://images.unsplash.com/photo-1530549387789-4c1017266635?w=500"},
			},
			PlayOptions: []map[string]interface{}{
				{
					"sport":          "Badminton",
					"court_type":     "Indoor",
					"surface":        "Wooden",
					"price_per_slot": 500.0,
				},
			},
			SlotSettings: map[string]interface{}{
				"slot_duration_minutes":    60,
				"open_time":                "06:00 AM",
				"close_time":               "10:00 PM",
				"max_days_advance_booking": 30,
			},
			FAQs: []map[string]interface{}{
				{
					"question": "Do you provide shuttles?",
					"answer":   "Yes, we provide premium shuttles at ₹50 per piece",
				},
				{
					"question": "Are instructors available?",
					"answer":   "Yes, professional coaches available at ₹500/hour",
				},
			},
			TermsAndConditions: []string{
				"Cancellation must be done 24 hours in advance",
				"Footwear is mandatory on courts",
				"No outside food allowed",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:                     uuid.New().String(),
			OrganizerID:            organizerID,
			Name:                   "Championship Tennis Arena",
			Slug:                   "championship-tennis-arena",
			Status:                 "active",
			About:                  "World-class tennis courts with grass, clay, and hard courts. AITA certified.",
			ShortAbout:             "Professional tennis facility with multiple court types",
			DurationPerSlotMinutes: 60,
			Location: map[string]interface{}{
				"venue_name": "Championship Tennis Arena",
				"address":    "45, Sports Complex Road, Koramangala",
				"city":       "Bangalore",
				"state":      "Karnataka",
				"latitude":   12.9352,
				"longitude":  77.6365,
				"map_url":    "https://maps.google.com/?q=12.9352,77.6365",
			},
			Images: map[string]interface{}{
				"hero":    "https://images.unsplash.com/photo-1554224311-beee415c201f?w=1200",
				"promo":   "https://images.unsplash.com/photo-1554224311-beee415c201f?w=600",
				"gallery": []string{"https://images.unsplash.com/photo-1461644437690-f90cc0c91efb?w=500"},
			},
			PlayOptions: []map[string]interface{}{
				{
					"sport":          "Tennis",
					"court_type":     "Hard Court",
					"surface":        "Synthetic",
					"price_per_slot": 800.0,
				},
				{
					"sport":          "Tennis",
					"court_type":     "Clay Court",
					"surface":        "Clay",
					"price_per_slot": 1000.0,
				},
			},
			SlotSettings: map[string]interface{}{
				"slot_duration_minutes":    60,
				"open_time":                "07:00 AM",
				"close_time":               "09:00 PM",
				"max_days_advance_booking": 30,
			},
			FAQs: []map[string]interface{}{
				{
					"question": "What's the court condition?",
					"answer":   "Courts are maintained daily by professional groundstaff",
				},
			},
			TermsAndConditions: []string{
				"Proper tennis attire required",
				"Membership available for monthly passes",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	batch := client.Batch()
	for _, venue := range playVenues {
		ref := client.Collection("play_venues").Doc(venue.ID)
		batch.Set(ref, venue)
	}

	_, err := batch.Commit(ctx)
	return err
}

func seedDiningVenues(ctx context.Context, client *firestore.Client, organizerID string) error {
	diningVenues := []DiningVenueData{
		{
			ID:               uuid.New().String(),
			OrganizerID:      organizerID,
			Name:             "The Golden Plate",
			Slug:             "the-golden-plate",
			Status:           "active",
			Description:      "Award-winning restaurant serving contemporary Indian cuisine with international influences. Chef-driven menu with seasonal ingredients.",
			ShortDescription: "Fine dining Indian restaurant with award-winning cuisine",
			Rating:           4.8,
			IsOpen:           true,
			OpeningTime:      "12:00 PM",
			ClosingTime:      "11:00 PM",
			Contact: map[string]interface{}{
				"phone": "+91-80-4567-8900",
				"email": "info@thegoldenplate.com",
			},
			Location: map[string]interface{}{
				"venue_name": "The Golden Plate",
				"address":    "123, Grand Central Mall, MG Road",
				"city":       "Bangalore",
				"state":      "Karnataka",
				"latitude":   12.9698,
				"longitude":  77.6109,
				"map_url":    "https://maps.google.com/?q=12.9698,77.6109",
			},
			Images: map[string]interface{}{
				"hero":    "https://images.unsplash.com/photo-1504674900152-b8d46e50f3e5?w=1200",
				"gallery": []string{"https://images.unsplash.com/photo-1546069901-ba9599a7e63c?w=500", "https://images.unsplash.com/photo-1555939594-58d7cb561621?w=500"},
			},
			MenuImages: []string{
				"https://images.unsplash.com/photo-1546069901-ba9599a7e63c?w=600",
			},
			PopularDishes: []string{
				"Butter Chicken",
				"Tandoori Salmon",
				"Paneer Tikka",
				"Biryani",
				"Dal Makhani",
			},
			Offers: []map[string]interface{}{
				{
					"title":       "Happy Hour Special",
					"code":        "HAPPYHOUR20",
					"description": "20% off on beverages 5-7 PM",
				},
				{
					"title":       "Birthday Celebration",
					"code":        "BIRTHDAY25",
					"description": "25% off for birthday guest during their birthday month",
				},
			},
			Facilities: []string{
				"WiFi",
				"Parking",
				"Private Dining Rooms",
				"Bar",
				"Live Music",
				"AC",
				"Family Friendly",
			},
			SeatingTypes: []map[string]interface{}{
				{
					"type":               "Indoor",
					"total_tables":       20,
					"available_tables":   15,
					"capacity_per_table": 4,
				},
				{
					"type":               "Outdoor Garden",
					"total_tables":       8,
					"available_tables":   6,
					"capacity_per_table": 4,
				},
			},
			BookingSettings: map[string]interface{}{
				"advance_booking_days":            30,
				"time_slots":                      []string{"12:00 PM", "1:00 PM", "2:00 PM", "6:00 PM", "7:00 PM", "8:00 PM", "9:00 PM"},
				"average_dining_duration_minutes": 90,
			},
			FAQs: []map[string]interface{}{
				{
					"question": "What is the dress code?",
					"answer":   "Formal attire recommended. No flip-flops allowed.",
				},
				{
					"question": "Do you accommodate dietary restrictions?",
					"answer":   "Yes! We cater to vegan, gluten-free, and all other dietary needs.",
				},
			},
			TermsAndConditions: []string{
				"Reservation must be confirmed 24 hours before",
				"Cancellation charges apply if cancelled within 12 hours",
				"No outside food or beverages allowed",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:               uuid.New().String(),
			OrganizerID:      organizerID,
			Name:             "Street Kitchen Dhaba",
			Slug:             "street-kitchen-dhaba",
			Status:           "active",
			Description:      "Authentic street food and traditional Indian dhaba experience. Specializing in North Indian cuisine with comfort food classics.",
			ShortDescription: "Authentic street food and traditional dhaba cuisine",
			Rating:           4.5,
			IsOpen:           true,
			OpeningTime:      "11:00 AM",
			ClosingTime:      "12:00 AM",
			Contact: map[string]interface{}{
				"phone": "+91-80-1234-5678",
				"email": "hello@streetkitchen.com",
			},
			Location: map[string]interface{}{
				"venue_name": "Street Kitchen Dhaba",
				"address":    "567, Indiranagar Park, 100 Feet Road",
				"city":       "Bangalore",
				"state":      "Karnataka",
				"latitude":   12.9716,
				"longitude":  77.6412,
				"map_url":    "https://maps.google.com/?q=12.9716,77.6412",
			},
			Images: map[string]interface{}{
				"hero":    "https://images.unsplash.com/photo-1567574883219-894e6becbbf7?w=1200",
				"gallery": []string{"https://images.unsplash.com/photo-1585937421456-2c89f1b61f19?w=500"},
			},
			MenuImages: []string{
				"https://images.unsplash.com/photo-1585937421456-2c89f1b61f19?w=600",
			},
			PopularDishes: []string{
				"Chole Bhature",
				"Dosa",
				"Samosas",
				"Chaat",
				"Naan & Curry",
			},
			Offers: []map[string]interface{}{
				{
					"title":       "Combo Deals",
					"code":        "COMBO30",
					"description": "30% off on combo meals",
				},
			},
			Facilities: []string{
				"WiFi",
				"Casual Dining",
				"Budget Friendly",
				"Quick Service",
			},
			SeatingTypes: []map[string]interface{}{
				{
					"type":               "Casual Dining",
					"total_tables":       30,
					"available_tables":   25,
					"capacity_per_table": 2,
				},
			},
			BookingSettings: map[string]interface{}{
				"advance_booking_days":            15,
				"time_slots":                      []string{"11:00 AM", "12:00 PM", "1:00 PM", "6:00 PM", "7:00 PM", "8:00 PM"},
				"average_dining_duration_minutes": 45,
			},
			FAQs: []map[string]interface{}{
				{
					"question": "Is parking available?",
					"answer":   "Street parking available. No dedicated parking lot.",
				},
			},
			TermsAndConditions: []string{
				"Casual dress code",
				"Fast service, great for quick meals",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	batch := client.Batch()
	for _, venue := range diningVenues {
		ref := client.Collection("dining_venues").Doc(venue.ID)
		batch.Set(ref, venue)
	}

	_, err := batch.Commit(ctx)
	return err
}

func seedEvents(ctx context.Context, client *firestore.Client, organizerID string) error {
	now := time.Now()

	events := []EventData{
		{
			ID:               uuid.New().String(),
			OrganizerID:      organizerID,
			Title:            "Live Jazz Night with The Blue Notes",
			Slug:             "live-jazz-night-blue-notes",
			Category:         "Music",
			Language:         "English",
			DurationMinutes:  180,
			AgeLimit:         "18+",
			Description:      "Experience an unforgettable evening of live jazz music with The Blue Notes. A legendary band bringing classic and contemporary jazz together. Perfect evening for jazz enthusiasts.",
			ShortDescription: "Live jazz performance with The Blue Notes",
			StartDatetime:    now.AddDate(0, 0, 15).Add(7 * time.Hour),
			EndDatetime:      now.AddDate(0, 0, 15).Add(10 * time.Hour),
			PriceStart:       500.0,
			Status:           "active",
			Venue: map[string]interface{}{
				"name":      "The Jazz Lounge",
				"address":   "789, Entertainment District, Whitefield",
				"city":      "Bangalore",
				"state":     "Karnataka",
				"latitude":  12.9698,
				"longitude": 77.7499,
				"map_url":   "https://maps.google.com/?q=12.9698,77.7499",
			},
			Images: map[string]interface{}{
				"hero":    "https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=1200",
				"poster":  "https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=600",
				"gallery": []string{"https://images.unsplash.com/photo-1514525253161-7a46d19cd819?w=500"},
			},
			Artists: []map[string]interface{}{
				{
					"name":             "Marcus Johnson",
					"role":             "Lead Saxophone",
					"image_url":        "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=300",
					"description":      "Grammy-nominated saxophonist",
					"genre":            "Jazz",
					"is_verified":      true,
					"rating":           4.9,
					"review_count":     450,
					"events_hosted":    85,
					"location":         "New York, USA",
					"contact_email":    "marcus@bluenotes.com",
					"contact_phone":    "+1-212-555-0123",
					"experience_years": 25,
					"specialties":      []string{"Bebop", "Modal Jazz", "Soul Jazz"},
					"follower_count":   50000,
					"social_links":     []string{"https://twitter.com/marcus_jazz"},
				},
				{
					"name":             "Sarah Williams",
					"role":             "Piano",
					"image_url":        "https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=300",
					"description":      "Concert pianist and composer",
					"genre":            "Jazz",
					"is_verified":      true,
					"rating":           4.8,
					"review_count":     380,
					"events_hosted":    72,
					"location":         "Boston, USA",
					"contact_email":    "sarah@bluenotes.com",
					"contact_phone":    "+1-617-555-0123",
					"experience_years": 20,
					"specialties":      []string{"Contemporary Jazz", "Classical", "Fusion"},
					"follower_count":   40000,
					"social_links":     []string{"https://twitter.com/sarah_jazz"},
				},
			},
			Tickets: []map[string]interface{}{
				{
					"ticket_type":        "General Admission",
					"seat_type":          "Regular",
					"price":              500.0,
					"total_quantity":     200,
					"available_quantity": 150,
				},
				{
					"ticket_type":        "VIP Seating",
					"seat_type":          "Front Row",
					"price":              1000.0,
					"total_quantity":     50,
					"available_quantity": 35,
				},
				{
					"ticket_type":        "Premium with Dinner",
					"seat_type":          "VIP Table",
					"price":              1500.0,
					"total_quantity":     30,
					"available_quantity": 20,
				},
			},
			FAQs: []map[string]interface{}{
				{
					"question": "Is there a dress code?",
					"answer":   "Smart casual attire recommended",
				},
				{
					"question": "When do gates open?",
					"answer":   "Gates open at 6:00 PM, event starts at 7:00 PM",
				},
			},
			TermsAndConditions: []string{
				"Tickets are non-refundable",
				"Valid ID required for entry",
				"No outside food or beverages",
				"Photography permitted for personal use only",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:               uuid.New().String(),
			OrganizerID:      organizerID,
			Title:            "Comedy Night: Stand-up Comedy Extravaganza",
			Slug:             "comedy-night-stand-up",
			Category:         "Comedy",
			Language:         "English",
			DurationMinutes:  120,
			AgeLimit:         "16+",
			Description:      "Laugh out loud with India's top comedians. An evening filled with hilarious stand-up performances covering topics from daily life to current affairs.",
			ShortDescription: "India's top comedians performing live",
			StartDatetime:    now.AddDate(0, 0, 20).Add(8 * time.Hour),
			EndDatetime:      now.AddDate(0, 0, 20).Add(10 * time.Hour),
			PriceStart:       400.0,
			Status:           "active",
			Venue: map[string]interface{}{
				"name":      "Comedy Central Theater",
				"address":   "234, Entertainment Plaza, Bellandur",
				"city":      "Bangalore",
				"state":     "Karnataka",
				"latitude":  12.9387,
				"longitude": 77.6821,
				"map_url":   "https://maps.google.com/?q=12.9387,77.6821",
			},
			Images: map[string]interface{}{
				"hero":    "https://images.unsplash.com/photo-1516575334481-f410b4ae5f69?w=1200",
				"poster":  "https://images.unsplash.com/photo-1516575334481-f410b4ae5f69?w=600",
				"gallery": []string{},
			},
			Artists: []map[string]interface{}{
				{
					"name":             "Aditya Verma",
					"role":             "Stand-up Comedian",
					"image_url":        "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=300",
					"description":      "India's funniest comedian",
					"genre":            "Comedy",
					"is_verified":      true,
					"rating":           4.9,
					"review_count":     520,
					"events_hosted":    120,
					"location":         "Mumbai, India",
					"contact_email":    "aditya@comedycentral.com",
					"contact_phone":    "+91-9999-111111",
					"experience_years": 12,
					"specialties":      []string{"Observational Comedy", "Political Satire"},
					"follower_count":   200000,
					"social_links":     []string{"https://twitter.com/aditya_comedy"},
				},
			},
			Tickets: []map[string]interface{}{
				{
					"ticket_type":        "Standard",
					"seat_type":          "Regular Seating",
					"price":              400.0,
					"total_quantity":     300,
					"available_quantity": 250,
				},
				{
					"ticket_type":        "Premium",
					"seat_type":          "Front Block",
					"price":              800.0,
					"total_quantity":     100,
					"available_quantity": 75,
				},
			},
			FAQs: []map[string]interface{}{
				{
					"question": "Can I bring kids?",
					"answer":   "Kids below 16 not allowed due to content",
				},
			},
			TermsAndConditions: []string{
				"No recording allowed",
				"Arrive 15 minutes early",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:               uuid.New().String(),
			OrganizerID:      organizerID,
			Title:            "Art Exhibition: Contemporary Indian Artists",
			Slug:             "art-exhibition-contemporary",
			Category:         "Art & Culture",
			Language:         "English",
			DurationMinutes:  240,
			AgeLimit:         "All Ages",
			Description:      "Explore fascinating works by contemporary Indian artists. This exhibition showcases modern art, installations, and interactive exhibits celebrating Indian culture and global perspectives.",
			ShortDescription: "Contemporary Indian art exhibition",
			StartDatetime:    now.AddDate(0, 0, 25).Add(10 * time.Hour),
			EndDatetime:      now.AddDate(0, 1, 5).Add(8 * time.Hour),
			PriceStart:       200.0,
			Status:           "active",
			Venue: map[string]interface{}{
				"name":      "Bangalore Art Museum",
				"address":   "890, Museum Road, Cubbon Park",
				"city":      "Bangalore",
				"state":     "Karnataka",
				"latitude":  12.9716,
				"longitude": 77.5946,
				"map_url":   "https://maps.google.com/?q=12.9716,77.5946",
			},
			Images: map[string]interface{}{
				"hero":    "https://images.unsplash.com/photo-1578301978162-7aae4d755744?w=1200",
				"poster":  "https://images.unsplash.com/photo-1578301978162-7aae4d755744?w=600",
				"gallery": []string{"https://images.unsplash.com/photo-1552820728-8ac41f1ce891?w=500"},
			},
			Artists: []map[string]interface{}{
				{
					"name":             "Priya Sharma",
					"role":             "Contemporary Artist & Curator",
					"image_url":        "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=300",
					"description":      "Award-winning contemporary artist",
					"genre":            "Contemporary Art",
					"is_verified":      true,
					"rating":           4.7,
					"review_count":     290,
					"events_hosted":    45,
					"location":         "New Delhi, India",
					"contact_email":    "priya@artgallery.com",
					"contact_phone":    "+91-9876-543210",
					"experience_years": 15,
					"specialties":      []string{"Mixed Media", "Installation Art", "Social Commentary"},
					"follower_count":   30000,
					"social_links":     []string{"https://instagram.com/priya_art"},
				},
			},
			Tickets: []map[string]interface{}{
				{
					"ticket_type":        "General Entry",
					"seat_type":          "Walking Tour",
					"price":              200.0,
					"total_quantity":     500,
					"available_quantity": 450,
				},
				{
					"ticket_type":        "Guided Tour with Artist",
					"seat_type":          "Group Tour",
					"price":              500.0,
					"total_quantity":     80,
					"available_quantity": 60,
				},
			},
			FAQs: []map[string]interface{}{
				{
					"question": "Can I take photographs?",
					"answer":   "Photography allowed except for special installations marked 'No Photography'",
				},
				{
					"question": "Is it wheel-chair accessible?",
					"answer":   "Yes, the museum is fully wheelchair accessible",
				},
			},
			TermsAndConditions: []string{
				"Respect artwork - do not touch",
				"Keep the venue clean",
				"Guided tours must start at scheduled times",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	batch := client.Batch()
	for _, event := range events {
		ref := client.Collection("events").Doc(event.ID)
		batch.Set(ref, event)
	}

	_, err := batch.Commit(ctx)
	return err
}
