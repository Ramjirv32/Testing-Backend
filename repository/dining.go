package repository

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"backend/config"
	"backend/models"
)

type DiningRepository struct{}

func (r *DiningRepository) c() *firestore.CollectionRef {
	return config.FirestoreClient.Collection("dining_venues")
}

func NewDiningRepository() *DiningRepository {
	return &DiningRepository{}
}

func (r *DiningRepository) Create(ctx context.Context, d *models.DiningVenue) error {
	if d.ID == "" {
		d.ID = config.FirestoreClient.Collection("dining_venues").NewDoc().ID
	}
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
	_, err := r.c().Doc(d.ID).Set(ctx, d)
	return err
}

func (r *DiningRepository) GetPaginated(ctx context.Context, limit int, lastID string, category string, city string, searchQuery string) ([]*models.DiningVenue, string, error) {
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
	var venues []*models.DiningVenue
	var lastDocID string

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, "", err
		}

		var d models.DiningVenue
		if err := doc.DataTo(&d); err != nil {
			continue
		}
		venues = append(venues, &d)
		lastDocID = doc.Ref.ID
	}

	return venues, lastDocID, nil
}

func (r *DiningRepository) GetAll(ctx context.Context) ([]*models.DiningVenue, error) {
	venues, _, err := r.GetPaginated(ctx, 100, "", "", "", "")
	return venues, err
}

func (r *DiningRepository) GetBySlug(ctx context.Context, slug string) (*models.DiningVenue, error) {
	iter := r.c().Where("slug", "==", slug).Limit(1).Documents(ctx)
	doc, err := iter.Next()
	if err != nil {
		return nil, err
	}

	var d models.DiningVenue
	if err := doc.DataTo(&d); err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DiningRepository) FindByID(ctx context.Context, id string) (*models.DiningVenue, error) {
	doc, err := r.c().Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	var d models.DiningVenue
	if err := doc.DataTo(&d); err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DiningRepository) Update(ctx context.Context, d *models.DiningVenue) error {
	d.UpdatedAt = time.Now()
	_, err := r.c().Doc(d.ID).Set(ctx, d)
	return err
}

func (r *DiningRepository) Delete(ctx context.Context, id string) error {
	_, err := r.c().Doc(id).Delete(ctx)
	return err
}

func (r *DiningRepository) FindPaginatedByOrganizerID(ctx context.Context, organizerID string, limit int, lastID string) ([]*models.DiningVenue, string, error) {
	q := r.c().Where("organizer_id", "==", organizerID).OrderBy("created_at", firestore.Desc).Limit(limit)

	if lastID != "" {
		lastDoc, err := r.c().Doc(lastID).Get(ctx)
		if err == nil {
			q = q.StartAfter(lastDoc)
		}
	}

	iter := q.Documents(ctx)
	var venues []*models.DiningVenue
	var lastDocID string

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, "", err
		}

		var d models.DiningVenue
		if err := doc.DataTo(&d); err != nil {
			continue
		}
		venues = append(venues, &d)
		lastDocID = doc.Ref.ID
	}

	return venues, lastDocID, nil
}

func (r *DiningRepository) FindByOrganizerID(ctx context.Context, organizerID string) ([]*models.DiningVenue, error) {
	venues, _, err := r.FindPaginatedByOrganizerID(ctx, organizerID, 100, "")
	return venues, err
}
