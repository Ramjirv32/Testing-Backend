package repository

import (
	"context"
	"fmt"
	"time"

	"backend/config"
	"backend/models"
)

// Extended offer repository methods

func (r *OfferRepository) GetByID(ctx context.Context, id string) (*models.Offer, error) {
	if config.FirestoreClient == nil {
		return nil, fmt.Errorf("Firestore client is not initialized")
	}

	doc, err := r.c().Doc(id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get offer: %w", err)
	}

	var o models.Offer
	if err := doc.DataTo(&o); err != nil {
		return nil, fmt.Errorf("failed to parse offer data: %w", err)
	}
	return &o, nil
}

func (r *OfferRepository) Update(ctx context.Context, o *models.Offer) error {
	o.UpdatedAt = time.Now()
	_, err := r.c().Doc(o.ID).Set(ctx, o)
	return err
}

func (r *OfferRepository) Delete(ctx context.Context, id string) error {
	_, err := r.c().Doc(id).Delete(ctx)
	return err
}
