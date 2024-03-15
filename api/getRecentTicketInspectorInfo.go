package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/FreiFahren/backend/database"
	. "github.com/FreiFahren/backend/structs"
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

func GetRecentTicketInspectorInfo(c echo.Context) error {

	// Get the latest ticket inspector information from the database
	TicketInfoList, err := database.GetLatestStationCoordinates()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// // We fill up TicketInfoList with historic data, if we have less than 10 current entries
	// if len(TicketInfoList) < 10 {
	// 	historicDataList, err := database.GetHistoricStations(time.Now())

	// 	if err != nil {
	// 		return c.JSON(http.StatusInternalServerError, err.Error())
	// 	}

	// 	for _, TicketInfo := range historicDataList {

	// 		fmt.Println("Adding historic data...")
	// 		if len(TicketInfoList) > 10 {
	// 			break
	// 		}

	// 		TicketInfoList = append(TicketInfoList, TicketInfo)
	// 		fmt.Println(TicketInfo.Station_ID)
	// 	}
	// }

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

		stationLat, stationLon, err := IdToCoordinates(cleanedStationId)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		// Get the names

		stationName, err := IdToStationName(cleanedStationId)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		directionName := ""
		var directionLat float64 = 0
		var directionLon float64 = 0

		if ticketInfo.Direction_ID.Valid {
			directionName, err = IdToStationName(cleanedDirectionId)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
			directionLat, directionLon, err = IdToCoordinates(cleanedDirectionId)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
		}

		// Create a new TicketInspector struct and append it to the slice

		TicketInspectorInfo := TicketInspector{
			Timestamp: ticketInfo.Timestamp,
			Station: Station{
				ID:   cleanedStationId,
				Name: stationName,
				Coordinates: Coordinates{
					Latitude:  stationLat,
					Longitude: stationLon,
				},
			},
			Direction: Station{
				ID:   cleanedDirectionId,
				Name: directionName,
				Coordinates: Coordinates{
					Latitude:  directionLat,
					Longitude: directionLon,
				},
			},
			Line: cleanedLine,
		}

		TicketInspectorList = append(TicketInspectorList, TicketInspectorInfo)

	}

	return c.JSONPretty(http.StatusOK, TicketInspectorList, "  ")
}
