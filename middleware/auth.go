package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"

	"backend/utils"
)

func Auth(c fiber.Ctx) error {
	var token string
	authHeader := c.Get("Authorization")

	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			token = parts[1]
		} else {
			return utils.ErrorResponse(c, 401, "Invalid authorization format")
		}
	} else {
	
		token = c.Cookies("authToken")
	}

	if token == "" {
		return utils.ErrorResponse(c, 401, "Authentication required")
	}

	if token == "" {
		return utils.ErrorResponse(c, 401, "Authentication required")
	}

	tokenParts := strings.Split(token, "_")
	if len(tokenParts) < 2 {
		return utils.ErrorResponse(c, 401, "Invalid token format")
	}

	userID := tokenParts[1]
	isAdmin := tokenParts[0] == "admin"

	c.Locals("uid", userID)
	c.Locals("isAdmin", isAdmin)

	return c.Next()
}

func OptionalAuth(c fiber.Ctx) error {
	var token string
	authHeader := c.Get("Authorization")

	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			token = parts[1]
		}
	} else {
		token = c.Cookies("authToken")
	}

	if token != "" {
		tokenParts := strings.Split(token, "_")
		if len(tokenParts) >= 2 {
			c.Locals("uid", tokenParts[1])
			c.Locals("isAdmin", tokenParts[0] == "admin")
		}
	}

	return c.Next()
}
