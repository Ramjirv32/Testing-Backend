package pass

import (
	passCtrl "backend/controllers/pass"
	"backend/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterPassRoutes(router fiber.Router) {
	pass := router.Group("/pass")
	// POST /api/v1/pass/activate — authenticated users only
	pass.Post("/activate", middleware.Auth, passCtrl.ActivatePass)
}
