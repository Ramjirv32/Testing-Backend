package repository

import (
	"context"
	"fmt"
	"sort"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"backend/config"
	"backend/models"
)

type EventRepository struct{}

func (r *EventRepository) c() *firestore.CollectionRef {
	if config.FirestoreClient == nil {
		panic("FirestoreClient is not initialized. Check InitFirebase logs.")
	}
	return config.FirestoreClient.Collection("events")
}

func NewEventRepository() *EventRepository {
	return &EventRepository{}
}

func (r *EventRepository) Create(ctx context.Context, e *models.Event) error {
	_, err := r.c().Doc(e.ID).Set(ctx, e)
	return err
}

func (r *EventRepository) FindByID(ctx context.Context, id string) (*models.Event, error) {
	if config.FirestoreClient == nil {
		return nil, fmt.Errorf("Firestore client is not initialized")
	}

	doc, err := r.c().Doc(id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get event by ID: %w", err)
	}

	var e models.Event
	if err := doc.DataTo(&e); err != nil {
		return nil, fmt.Errorf("failed to parse event data: %w", err)
	}
	return &e, nil
}

func (r *EventRepository) FindPaginatedByOrganizerID(ctx context.Context, organizerID string, limit int, lastID string) ([]*models.Event, string, error) {
	// Fetch all events for this organizer (no OrderBy to avoid composite index requirement)
	q := r.c().Where("organizer_id", "==", organizerID)

	iter := q.Documents(ctx)
	allEvents := []*models.Event{}

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, "", err
		}
		var e models.Event
		if err := doc.DataTo(&e); err != nil {
			continue
		}
		allEvents = append(allEvents, &e)
	}

	// Sort by created_at descending in memory
	sort.Slice(allEvents, func(i, j int) bool {
		return allEvents[i].CreatedAt.After(allEvents[j].CreatedAt)
	})

	// Apply cursor-based pagination
	startIdx := 0
	if lastID != "" {
		for i, e := range allEvents {
			if e.ID == lastID {
				startIdx = i + 1
				break
			}
		}
	}

	end := startIdx + limit
	if end > len(allEvents) {
		end = len(allEvents)
	}

	page := allEvents[startIdx:end]
	nextCursor := ""
	if len(page) > 0 && end < len(allEvents) {
		nextCursor = page[len(page)-1].ID
	}

	return page, nextCursor, nil
}

func (r *EventRepository) FindByOrganizerID(ctx context.Context, organizerID string) ([]*models.Event, error) {
	events, _, err := r.FindPaginatedByOrganizerID(ctx, organizerID, 100, "")
	return events, err
}

func (r *EventRepository) GetPaginated(ctx context.Context, limit int, lastID string, category string, city string, searchQuery string, status string) ([]*models.Event, string, error) {
	q := r.c().Limit(limit)

	if status != "" {
		q = q.Where("status", "==", status)
	}

	if category != "" {
		q = q.Where("category", "==", category)
	}

	if city != "" {
		q = q.Where("venue.city", "==", city)
	}

	if lastID != "" {
		lastDoc, err := r.c().Doc(lastID).Get(ctx)
		if err == nil {
			q = q.StartAfter(lastDoc)
		}
	}

	iter := q.Documents(ctx)
	events := []*models.Event{}
	var lastDocID string

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, "", err
		}

		var e models.Event
		if err := doc.DataTo(&e); err != nil {
			continue
		}
		events = append(events, &e)
		lastDocID = doc.Ref.ID
	}

	return events, lastDocID, nil
}

func (r *EventRepository) GetAll(ctx context.Context, status string) ([]*models.Event, error) {
	events, _, err := r.GetPaginated(ctx, 100, "", "", "", "", status)
	return events, err
}

func (r *EventRepository) Update(ctx context.Context, e *models.Event) error {
	e.UpdatedAt = time.Now()
	_, err := r.c().Doc(e.ID).Set(ctx, e)
	return err
}

func (r *EventRepository) Delete(ctx context.Context, id string) error {
	_, err := r.c().Doc(id).Delete(ctx)
	return err
}
