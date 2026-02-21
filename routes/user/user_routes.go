package user

import (
	"backend/controllers/user"
	"backend/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterUserRoutes(router fiber.Router) {
	usersGroup := router.Group("/users")
	usersGroup.Get("/", user.GetAllUsers)
	usersGroup.Post("/", user.CreateUser)
	usersGroup.Get("/:id", user.GetUserByID)
	usersGroup.Put("/:id", middleware.Auth, user.UpdateUser)
}
