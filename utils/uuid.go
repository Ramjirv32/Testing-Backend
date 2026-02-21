package utils

import (
	"sync/atomic"

	"github.com/google/uuid"
)

var SeqID int64 = 0

func GenerateUUIDv7() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.New().String() 
	}
	return id.String()
}

func GetNextSeqID() int64 {
	return atomic.AddInt64(&SeqID, 1)
}
