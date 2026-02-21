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

func (r *PartnerRepository) FindByPAN(ctx context.Context, pan string) (*models.PartnerProfile, error) {
	iter := r.c().Where("organization_details.pan", "==", pan).Limit(1).Documents(ctx)
	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query partner by PAN: %w", err)
	}

	var ep models.PartnerProfile
	if err := doc.DataTo(&ep); err != nil {
		return nil, fmt.Errorf("failed to parse partner data: %w", err)
	}
	return &ep, nil
}

func (r *PartnerRepository) FindByBankAccountNumber(ctx context.Context, accountNumber string) (*models.PartnerProfile, error) {
	iter := r.c().Where("bank_details.account_number", "==", accountNumber).Limit(1).Documents(ctx)
	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query partner by bank account number: %w", err)
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

func (r *PartnerRepository) FindAllByUserID(ctx context.Context, userID string) ([]*models.PartnerProfile, error) {
	iter := r.c().Where("user_id", "==", userID).Documents(ctx)
	var profiles []*models.PartnerProfile
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to query partners by user ID: %w", err)
		}
		var ep models.PartnerProfile
		if err := doc.DataTo(&ep); err != nil {
			continue
		}
		profiles = append(profiles, &ep)
	}
	return profiles, nil
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

func (r *PartnerRepository) GetPaginated(ctx context.Context, limit int, lastID string, category string, status string) ([]*models.PartnerProfile, string, error) {
	var q firestore.Query

	q = r.c().Query

	if category != "" && category != "all" {
		q = q.Where("organization_details.category", "==", category)
	}

	if status != "" && status != "all" {
		q = q.Where("status", "==", status)
	}

	if category == "" && status == "" {
		q = q.OrderBy("created_at", firestore.Desc)
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

	posters, _, err := r.GetPaginated(ctx, 100, "", "", "")
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
