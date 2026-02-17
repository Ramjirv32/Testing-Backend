package main

import (
	"log"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"

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
	config.InitFirestore(cfg)
	config.InitStorage(config.FirebaseApp)
	utils.InitEmail(cfg)
	tasks.StartEmailWorker()

	log.Println("✅ Firebase, Firestore, and Email initialized successfully")

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Allow frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	routes.SetupRoutes(app)

	log.Printf("🚀 Server starting on port %s\n", cfg.Port)

	if err := app.Listen(cfg.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
