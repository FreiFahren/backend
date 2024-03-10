package api

import (
	"fmt"
	"net/http"
	"time"
	"github.com/labstack/echo/v4"
	"github.com/FreiFahren/backend/database"
)

// Moved outside to make it accessible throughout the package.
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

    data := &ResponseData{Line: req.Line}

    if stationID, found := FindStationId(req.StationName, stations); found {
        data.Station = Station{Name: req.StationName, ID: stationID}
    } else if req.StationName != "" {
        return nil, fmt.Errorf("Station not found")
    }

    if directionID, found := FindStationId(req.Direction, stations); found {
        data.Direction = Station{Name: req.Direction, ID: directionID}
    } else if req.Direction != "" {
        return nil, fmt.Errorf("Direction not found")
    }

    // Insert into database
    now := time.Now()
    if err := database.InsertTicketInfo(
        &now,
		nil, // no message was provided
		nil, // no author was provided
        &data.Line,
        &data.Station.Name, &data.Station.ID,
        &data.Direction.Name, &data.Direction.ID,
    ); err != nil {
        return nil, fmt.Errorf("Failed to insert ticket info into database: %v", err)
    }

    return data, nil
}
