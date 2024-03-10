package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/FreiFahren/backend/database"
	"github.com/labstack/echo/v4"
)

func IdToCoordinates(id string) (float64, float64, error) {
	stations, err := ReadFromFile("data/stations.json")
	if err != nil {
		return 0, 0, err
	}

	station, ok := stations[id]
	if !ok {
		return 0, 0, fmt.Errorf("station ID %s not found", id)
	}

	return station.Coordinates.Latitude, station.Coordinates.Longitude, nil
}

func GetData(c echo.Context) error {

	// Get the latest ticket inspector information from the database
	TicketInfoList, err := database.GetLatestInfo()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	historicData, err := database.GetHistoricStations(time.Now())

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// go through all the id's, remove the \n and get the coordinates
	// then appends it to the slice
	coordinates := make([][]float64, 0)

	for _, ticketInfo := range TicketInfoList {

		cleanStr := strings.ReplaceAll(ticketInfo.Station_ID, "\n", "")

		latitude, longitude, err := IdToCoordinates(cleanStr)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		coordinates = append(coordinates, []float64{latitude, longitude})
	}

	// if we have historic data, we use it
	if len(historicData) != 0 {
		print("Using historic data")

		// Add historic data to the coordinates slice
		for _, id := range historicData {

			cleanStr := strings.ReplaceAll(id, "\n", "")

			latitude, longitude, err := IdToCoordinates(cleanStr)

			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
			coordinates = append(coordinates, []float64{latitude, longitude})
		}
	}

	return c.JSON(http.StatusOK, coordinates)
}
