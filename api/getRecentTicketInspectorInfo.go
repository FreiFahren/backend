package api

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/FreiFahren/backend/database"
	"github.com/FreiFahren/backend/structs"
	"github.com/labstack/echo/v4"
)

func GetRecentTicketInspectorInfo(c echo.Context) error {
	// Check if the data has been modified since the provided time
	modifiedSince, err := CheckIfModifiedSince(c)
	if err != nil {
		fmt.Printf("Error checking if the data has been modified: %v\n", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if modifiedSince {
		// Return 304 Not Modified if the data hasn't been modified since the provided time
		return c.NoContent(http.StatusNotModified)
	}

	// Proceed with fetching and processing the data if it was modified
	// or if the If-Modified-Since header was not provided
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
		ticketInspector, err := constructTicketInspectorInfo(ticketInfo)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		ticketInspectorList = append(ticketInspectorList, ticketInspector)
	}

	filteredTicketInspectorList := RemoveDuplicateStations(ticketInspectorList)

	return c.JSONPretty(http.StatusOK, filteredTicketInspectorList, "  ")
}

func CheckIfModifiedSince(c echo.Context) (bool, error) {
    databaseLastModified, err := database.GetLatestUpdateTime()
    if err != nil {
        return false, err
    }

    ifModifiedSince := c.Request().Header.Get("If-Modified-Since")

    // If the header is empty, proceed with fetching the data
    if ifModifiedSince == "" {
        return false, nil
    }

    // Use time.RFC3339 to parse ISO 8601 format
    requestedModificationTime, err := time.Parse(time.RFC3339, ifModifiedSince)
    if err != nil {
        return false, fmt.Errorf("error parsing If-Modified-Since header: %v", err)
    }

    // Check if the database last modified time is after the requested modification time
    if !databaseLastModified.After(requestedModificationTime) {
        return true, nil
    }
    return false, nil
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

func RemoveDuplicateStations(ticketInspectorList []structs.TicketInspector) []structs.TicketInspector {
	uniqueStations := make(map[string]structs.TicketInspector)
	for _, ticketInspector := range ticketInspectorList {
		stationID := ticketInspector.Station.ID
		if existingInspector, ok := uniqueStations[stationID]; !ok || ticketInspector.Timestamp.After(existingInspector.Timestamp) {
			uniqueStations[stationID] = ticketInspector
		}
	}

	filteredTicketInspectorList := make([]structs.TicketInspector, 0, len(uniqueStations))
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
		return structs.TicketInspector{}, err
	}

	stationName, err := IdToStationName(cleanedStationId)
	if err != nil {
		return structs.TicketInspector{}, err
	}

	directionName, directionLat, directionLon := "", float64(0), float64(0)
	if ticketInfo.Direction_ID.Valid {
		directionName, err = IdToStationName(cleanedDirectionId)
		if err != nil {
			return structs.TicketInspector{}, err
		}
		directionLat, directionLon, err = IdToCoordinates(cleanedDirectionId)
		if err != nil {
			return structs.TicketInspector{}, err
		}
	}

	ticketInspectorInfo := structs.TicketInspector{
		Timestamp: ticketInfo.Timestamp,
		Station: structs.Station{
			ID:          cleanedStationId,
			Name:        stationName,
			Coordinates: structs.Coordinates{Latitude: stationLat, Longitude: stationLon},
		},
		Direction: structs.Station{
			ID:          cleanedDirectionId,
			Name:        directionName,
			Coordinates: structs.Coordinates{Latitude: directionLat, Longitude: directionLon},
		},
		Line: cleanedLine,
	}
	return ticketInspectorInfo, nil
}
