package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type Station struct {
	Name        string            `json:"name"`
	Coordinates map[string]string `json:"coordinates"`
}

func GetStation(c echo.Context) error {
	name := c.QueryParam("name")

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
		if station.Name == name {
			return c.JSON(http.StatusOK, id)
		}
	}

	return echo.ErrNotFound
}
