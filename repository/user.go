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

type UserRepository struct{}

func (r *UserRepository) c() *firestore.CollectionRef {
	if config.FirestoreClient == nil {
		panic("FirestoreClient is not initialized. Check InitFirebase logs.")
	}
	return config.FirestoreClient.Collection("users")
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	_, err := r.c().Doc(user.ID).Set(ctx, user)
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	if config.FirestoreClient == nil {
		return nil, fmt.Errorf("Firestore client is not initialized")
	}

	doc, err := r.c().Doc(id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return nil, fmt.Errorf("failed to parse user data: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (*models.User, error) {
	if config.FirestoreClient == nil {
		return nil, fmt.Errorf("Firestore client is not initialized")
	}

	iter := r.c().Where("phone", "==", phone).Limit(1).Documents(ctx)
	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user by phone: %w", err)
	}

	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return nil, fmt.Errorf("failed to parse user data: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	if config.FirestoreClient == nil {
		return nil, fmt.Errorf("Firestore client is not initialized")
	}

	iter := r.c().Where("email", "==", email).Limit(1).Documents(ctx)
	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user by email: %w", err)
	}

	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return nil, fmt.Errorf("failed to parse user data: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	_, err := r.c().Doc(user.ID).Set(ctx, user)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.c().Doc(id).Delete(ctx)
	return err
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	iter := r.c().Documents(ctx)
	var users []*models.User

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate users: %w", err)
		}

		var user models.User
		if err := doc.DataTo(&user); err != nil {
			continue
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *UserRepository) DeleteAll(ctx context.Context) error {
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
