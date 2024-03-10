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

	database.CreateConnection()

	// Close the database connection when the main function returns
	defer database.CloseConnection()

	// Ensure the required table exists
	database.CreateTicketInfoTable()

	e := echo.New()

	// Return the id for given name
	e.GET("/id", api.GetStationId)

	e.GET("/data", api.GetData)

	// Post a new ticket inspector
	e.POST("/newInspector", api.PostInspector)

	e.Start(":8080")
}
