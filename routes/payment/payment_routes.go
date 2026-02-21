package paymentroutes

import (
	"backend/controllers/payment"
	"backend/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterPaymentRoutes(router fiber.Router) {
	payGroup := router.Group("/payment", middleware.Auth)
	payGroup.Post("/create-order", payment.CreateOrder)
	payGroup.Post("/verify", payment.VerifyPayment)
}
