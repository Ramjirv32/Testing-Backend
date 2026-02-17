package utils

import (
	"github.com/gofiber/fiber/v3"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(c fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Status:  statusCode,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(Response{
		Status:  statusCode,
		Message: message,
	})
}
