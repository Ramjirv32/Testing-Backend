package routes

import (
	"github.com/gofiber/fiber/v3"

	"backend/controllers"
	"backend/middleware"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	api.Post("/upload", middleware.Auth, controllers.UploadFile)

	// Organizer Auth
	orgAuth := api.Group("/organizer/auth")
	orgAuth.Post("/register", controllers.OrganizerRegister)
	orgAuth.Post("/login", controllers.OrganizerLogin)
	orgAuth.Post("/google", controllers.OrganizerGoogleLogin)
	orgAuth.Post("/verify-otp", controllers.OrganizerVerifyOTP)
	orgAuth.Post("/resend-otp", controllers.ResendOrganizerOTP)

	health := api.Group("/health")
	health.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	auth := api.Group("/auth")
	auth.Post("/send-otp", controllers.SendOTP)
	auth.Post("/login", controllers.Login)
	auth.Post("/logout", controllers.Logout)
	auth.Get("/profile", middleware.Auth, controllers.GetProfile)
	auth.Put("/profile", middleware.Auth, controllers.UpdateProfile)
	auth.Post("/email/send-otp", middleware.Auth, controllers.SendEmailOTP)
	auth.Post("/email/verify", middleware.Auth, controllers.VerifyEmail)

	bookings := api.Group("/bookings", middleware.Auth)
	bookings.Post("/play", controllers.CreatePlayBooking)
	bookings.Post("/dining", controllers.CreateDiningBooking)
	bookings.Get("/", controllers.GetUserBookings)
	bookings.Get("/:id", controllers.GetBookingByID)
	bookings.Delete("/:id", controllers.CancelBooking)

	admin := api.Group("/admin", middleware.Auth)
	admin.Delete("/play-bookings/all", controllers.DeleteAllPlayBookings)
	admin.Delete("/dining-bookings/all", controllers.DeleteAllDiningBookings)
	admin.Delete("/users/all", controllers.DeleteAllUsers)
	admin.Get("/partners", controllers.GetEventPosters)
	admin.Patch("/partners/:id/approve", controllers.ApproveEventPoster)
	admin.Put("/partners/:id", controllers.UpdatePartnerProfile)

	partners := api.Group("/partners", middleware.Auth)
	partners.Post("/verify", controllers.SubmitVerification)
	partners.Get("/my-status", controllers.GetMyVerificationStatus)

	events := api.Group("/events")
	events.Get("/organizer/my", middleware.Auth, controllers.GetOrganizerEvents)
	events.Get("/", controllers.GetAllEvents)
	events.Get("/:id", controllers.GetEventByID)
	events.Post("/", middleware.Auth, controllers.CreateEvent)
	events.Put("/:id", middleware.Auth, controllers.UpdateEvent)

	play := api.Group("/play")
	play.Get("/organizer/my", middleware.Auth, controllers.GetOrganizerPlayVenues)
	play.Get("/", controllers.GetAllPlayVenues)
	play.Get("/:slug", controllers.GetPlayVenueBySlug)
	play.Get("/id/:id", controllers.GetPlayVenueByID)
	play.Post("/", middleware.Auth, controllers.CreatePlayVenue)
	play.Post("/seed", controllers.SeedPlayVenues)

	dining := api.Group("/dining")
	dining.Get("/organizer/my", middleware.Auth, controllers.GetOrganizerDiningVenues)
	dining.Get("/", controllers.GetAllDiningVenues)
	dining.Get("/:slug", controllers.GetDiningVenueBySlug)
	dining.Get("/id/:id", controllers.GetDiningVenueByID)
	dining.Post("/", middleware.Auth, controllers.CreateDiningVenue)
	dining.Post("/seed", controllers.SeedDiningVenues)

	users := api.Group("/users")
	users.Get("/", controllers.GetAllUsers)
	users.Post("/", controllers.CreateUser)
	users.Get("/:id", controllers.GetUserByID)
	users.Put("/:id", middleware.Auth, controllers.UpdateUser)

	artists := api.Group("/artists")
	artists.Get("/", controllers.GetAllArtists)
	artists.Get("/:id", controllers.GetArtistByID)
	artists.Post("/", middleware.Auth, controllers.CreateArtist)
	artists.Post("/seed", controllers.SeedArtists)

	offers := api.Group("/offers")
	offers.Get("/", controllers.GetAllOffers)
	offers.Get("/user/:userId", controllers.GetUserOffers)
	offers.Post("/", middleware.Auth, controllers.CreateOffer)
	offers.Post("/seed", controllers.SeedOffers)
}
