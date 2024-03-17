package api

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"
	"io/ioutil"
	"os"

	types "github.com/FreiFahren/backend/structs"
	"github.com/labstack/echo/v4"
)

func GetAllStationsAndLines(c echo.Context) error {
	var StationsAndLinesList = types.StationsAndLinesList{}

	var stationList = make([]types.StationList, 0)

	// Open the files and read lines and stations
	stations, err := ReadFromFile("data/stations.json")
	if err != nil {
		fmt.Println("Error reading stations:", err)
		return err
	}

	lines, err := ReadLines("data/lines.json")
	if err != nil {
		fmt.Println("Error reading lines:", err)
		return err
	}

	// Create the station list with all the names and ids
	for id, station := range stations {
		stationList = append(stationList, types.StationList{StationName: station.Name, StationId: id})
	}

	// Create the lines list
	for line, stations := range lines {
		StationsAndLinesList.Lines = append(StationsAndLinesList.Lines, types.Lines{Name: line, Stations: stations})

		fmt.Println("Line:", line)
		fmt.Println("Stations:", stations)
	}

	// Add the station list to the suggestions
	StationsAndLinesList.StationList = stationList

	// // Sort the lines alphabetically
	// sort.Strings(suggestions.Lines)
	// sort.Slice(suggestions.StationList, func(i, j int) bool {
	// 	return suggestions.StationList[i].StationName < suggestions.StationList[j].StationName
	// })

	// Return the suggestions

	return c.JSONPretty(http.StatusOK, StationsAndLinesList, "  ")
}

func ReadLines(filepath string) (map[string][]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var lines map[string][]string
	err = json.Unmarshal(byteValue, &lines)
	if err != nil {
		log.Fatalf("Error reading lines: %v", err)
	}

	return lines, nil
}
