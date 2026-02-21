package ai

import (
	"backend/controllers/ai"

	"github.com/gofiber/fiber/v3"
)

func RegisterAIRoutes(router fiber.Router) {
	aiGroup := router.Group("/ai")
	aiGroup.Post("/chat", ai.HandleAIChat)
}
