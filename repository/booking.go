package repository

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"backend/config"
	"backend/models"
)

type PlayBookingRepository struct{}

func (r *PlayBookingRepository) c() *firestore.CollectionRef {
	if config.FirestoreClient == nil {
		panic("FirestoreClient is not initialized. Check InitFirebase logs.")
	}
	return config.FirestoreClient.Collection("play_bookings")
}

func NewPlayBookingRepository() *PlayBookingRepository {
	return &PlayBookingRepository{}
}

func (r *PlayBookingRepository) Create(ctx context.Context, booking *models.PlayBooking) error {
	_, err := r.c().Doc(booking.ID).Set(ctx, booking)
	return err
}

func (r *PlayBookingRepository) FindByID(ctx context.Context, id string) (*models.PlayBooking, error) {
	if config.FirestoreClient == nil {
		return nil, fmt.Errorf("Firestore client is not initialized")
	}

	doc, err := r.c().Doc(id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get play booking by ID: %w", err)
	}

	var booking models.PlayBooking
	if err := doc.DataTo(&booking); err != nil {
		return nil, fmt.Errorf("failed to parse play booking data: %w", err)
	}
	return &booking, nil
}

func (r *PlayBookingRepository) FindByUserID(ctx context.Context, userID string) ([]*models.PlayBooking, error) {
	iter := r.c().Where("user_id", "==", userID).Documents(ctx)
	var bookings []*models.PlayBooking

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate play bookings: %w", err)
		}

		var booking models.PlayBooking
		if err := doc.DataTo(&booking); err != nil {
			continue
		}
		bookings = append(bookings, &booking)
	}

	return bookings, nil
}

func (r *PlayBookingRepository) Update(ctx context.Context, booking *models.PlayBooking) error {
	booking.UpdatedAt = time.Now()
	_, err := r.c().Doc(booking.ID).Set(ctx, booking)
	return err
}

func (r *PlayBookingRepository) Delete(ctx context.Context, id string) error {
	_, err := r.c().Doc(id).Delete(ctx)
	return err
}

func (r *PlayBookingRepository) CheckPlayAvailability(ctx context.Context, venueID, date, timeSlot string) (bool, error) {
	iter := r.c().Where("venue_id", "==", venueID).Where("date", "==", date).Where("time_slot", "==", timeSlot).Where("status", "==", models.BookingConfirmed).Documents(ctx)
	snaps, err := iter.GetAll()
	if err != nil {
		return false, err
	}
	// For now, assume 1 court per slot. If anyone booked, it's unavailable.
	return len(snaps) == 0, nil
}

func (r *PlayBookingRepository) DeleteAll(ctx context.Context) error {
	iter := r.c().Documents(ctx)
	batch := config.FirestoreClient.Batch()
	count := 0

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		batch.Delete(doc.Ref)
		count++

		if count >= 500 {
			if _, err := batch.Commit(ctx); err != nil {
				return err
			}
			batch = config.FirestoreClient.Batch()
			count = 0
		}
	}

	if count > 0 {
		_, err := batch.Commit(ctx)
		return err
	}

	return nil
}

type DiningBookingRepository struct{}

func (r *DiningBookingRepository) c() *firestore.CollectionRef {
	if config.FirestoreClient == nil {
		panic("FirestoreClient is not initialized. Check InitFirebase logs.")
	}
	return config.FirestoreClient.Collection("dining_bookings")
}

func NewDiningBookingRepository() *DiningBookingRepository {
	return &DiningBookingRepository{}
}

func (r *DiningBookingRepository) Create(ctx context.Context, booking *models.DiningBooking) error {
	_, err := r.c().Doc(booking.ID).Set(ctx, booking)
	return err
}

func (r *DiningBookingRepository) FindByID(ctx context.Context, id string) (*models.DiningBooking, error) {
	if config.FirestoreClient == nil {
		return nil, fmt.Errorf("Firestore client is not initialized")
	}

	doc, err := r.c().Doc(id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get dining booking by ID: %w", err)
	}

	var booking models.DiningBooking
	if err := doc.DataTo(&booking); err != nil {
		return nil, fmt.Errorf("failed to parse dining booking data: %w", err)
	}
	return &booking, nil
}

func (r *DiningBookingRepository) FindByUserID(ctx context.Context, userID string) ([]*models.DiningBooking, error) {
	iter := r.c().Where("user_id", "==", userID).Documents(ctx)
	var bookings []*models.DiningBooking

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate dining bookings: %w", err)
		}

		var booking models.DiningBooking
		if err := doc.DataTo(&booking); err != nil {
			continue
		}
		bookings = append(bookings, &booking)
	}

	return bookings, nil
}

func (r *DiningBookingRepository) Update(ctx context.Context, booking *models.DiningBooking) error {
	booking.UpdatedAt = time.Now()
	_, err := r.c().Doc(booking.ID).Set(ctx, booking)
	return err
}

func (r *DiningBookingRepository) Delete(ctx context.Context, id string) error {
	_, err := r.c().Doc(id).Delete(ctx)
	return err
}

func (r *DiningBookingRepository) CheckDiningAvailability(ctx context.Context, restaurantID, date, timeSlot string) (bool, error) {
	iter := r.c().Where("restaurant_id", "==", restaurantID).Where("date", "==", date).Where("time_slot", "==", timeSlot).Where("status", "==", models.BookingConfirmed).Documents(ctx)
	snaps, err := iter.GetAll()
	if err != nil {
		return false, err
	}
	// For now, assume a limit of 5 tables per slot.
	return len(snaps) < 5, nil
}

func (r *DiningBookingRepository) DeleteAll(ctx context.Context) error {
	iter := r.c().Documents(ctx)
	batch := config.FirestoreClient.Batch()
	count := 0

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		batch.Delete(doc.Ref)
		count++

		if count >= 500 {
			if _, err := batch.Commit(ctx); err != nil {
				return err
			}
			batch = config.FirestoreClient.Batch()
			count = 0
		}
	}

	if count > 0 {
		_, err := batch.Commit(ctx)
		return err
	}

	return nil
}

type EventBookingRepository struct{}

func (r *EventBookingRepository) c() *firestore.CollectionRef {
	if config.FirestoreClient == nil {
		panic("FirestoreClient is not initialized. Check InitFirebase logs.")
	}
	return config.FirestoreClient.Collection("event_bookings")
}

func NewEventBookingRepository() *EventBookingRepository {
	return &EventBookingRepository{}
}

func (r *EventBookingRepository) Create(ctx context.Context, booking *models.EventBooking) error {
	_, err := r.c().Doc(booking.ID).Set(ctx, booking)
	return err
}

func (r *EventBookingRepository) FindByID(ctx context.Context, id string) (*models.EventBooking, error) {
	if config.FirestoreClient == nil {
		return nil, fmt.Errorf("Firestore client is not initialized")
	}

	doc, err := r.c().Doc(id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get event booking by ID: %w", err)
	}

	var booking models.EventBooking
	if err := doc.DataTo(&booking); err != nil {
		return nil, fmt.Errorf("failed to parse event booking data: %w", err)
	}
	return &booking, nil
}

func (r *EventBookingRepository) FindByUserID(ctx context.Context, userID string) ([]*models.EventBooking, error) {
	iter := r.c().Where("user_id", "==", userID).Documents(ctx)
	var bookings []*models.EventBooking

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate event bookings: %w", err)
		}

		var booking models.EventBooking
		if err := doc.DataTo(&booking); err != nil {
			continue
		}
		bookings = append(bookings, &booking)
	}

	return bookings, nil
}

func (r *EventBookingRepository) Update(ctx context.Context, booking *models.EventBooking) error {
	booking.UpdatedAt = time.Now()
	_, err := r.c().Doc(booking.ID).Set(ctx, booking)
	return err
}

func (r *EventBookingRepository) Delete(ctx context.Context, id string) error {
	_, err := r.c().Doc(id).Delete(ctx)
	return err
}
