package main

import (
	"fmt"
	"log"
	"os"

	"github.com/3blakejohnson/survivor-fantasy-league/backend/controllers"
	"github.com/3blakejohnson/survivor-fantasy-league/backend/dao"
	"github.com/3blakejohnson/survivor-fantasy-league/backend/db"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	supa "github.com/nedpals/supabase-go"
)

// func dbConnect() {
// 	// Get connection string from environment variable
// 	connStr := os.Getenv("DATABASE_URL")
// 	if connStr == "" {
// 		log.Fatal("DATABASE_URL is not set in .env file")
// 	}

// 	db, err := sql.Open("postgres", connStr)
// 	if err != nil {
// 		panic("Couldn't establish connection to Database")
// 	}
// 	defer db.Close()

// 	err = db.Ping()
// 	if err != nil {
// 		panic("Ping to database failed")
// 	} else {
// 		fmt.Println("Succesfully Established connection to database!")
// 	}
// }

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	db.ConnectDB()
	defer db.CloseDB()

	dm := dao.NewDAOManager(db.DB)

	// // Example: add user to db
	// fmt.Println("Starting user creation")
	// userDAO := dao.NewUserDAO(db.DB)
	// newUser := models.User{
	// 	Username:     "3blakejohnson",
	// 	PasswordHash: "password",
	// 	FirstName:    "Blake",
	// 	LastName:     "Johnson",
	// }
	// err = userDAO.Create(newUser)
	// if err != nil {
	// 	log.Fatal("Error creating user")
	// }

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
