package api

import (
	"time"
    "net/http"
    "github.com/labstack/echo/v4"
    "github.com/FreiFahren/backend/database"
)

type InspectorRequest struct {
    Line        string `json:"line"`
    StationName string `json:"station"`
    Direction   string `json:"direction"`
}


func PostInspector(c echo.Context) error {
    // Decode the request body into an InspectorRequest struct
    var req InspectorRequest
    if err := c.Bind(&req); err != nil {
        return err
    }

    // Initialize the response data structure
    type ResponseData struct {
        Line      string   `json:"line"`
        Station   Station  `json:"station"`
        Direction Station  `json:"direction"`
    }
    var data ResponseData

    // Set the line directly from the request
    data.Line = req.Line

    // Read stations from file
    stations, err := ReadFromFile("data/stations.json")
    if err != nil {
        return err
    }

    // Find station ID for both station and direction
    var stationID, directionID string
    var found bool
    if req.StationName != "" {
        stationID, found = FindStationId(req.StationName, stations)
        if !found {
            return echo.NewHTTPError(http.StatusNotFound, "Station not found")
        }
        data.Station = Station{Name: req.StationName, ID: stationID}
    }

    if req.Direction != "" {
        directionID, found = FindStationId(req.Direction, stations)
        if !found {
            return echo.NewHTTPError(http.StatusNotFound, "Direction not found")
        }
        data.Direction = Station{Name: req.Direction, ID: directionID}
    }

	// get current time
	now := time.Now()

    // Insert the information into the database
    err = database.InsertTicketInfo(
		&now,
		nil, // message is not provided
		nil, // author is not provided
		&data.Line,
		&data.Station.Name,
		&data.Station.ID,
		&data.Direction.Name,
		&data.Direction.ID,
	)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to insert ticket info into database")
    }

    // Return the assembled data as JSON
    return c.JSON(http.StatusOK, data)
}
