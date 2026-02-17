package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateToken(userID string, isAdmin bool) string {
	randomBytes := make([]byte, 32)
	rand.Read(randomBytes)
	token := base64.URLEncoding.EncodeToString(randomBytes)

	prefix := "user"
	if isAdmin {
		prefix = "admin"
	}

	return fmt.Sprintf("%s_%s_%s", prefix, userID, token)
}

func ValidateToken(token string) (string, bool, error) {
	return "", false, nil
}
