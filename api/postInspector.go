package api

import (
    "encoding/json"
    "net/http"

    "github.com/labstack/echo/v4"
)

// Define the structures for the request and station
type InspectorRequest struct {
    Line        string `json:"line"`
    StationName string `json:"station"`
    Direction   string `json:"direction"`
}

// Your PostInspector function
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

    // Find station ID
    if req.StationName != "" {
        id, found := FindStationId(req.StationName, stations)
        if found {
            data.Station = Station{Name: req.StationName, ID: id}
        } else {
            return echo.NewHTTPError(http.StatusNotFound, "Station not found")
        }
    }

    // Find direction ID
    if req.Direction != "" {
        id, found := FindStationId(req.Direction, stations)
        if found {
            data.Direction = Station{Name: req.Direction, ID: id}
        } else {
            return echo.NewHTTPError(http.StatusNotFound, "Direction not found")
        }
    }

    // For demonstration, we'll print the assembled data
    // In your real application, you'd likely save this data to a database or another storage
    jsonData, err := json.Marshal(data)
    if err != nil {
        return err
    }
    println(string(jsonData))

    // Return the assembled data as JSON
    return c.JSON(http.StatusOK, data)
}
