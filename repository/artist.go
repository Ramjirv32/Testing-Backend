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

type ArtistRepository struct{}

func (r *ArtistRepository) c() *firestore.CollectionRef {
	if config.FirestoreClient == nil {
		panic("FirestoreClient is not initialized. Check InitFirebase logs.")
	}
	return config.FirestoreClient.Collection("artists")
}

func NewArtistRepository() *ArtistRepository {
	return &ArtistRepository{}
}

func (r *ArtistRepository) Create(ctx context.Context, a *models.Artist) error {
	if a.ID == "" {
		a.ID = config.FirestoreClient.Collection("artists").NewDoc().ID
	}
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	_, err := r.c().Doc(a.ID).Set(ctx, a)
	return err
}

func (r *ArtistRepository) GetAll(ctx context.Context) ([]*models.Artist, error) {
	if config.FirestoreClient == nil {
		return nil, fmt.Errorf("Firestore client is not initialized")
	}

	iter := r.c().OrderBy("name", firestore.Asc).Documents(ctx)
	var artists []*models.Artist

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate artists: %w", err)
		}

		var a models.Artist
		if err := doc.DataTo(&a); err != nil {
			continue
		}
		artists = append(artists, &a)
	}

	return artists, nil
}

func (r *ArtistRepository) FindByID(ctx context.Context, id string) (*models.Artist, error) {
	if config.FirestoreClient == nil {
		return nil, fmt.Errorf("Firestore client is not initialized")
	}

	doc, err := r.c().Doc(id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get artist by ID: %w", err)
	}

	var a models.Artist
	if err := doc.DataTo(&a); err != nil {
		return nil, fmt.Errorf("failed to parse artist data: %w", err)
	}
	return &a, nil
}

func (r *ArtistRepository) Update(ctx context.Context, a *models.Artist) error {
	a.UpdatedAt = time.Now()
	_, err := r.c().Doc(a.ID).Set(ctx, a)
	return err
}

func (r *ArtistRepository) Delete(ctx context.Context, id string) error {
	_, err := r.c().Doc(id).Delete(ctx)
	return err
}
