package utils

import (
	"sync/atomic"

	"github.com/google/uuid"
)

var SeqID int64 = 0

func GenerateUUIDv7() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.New().String() // Fallback to v4 if v7 fails
	}
	return id.String()
}

func GetNextSeqID() int64 {
	return atomic.AddInt64(&SeqID, 1)
}
