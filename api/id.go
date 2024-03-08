package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

type Station struct {
	Name string `json:"name"`
}

func GetStationId(c echo.Context) error {
	name := c.QueryParam("name")

	// Convert the input to lowercase and remove all spaces
	name = strings.ToLower(strings.ReplaceAll(name, " ", ""))

	jsonFile, err := os.Open("data/stations.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var stations map[string]Station
	json.Unmarshal(byteValue, &stations)

	for id, station := range stations {
		// Convert the station name to lowercase and remove all spaces before comparing
		stationName := strings.ToLower(strings.ReplaceAll(station.Name, " ", ""))
		if stationName == name {
			return c.JSON(http.StatusOK, id)
		}
	}

	return echo.ErrNotFound
}
