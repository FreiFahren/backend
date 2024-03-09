package main

import (
	"github.com/FreiFahren/backend/api"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// return the id for given name
	e.GET("/id", api.GetStationId)

	// post a new ticket inspector
	e.POST("/newInspector", api.PostInspector)

	e.Start(":8080")
}
