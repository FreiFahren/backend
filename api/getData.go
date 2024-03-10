package api

import (
	"fmt"
)

func IdToCoordinates(id string) (float64, float64, error) {
	stations, err := ReadFromFile("../data/stations.json")
	if err != nil {
		return 0, 0, err
	}

	station, ok := stations[id]
	if !ok {
		return 0, 0, fmt.Errorf("station ID %s not found", id)
	}

	return station.Coordinates.Latitude, station.Coordinates.Longitude, nil
}
