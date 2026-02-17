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

type PartnerRepository struct{}

func (r *PartnerRepository) c() *firestore.CollectionRef {
	if config.FirestoreClient == nil {
		panic("FirestoreClient is not initialized. Check InitFirebase logs.")
	}
	return config.FirestoreClient.Collection("partners")
}

func NewPartnerRepository() *PartnerRepository {
	return &PartnerRepository{}
}

func (r *PartnerRepository) Create(ctx context.Context, ep *models.PartnerProfile) error {
	_, err := r.c().Doc(ep.ID).Set(ctx, ep)
	return err
}

func (r *PartnerRepository) FindByUserIDAndCategory(ctx context.Context, userID string, category string) (*models.PartnerProfile, error) {
	q := r.c().Where("user_id", "==", userID)
	if category != "" {
		q = q.Where("organization_details.category", "==", category)
	}
	iter := q.Limit(1).Documents(ctx)
	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query partner by user ID: %w", err)
	}

	var ep models.PartnerProfile
	if err := doc.DataTo(&ep); err != nil {
		return nil, fmt.Errorf("failed to parse partner data: %w", err)
	}
	return &ep, nil
}

func (r *PartnerRepository) FindByUserID(ctx context.Context, userID string) (*models.PartnerProfile, error) {
	return r.FindByUserIDAndCategory(ctx, userID, "")
}

func (r *PartnerRepository) FindByID(ctx context.Context, id string) (*models.PartnerProfile, error) {
	if config.FirestoreClient == nil {
		return nil, fmt.Errorf("Firestore client is not initialized")
	}

	doc, err := r.c().Doc(id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get partner by ID: %w", err)
	}

	var ep models.PartnerProfile
	if err := doc.DataTo(&ep); err != nil {
		return nil, fmt.Errorf("failed to parse partner data: %w", err)
	}
	return &ep, nil
}

func (r *PartnerRepository) GetPaginated(ctx context.Context, limit int, lastID string, category string) ([]*models.PartnerProfile, string, error) {
	var q firestore.Query

	if category != "" && category != "all" {
		// When filtering by category, we skip OrderBy to avoid requiring a composite index
		// Note: Using capitalized 'Category' to match existing data stored without firestore tags
		q = r.c().Where("organization_details.category", "==", category)
	} else {
		q = r.c().OrderBy("created_at", firestore.Desc)
	}

	q = q.Limit(limit)

	if lastID != "" {
		lastDoc, err := r.c().Doc(lastID).Get(ctx)
		if err == nil {
			q = q.StartAfter(lastDoc)
		}
	}

	iter := q.Documents(ctx)
	var posters []*models.PartnerProfile
	var lastDocID string

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, "", fmt.Errorf("failed to iterate partners: %w", err)
		}

		var ep models.PartnerProfile
		if err := doc.DataTo(&ep); err != nil {
			continue
		}
		posters = append(posters, &ep)
		lastDocID = doc.Ref.ID
	}

	return posters, lastDocID, nil
}

func (r *PartnerRepository) GetAll(ctx context.Context) ([]*models.PartnerProfile, error) {
	// For enterprise scale, we limit GetAll to a reasonable default or use GetPaginated
	posters, _, err := r.GetPaginated(ctx, 100, "", "")
	return posters, err
}

func (r *PartnerRepository) Update(ctx context.Context, profile *models.PartnerProfile) error {
	profile.UpdatedAt = time.Now()
	_, err := r.c().Doc(profile.ID).Set(ctx, profile)
	return err
}

func (r *PartnerRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	_, err := r.c().Doc(id).Update(ctx, []firestore.Update{
		{Path: "status", Value: status},
		{Path: "updated_at", Value: time.Now()},
	})
	return err
}
