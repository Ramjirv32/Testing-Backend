package main

import (
	"log"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"

	"backend/config"
	"backend/routes"
	"backend/tasks"
	"backend/utils"
)

func main() {
	cfg := config.LoadConfig()

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		AppName:     "Website Backend",
		BodyLimit:   10 * 1024 * 1024, // 10MB
	})

	config.InitFirebase(cfg)
	utils.InitEmail(cfg)
	tasks.StartEmailWorker()

	log.Println("✅ Firebase (Auth, Firestore, Storage) and Email initialized successfully")

	// Dynamic CORS Middleware - Critical for cross-domain testing with credentials
	app.Use(func(c fiber.Ctx) error {
		origin := c.Get("Origin")
		if origin != "" {
			c.Set("Access-Control-Allow-Origin", origin)
			c.Set("Access-Control-Allow-Credentials", "true")
		} else {
			c.Set("Access-Control-Allow-Origin", "*")
		}
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin, X-Requested-With, X-Firebase-Token")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(204)
		}
		return c.Next()
	})

	routes.SetupRoutes(app)

	// Ensure port starts with colon if it's just a number
	port := cfg.Port
	if port != "" && port[0] != ':' {
		port = ":" + port
	}

	log.Printf("🚀 Server starting on port %s (Environment: %s)\n", port, cfg.Env)

	if err := app.Listen(port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
