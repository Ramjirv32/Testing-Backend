package repository

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"backend/config"
	"backend/models"
)

type PlayRepository struct{}

func (r *PlayRepository) c() *firestore.CollectionRef {
	return config.FirestoreClient.Collection("play_venues")
}

func NewPlayRepository() *PlayRepository {
	return &PlayRepository{}
}

func (r *PlayRepository) Create(ctx context.Context, p *models.PlayVenue) error {
	if p.ID == "" {
		p.ID = config.FirestoreClient.Collection("play_venues").NewDoc().ID
	}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	_, err := r.c().Doc(p.ID).Set(ctx, p)
	return err
}

func (r *PlayRepository) GetPaginated(ctx context.Context, limit int, lastID string, category string, city string, searchQuery string) ([]*models.PlayVenue, string, error) {
	q := r.c().Limit(limit)

	if category != "" {
		q = q.Where("category", "==", category)
	}

	if city != "" {
		q = q.Where("location.city", "==", city)
	}

	// Removed OrderBy to avoid index requirement for simple city/category filters
	// q = q.OrderBy("created_at", firestore.Desc)

	if lastID != "" {
		lastDoc, err := r.c().Doc(lastID).Get(ctx)
		if err == nil {
			q = q.StartAfter(lastDoc)
		}
	}

	iter := q.Documents(ctx)
	var venues []*models.PlayVenue
	var lastDocID string

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, "", err
		}

		var v models.PlayVenue
		if err := doc.DataTo(&v); err != nil {
			continue
		}
		venues = append(venues, &v)
		lastDocID = doc.Ref.ID
	}

	return venues, lastDocID, nil
}

func (r *PlayRepository) GetAll(ctx context.Context) ([]*models.PlayVenue, error) {
	venues, _, err := r.GetPaginated(ctx, 100, "", "", "", "")
	return venues, err
}

func (r *PlayRepository) GetBySlug(ctx context.Context, slug string) (*models.PlayVenue, error) {
	iter := r.c().Where("slug", "==", slug).Limit(1).Documents(ctx)
	doc, err := iter.Next()
	if err != nil {
		return nil, err
	}

	var v models.PlayVenue
	if err := doc.DataTo(&v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *PlayRepository) FindByID(ctx context.Context, id string) (*models.PlayVenue, error) {
	doc, err := r.c().Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	var v models.PlayVenue
	if err := doc.DataTo(&v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *PlayRepository) Update(ctx context.Context, v *models.PlayVenue) error {
	v.UpdatedAt = time.Now()
	_, err := r.c().Doc(v.ID).Set(ctx, v)
	return err
}

func (r *PlayRepository) Delete(ctx context.Context, id string) error {
	_, err := r.c().Doc(id).Delete(ctx)
	return err
}

func (r *PlayRepository) FindPaginatedByOrganizerID(ctx context.Context, organizerID string, limit int, lastID string) ([]*models.PlayVenue, string, error) {
	q := r.c().Where("organizer_id", "==", organizerID).OrderBy("created_at", firestore.Desc).Limit(limit)

	if lastID != "" {
		lastDoc, err := r.c().Doc(lastID).Get(ctx)
		if err == nil {
			q = q.StartAfter(lastDoc)
		}
	}

	iter := q.Documents(ctx)
	var venues []*models.PlayVenue
	var lastDocID string

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, "", err
		}

		var v models.PlayVenue
		if err := doc.DataTo(&v); err != nil {
			continue
		}
		venues = append(venues, &v)
		lastDocID = doc.Ref.ID
	}

	return venues, lastDocID, nil
}

func (r *PlayRepository) FindByOrganizerID(ctx context.Context, organizerID string) ([]*models.PlayVenue, error) {
	venues, _, err := r.FindPaginatedByOrganizerID(ctx, organizerID, 100, "")
	return venues, err
}
