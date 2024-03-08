package main

import (
	"github.com/FreiFahren/backend/api/id"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/station", id.GetStation)

	e.Start(":8080")
}
