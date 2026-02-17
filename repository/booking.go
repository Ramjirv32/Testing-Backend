package repository

import (
	"context"
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
	doc, err := r.c().Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	var booking models.PlayBooking
	if err := doc.DataTo(&booking); err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *PlayBookingRepository) FindByUserID(ctx context.Context, userID string) ([]*models.PlayBooking, error) {
	iter := r.c().Where("user_id", "==", userID).OrderBy("created_at", firestore.Desc).Documents(ctx)
	var bookings []*models.PlayBooking

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
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
	doc, err := r.c().Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	var booking models.DiningBooking
	if err := doc.DataTo(&booking); err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *DiningBookingRepository) FindByUserID(ctx context.Context, userID string) ([]*models.DiningBooking, error) {
	iter := r.c().Where("user_id", "==", userID).OrderBy("created_at", firestore.Desc).Documents(ctx)
	var bookings []*models.DiningBooking

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
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
