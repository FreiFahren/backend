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
	TicketInfoList, err := database.GetLatestStationCoordinates()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// We fill up TicketInfoList with historic data, if we have less than 10 current entries
	if len(TicketInfoList) < 10 {
		historicDataList, err := database.GetHistoricStations(time.Now())

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		for _, TicketInfo := range historicDataList {

			fmt.Println("Adding historic data...")
			if len(TicketInfoList) > 10 {
				break
			}

			TicketInfoList = append(TicketInfoList, TicketInfo)
			fmt.Println(TicketInfo.Station_ID)
		}
	}

	coordinates := make([][]float64, 0)

	// go through all the id's, remove the \n and get the coordinates
	// then appends it to the slice
	for _, ticketInfo := range TicketInfoList {

		cleanedStationId := strings.ReplaceAll(ticketInfo.Station_ID, "\n", "")

		latitude, longitude, err := IdToCoordinates(cleanedStationId)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		coordinates = append(coordinates, []float64{latitude, longitude})
	}

	return c.JSON(http.StatusOK, coordinates)
}
