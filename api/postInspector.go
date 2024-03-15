package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/FreiFahren/backend/database"
	"github.com/labstack/echo/v4"
)

type InspectorRequest struct {
	Line        string `json:"line"`
	StationName string `json:"station"`
	Direction   string `json:"direction"`
}

type ResponseData struct {
	Line      string  `json:"line"`
	Station   Station `json:"station"`
	Direction Station `json:"direction"`
}

func PostInspector(c echo.Context) error {
	var req InspectorRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	// Check if all parameters are empty
	if req.Line == "" && req.StationName == "" && req.DirectionName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "At least one of 'line', 'station', or 'direction' must be provided")
	}

	data, err := processRequestData(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, data)
}

func processRequestData(req InspectorRequest) (*ResponseData, error) {
	stations, err := ReadFromFile("data/stations.json")
	if err != nil {
		return nil, err
	}

	data := &ResponseData{}

	// Use pointers for all fields that can be empty and thus should be inserted as NULL.
	var linePtr, stationNamePtr, stationIDPtr, directionNamePtr, directionIDPtr *string

	// Assign the line pointer if the line is not an empty string.
	if req.Line != "" {
		linePtr = &req.Line
		data.Line = req.Line // Assign to data for response.
	}

	// Only assign other pointers if the value is found and not an empty string.
	if req.StationName != "" {
		if stationID, found := FindStationId(req.StationName, stations); found {
			stationNamePtr = &req.StationName
			stationIDPtr = &stationID
			data.Station = Station{Name: req.StationName, ID: stationID}
		} else {
			return nil, fmt.Errorf("Station not found")
		}
	}

	if req.DirectionName != "" {
		if directionID, found := FindStationId(req.DirectionName, stations); found {
			directionNamePtr = &req.DirectionName
			directionIDPtr = &directionID
			data.Direction = Station{Name: req.DirectionName, ID: directionID}
		} else {
			return nil, fmt.Errorf("direction not found")
		}
	}

	now := time.Now()

	log.Printf("Inserted ticket info: %v", data)

	// Directly pass the pointers for all parameters.
	if err := database.InsertTicketInfo(
		&now,
		nil,
		nil,
		linePtr,
		stationNamePtr,
		stationIDPtr,
		directionNamePtr,
		directionIDPtr,
	); err != nil {
		return nil, fmt.Errorf("failed to insert ticket info into database: %v", err)
	}

	return data, nil
}
