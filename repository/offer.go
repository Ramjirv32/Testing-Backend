package repository

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"backend/config"
	"backend/models"
)

type OfferRepository struct{}

func (r *OfferRepository) c() *firestore.CollectionRef {
	if config.FirestoreClient == nil {
		panic("FirestoreClient is not initialized. Check InitFirebase logs.")
	}
	return config.FirestoreClient.Collection("offers")
}

func NewOfferRepository() *OfferRepository {
	return &OfferRepository{}
}

func (r *OfferRepository) Create(ctx context.Context, o *models.Offer) error {
	if o.ID == "" {
		o.ID = r.c().NewDoc().ID
	}
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	_, err := r.c().Doc(o.ID).Set(ctx, o)
	return err
}

func (r *OfferRepository) GetByUserID(ctx context.Context, userID string) ([]*models.Offer, error) {
	iter := r.c().Where("user_id", "==", userID).Documents(ctx)
	var offers []*models.Offer

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var o models.Offer
		if err := doc.DataTo(&o); err != nil {
			continue
		}
		offers = append(offers, &o)
	}

	return offers, nil
}

func (r *OfferRepository) GetAll(ctx context.Context) ([]*models.Offer, error) {
	iter := r.c().Documents(ctx)
	var offers []*models.Offer

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var o models.Offer
		if err := doc.DataTo(&o); err != nil {
			continue
		}
		offers = append(offers, &o)
	}

	return offers, nil
}
