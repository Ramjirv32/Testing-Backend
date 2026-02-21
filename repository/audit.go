package repository

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"

	"backend/config"
	"backend/models"
	"backend/utils"
)

type AuditRepository struct{}

func (r *AuditRepository) c() *firestore.CollectionRef {
	if config.FirestoreClient == nil {
		panic("FirestoreClient is not initialized.")
	}
	return config.FirestoreClient.Collection("audit_logs")
}

func NewAuditRepository() *AuditRepository {
	return &AuditRepository{}
}

func (r *AuditRepository) Log(ctx context.Context, log *models.AuditLog) error {
	log.ID = utils.GenerateUUIDv7()
	log.CreatedAt = time.Now()
	_, err := r.c().Doc(log.ID).Set(ctx, log)
	return err
}
