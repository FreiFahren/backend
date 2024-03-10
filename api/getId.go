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
	ID          string `json:"id"`
	Name        string `json:"name"`
	Coordinates struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"coordinates"`
}

func ReadFromFile(filePath string) (map[string]Station, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var stations map[string]Station
	if err := json.Unmarshal(byteValue, &stations); err != nil {
		return nil, err
	}

	return stations, nil
}

func FindStationId(name string, stations map[string]Station) (string, bool) {
	name = strings.ToLower(strings.ReplaceAll(name, " ", ""))
	for id, station := range stations {
		stationName := strings.ToLower(strings.ReplaceAll(station.Name, " ", ""))
		if stationName == name {
			return id, true
		}
	}
	return "", false
}

func GetStationId(c echo.Context) error {
	name := c.QueryParam("name")

	stations, err := ReadFromFile("data/stations.json")
	if err != nil {
		return err
	}

	id, found := FindStationId(name, stations)
	if found {
		return c.JSON(http.StatusOK, id)
	}

	return echo.ErrNotFound
}
