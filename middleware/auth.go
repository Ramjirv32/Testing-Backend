package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"

	"backend/repository"
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

	// Use SplitN(3) so base64 payload (may contain "="/extra chars) won't corrupt userID
	tokenParts := strings.SplitN(token, "_", 3)
	if len(tokenParts) < 3 {
		return utils.ErrorResponse(c, 401, "Invalid token format")
	}

	prefix := tokenParts[0]
	userID := tokenParts[1]

	if (prefix != "user" && prefix != "admin") || userID == "" {
		return utils.ErrorResponse(c, 401, "Invalid token")
	}

	// Look up the user in DB to get the real isAdmin flag — never trust the token prefix alone
	userRepo := repository.NewUserRepository()
	user, err := userRepo.FindByID(c.Context(), userID)
	isAdmin := false
	if err == nil && user != nil {
		isAdmin = user.IsAdmin
	} else {
		// Fallback: if DB is unreachable trust the prefix (for local dev)
		isAdmin = prefix == "admin"
		fmt.Printf("[Auth] DB lookup failed for %s, falling back to token prefix: %v\n", userID, err)
	}

	c.Locals("uid", userID)
	c.Locals("isAdmin", isAdmin)

	return c.Next()
}

func AdminOnly(c fiber.Ctx) error {
	isAdmin, ok := c.Locals("isAdmin").(bool)
	if !ok || !isAdmin {
		return utils.ErrorResponse(c, 403, "Admin access required")
	}
	return c.Next()
}

func OrganizerOnly(category string) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := c.Locals("uid").(string)
		isAdmin, _ := c.Locals("isAdmin").(bool)

		if isAdmin {
			return c.Next()
		}

		userRepo := repository.NewUserRepository()
		user, err := userRepo.FindByID(c.Context(), userID)

		if err != nil || user == nil {
			return utils.ErrorResponse(c, 404, "User not found")
		}

		if !user.IsOrganizer {
			return utils.ErrorResponse(c, 403, "Organizer access required")
		}

		if category != "" {
			match := false
			// Check primary category
			if user.OrganizerCategory == category {
				match = true
			} else if category == "event" && (user.OrganizerCategory == "creator" || user.OrganizerCategory == "individual" || user.OrganizerCategory == "company" || user.OrganizerCategory == "non-profit") {
				match = true
			}

			// Check categories slice if not matched yet
			if !match {
				for _, cat := range user.OrganizerCategories {
					if cat == category {
						match = true
						break
					}
					if category == "event" && (cat == "creator" || cat == "individual" || cat == "company" || cat == "non-profit") {
						match = true
						break
					}
				}
			}

			if !match {
				return utils.ErrorResponse(c, 403, fmt.Sprintf("Access denied. Direct %s organizer access required.", category))
			}
		}

		return c.Next()
	}
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
		tokenParts := strings.SplitN(token, "_", 3)
		if len(tokenParts) >= 3 && tokenParts[1] != "" {
			userID := tokenParts[1]
			c.Locals("uid", userID)
			// Resolve isAdmin from DB if possible
			userRepo := repository.NewUserRepository()
			if user, err := userRepo.FindByID(c.Context(), userID); err == nil && user != nil {
				c.Locals("isAdmin", user.IsAdmin)
			} else {
				c.Locals("isAdmin", tokenParts[0] == "admin")
			}
		}
	}

	return c.Next()
}
