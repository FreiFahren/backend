package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/id", GetStation)

	e.Start(":8080")
}
