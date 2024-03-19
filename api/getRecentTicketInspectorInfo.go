package api

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/FreiFahren/backend/database"
	"github.com/FreiFahren/backend/structs"
	. "github.com/FreiFahren/backend/structs"
	"github.com/labstack/echo/v4"
)

func GetRecentTicketInspectorInfo(c echo.Context) error {
	ticketInfoList, err := database.GetLatestStationCoordinates()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	ticketInfoList, err = FetchAndAddHistoricData(ticketInfoList)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	ticketInspectorList := []structs.TicketInspector{}
	for _, ticketInfo := range ticketInfoList {
		//
		ticketInspector, err := constructTicketInspectorInfo(ticketInfo)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		ticketInspectorList = append(ticketInspectorList, ticketInspector)
	}

	filteredTicketInspectorList := RemoveDuplicateStations(ticketInspectorList)

	return c.JSONPretty(http.StatusOK, filteredTicketInspectorList, "  ")
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

func RemoveDuplicateStations(ticketInspectorList []TicketInspector) []TicketInspector {
	uniqueStations := make(map[string]TicketInspector)
	for _, ticketInspector := range ticketInspectorList {
		stationID := ticketInspector.Station.ID
		if existingInspector, ok := uniqueStations[stationID]; !ok || ticketInspector.Timestamp.After(existingInspector.Timestamp) {
			uniqueStations[stationID] = ticketInspector
		}
	}

	filteredTicketInspectorList := make([]TicketInspector, 0, len(uniqueStations))
	for _, ticketInspector := range uniqueStations {
		filteredTicketInspectorList = append(filteredTicketInspectorList, ticketInspector)
	}

	// Sort the list by timestamp, then by station name if timestamps are equal
	sort.Slice(filteredTicketInspectorList, func(i, j int) bool {
		if filteredTicketInspectorList[i].Timestamp.Equal(filteredTicketInspectorList[j].Timestamp) {
			return filteredTicketInspectorList[i].Station.Name < filteredTicketInspectorList[j].Station.Name
		}
		return filteredTicketInspectorList[i].Timestamp.After(filteredTicketInspectorList[j].Timestamp)
	})

	return filteredTicketInspectorList
}

func FetchAndAddHistoricData(ticketInfoList []database.TicketInfo) ([]database.TicketInfo, error) {
	if len(ticketInfoList) < 10 {
		historicDataList, err := database.GetHistoricStations(time.Now())
		if err != nil {
			return nil, err
		}

		for _, ticketInfo := range historicDataList {
			if len(ticketInfoList) >= 10 {
				break
			}
			ticketInfoList = append(ticketInfoList, ticketInfo)
		}
	}
	return ticketInfoList, nil
}

func constructTicketInspectorInfo(ticketInfo database.TicketInfo) (structs.TicketInspector, error) {
	cleanedStationId := strings.ReplaceAll(ticketInfo.Station_ID, "\n", "")
	cleanedDirectionId := strings.ReplaceAll(ticketInfo.Direction_ID.String, "\n", "")
	cleanedLine := strings.ReplaceAll(ticketInfo.Line.String, "\n", "")

	stationLat, stationLon, err := IdToCoordinates(cleanedStationId)
	if err != nil {
		return TicketInspector{}, err
	}

	stationName, err := IdToStationName(cleanedStationId)
	if err != nil {
		return TicketInspector{}, err
	}

	directionName, directionLat, directionLon := "", float64(0), float64(0)
	if ticketInfo.Direction_ID.Valid {
		directionName, err = IdToStationName(cleanedDirectionId)
		if err != nil {
			return TicketInspector{}, err
		}
		directionLat, directionLon, err = IdToCoordinates(cleanedDirectionId)
		if err != nil {
			return TicketInspector{}, err
		}
	}

	ticketInspectorInfo := TicketInspector{
		Timestamp: ticketInfo.Timestamp,
		Station: Station{
			ID:          cleanedStationId,
			Name:        stationName,
			Coordinates: Coordinates{Latitude: stationLat, Longitude: stationLon},
		},
		Direction: Station{
			ID:          cleanedDirectionId,
			Name:        directionName,
			Coordinates: Coordinates{Latitude: directionLat, Longitude: directionLon},
		},
		Line: cleanedLine,
	}
	return ticketInspectorInfo, nil
}
