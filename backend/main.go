package main

import (
	"fmt"
	"log"
	"os"

	"github.com/3blakejohnson/survivor-fantasy-league/backend/controllers"
	"github.com/3blakejohnson/survivor-fantasy-league/backend/dao"
	"github.com/3blakejohnson/survivor-fantasy-league/backend/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	supa "github.com/nedpals/supabase-go"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	db.ConnectDB()
	defer db.CloseDB()

	dm := dao.NewDAOManager(db.DB)

	// Initialize Supabase
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")
	if supabaseURL == "" || supabaseKey == "" {
		log.Fatal("Supabase credentials not found in environment variables")
	}
	supabase := supa.CreateClient(supabaseURL, supabaseKey)
	fmt.Printf("using supabase %v", supabase)

	// Create a new Fiber instance
	app := fiber.New()
	app.Use(logger.New())

	// Serve static files (CSS, images, etc.)
	app.Static("/static", "./static")

	// Serve HTML Homepage
	app.Get("/", func(c *fiber.Ctx) error {
		return templates.HomePage().Render(c.Context(), c.Response().BodyWriter()) // Render Templ component
	})

	app.Get("/user/:id", controllers.GetUser(dm))
	app.Get("/episode/:season/:episode", controllers.GetEpisode(dm))
	app.Post("/episode", controllers.CreateEpisode(dm))

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080
	}
	log.Fatal(app.Listen(":" + port))
}
