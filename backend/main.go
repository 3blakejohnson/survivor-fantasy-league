package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	supa "github.com/nedpals/supabase-go"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

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

	// Health check route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Survivor Fantasy League API is running!")
	})

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080
	}
	log.Fatal(app.Listen(":" + port))
}
