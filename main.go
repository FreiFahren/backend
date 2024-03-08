package main

import (
	"github.com/FreiFahren/backend/api"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/id", api.GetStationId)

	e.Start(":8080")
}
