package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/FreiFahren/backend/database"
	"github.com/labstack/echo/v4"
)

// This is the struct that we will use to store the coordinates and the direction id and line
// to return to the frontend
type TicketInspector struct {
	Coordinates []float64
	StationID   string
	DirectionID string
	Line        string
}

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

	TicketInspectorList := make([]TicketInspector, 0)

	// go through all the id's, remove the \n and get the coordinates
	// then appends it to the slice
	for _, ticketInfo := range TicketInfoList {

		cleanedStationId := strings.ReplaceAll(ticketInfo.Station_ID, "\n", "")
		cleanedDirectionId := ""
		cleanedLine := ""

		if ticketInfo.Direction_ID.Valid {
			cleanedDirectionId = strings.ReplaceAll(ticketInfo.Direction_ID.String, "\n", "")
		}

		if ticketInfo.Line.Valid {
			cleanedLine = strings.ReplaceAll(ticketInfo.Line.String, "\n", "")
		}

		latitude, longitude, err := IdToCoordinates(cleanedStationId)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		TicketInspectorInfo := TicketInspector{
			StationID:   cleanedStationId,
			Coordinates: []float64{latitude, longitude},
			DirectionID: cleanedDirectionId,
			Line:        cleanedLine,
		}

		TicketInspectorList = append(TicketInspectorList, TicketInspectorInfo)

	}

	return c.JSONPretty(http.StatusOK, TicketInspectorList, "  ")
}
