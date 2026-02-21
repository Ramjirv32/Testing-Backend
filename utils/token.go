package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

func GenerateToken(userID string, isAdmin bool) string {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
	
		randomBytes = []byte(fmt.Sprintf("%064d", time.Now().UnixNano()))
	}
	// Use StdEncoding (no _ chars) so token splitting on "_" is unambiguous
	token := base64.StdEncoding.EncodeToString(randomBytes)

	prefix := "user"
	if isAdmin {
		prefix = "admin"
	}

	return fmt.Sprintf("%s_%s_%s", prefix, userID, token)
}

// ValidateToken parses a token and returns (userID, isAdminPrefix, error).
// SplitN(3) ensures the base64 segment (which may contain "=") is kept intact.
func ValidateToken(token string) (string, bool, error) {
	parts := strings.SplitN(token, "_", 3)
	if len(parts) < 3 {
		return "", false, fmt.Errorf("invalid token format")
	}
	prefix := parts[0]
	userID := parts[1]
	if prefix != "user" && prefix != "admin" {
		return "", false, fmt.Errorf("unknown token prefix")
	}
	if userID == "" {
		return "", false, fmt.Errorf("empty user ID in token")
	}
	return userID, prefix == "admin", nil
}
