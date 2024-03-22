package main

import (
	"log"
	"net/http"

	"github.com/FreiFahren/backend/api"
	"github.com/FreiFahren/backend/database"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme/autocert"
)

type (
	Host struct {
		Echo *echo.Echo
	}
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
	// Hosts
	hosts := map[string]*Host{}

	apiHOST := echo.New()
	apiHOST.Use(middleware.Logger())
	apiHOST.Use(middleware.Recover())

	hosts["api.freifahren.org:443"] = &Host{apiHOST}

	apiHOST.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "API")
	})

	// Close the database connection when the main function returns
	defer database.ClosePool()

	// Ensure the required table exists
	database.CreateTicketInfoTable()

	// Return the id for given name
	apiHOST.GET("/id", api.GetStationId)

	// Return the last known inspectors 15 mins ago
	apiHOST.GET("/recent", api.GetRecentTicketInspectorInfo)

	// Return the name for given id
	apiHOST.GET("/station", api.GetStationName)

	// Return all stations with their id (used for suggestions on the frontend)
	apiHOST.GET("/list", api.GetAllStationsAndLines)

	// Post a new ticket inspector
	apiHOST.POST("/newInspector", api.PostInspector)

	site := echo.New()
	site.Use(middleware.Logger())
	site.Use(middleware.Recover())

	hosts["freifahren.org:443"] = &Host{site}

	site.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Website")
	})

	// Server
	e := echo.New()
	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		host := hosts[req.Host]

		if host == nil {
			err = echo.ErrNotFound
		} else {
			host.Echo.ServeHTTP(res, req)
		}

		return
	})

	// Start the server with AutoTLS
	e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")

	e.Logger.Fatal(e.StartAutoTLS(":443"))

	defer apiHOST.Close()
}
