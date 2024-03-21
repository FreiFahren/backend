package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func IdToStationName(id string) (string, error) {

	stations, err := ReadFromFile("data/Stations.json")
	if err != nil {
		return "", err
	}

	station, ok := stations[id]
	if !ok {
		return "", fmt.Errorf("station ID %s not found", id)
	}

	return station.Name, nil
}

func GetStationName(c echo.Context) error {
	id := c.QueryParam("id")

	id, err := IdToStationName(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, id)
}
