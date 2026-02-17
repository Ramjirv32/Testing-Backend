package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
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
	parts := strings.Split(token, "_")
	if len(parts) < 2 {
		return "", false, fmt.Errorf("invalid token format")
	}
	userID := parts[1]
	isAdmin := parts[0] == "admin"
	return userID, isAdmin, nil
}
