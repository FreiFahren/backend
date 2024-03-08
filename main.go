package main

import (
	"github.com/FreiFahren/backend/api"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// return the id for given name
	e.GET("/id", api.GetStationId)

	e.Start(":8085")
}
