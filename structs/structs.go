package structs

import (
	"database/sql"
	"time"
)

// This is the struct that we will use to store the data from the stations.json file
// getId.go

type Station struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Coordinates Coordinates `json:"coordinates"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// This is the struct that we will use to store the coordinates and the direction id and line
// to return to the frontend
// getData.go

type TicketInspector struct {
	Timestamp time.Time `json:"timestamp"`
	Station   Station   `json:"station"`
	Direction Station   `json:"direction"`
	Line      string    `json:"line"`
}

// For the data received from the database query we will use this struct
// to store the data and return it to the frontend
// database.go

type TicketInfo struct {
	Timestamp    time.Time      `json:"timestamp"`
	Station_ID   string         `json:"station_id"`
	Line         sql.NullString `json:"line"`
	Direction_ID sql.NullString `json:"direction_id"`
}

// PostInspector.go

type InspectorRequest struct {
	Line          string `json:"line"`
	StationName   string `json:"station"`
	DirectionName string `json:"direction"`
}

type ResponseData struct {
	Line      string  `json:"line"`
	Station   Station `json:"station"`
	Direction Station `json:"direction"`
}

// getAllStationsAndLines.go

type StationListEntry struct {
	Name        string           `json:"name"`
	Coordinates CoordinatesEntry `json:"coordinates"`
	Lines       []string         `json:"lines"`
}

type CoordinatesEntry struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Data struct {
	Lines    []map[string][]string       `json:"lines"`
	Stations map[string]StationListEntry `json:"stations"`
}
