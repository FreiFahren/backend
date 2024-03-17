package main

import (
	"log"

	"github.com/FreiFahren/backend/api"
	"github.com/FreiFahren/backend/database"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create a new connection pool, for concurrency
	database.CreatePool()

	if err != nil {
		log.Fatal("Error while creating a pool :(")
	}

	e := echo.New()

	// Close the database connection when the main function returns
	defer database.ClosePool()

	// Ensure the required table exists
	database.CreateTicketInfoTable()

	// Return the id for given name
	e.GET("/id", api.GetStationId)

	// Return the last known inspectors 15 mins ago
	e.GET("/recent", api.GetRecentTicketInspectorInfo)

	// Return the name for given id
	e.GET("/station", api.GetStationName)

	// Return all stations with their id (used for suggestions on the frontend)
	e.GET("/list", api.GetAllStationsAndLines)

	// Post a new ticket inspector
	e.POST("/newInspector", api.PostInspector)

	e.Start(":8080")

	defer e.Close()
}
