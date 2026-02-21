package booking

import (
	"backend/controllers/booking"
	"backend/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterBookingRoutes(router fiber.Router) {
	bookingsGroup := router.Group("/bookings", middleware.Auth)
	bookingsGroup.Post("/play", booking.CreatePlayBooking)
	bookingsGroup.Post("/dining", booking.CreateDiningBooking)
	bookingsGroup.Post("/event", booking.CreateEventBooking)
	bookingsGroup.Get("/", booking.GetUserBookings)
	bookingsGroup.Get("/:id", booking.GetBookingByID)
	bookingsGroup.Delete("/:id", booking.CancelBooking)
}
