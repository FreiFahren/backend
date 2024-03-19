package api

import (
	"log"
	"net/http"

	"encoding/json"
	"io/ioutil"
	"os"

	types "github.com/FreiFahren/backend/structs"
	"github.com/labstack/echo/v4"
)

func GetAllStationsAndLines(c echo.Context) error {
	var StationsAndLinesList types.AllStationsAndLinesList
	StationsAndLinesList, err := ReadStationsAndLinesList("data/StationsAndLinesList.json")

	// only get the lines
	isLineList := c.QueryParam("lines")

	if isLineList == "true" {
		linesList, err := ReadLinesList("data/Lines.json")
		if err != nil {
			log.Fatalf("Error reading lines: %v", err)
		}

		return c.JSONPretty(http.StatusOK, linesList, "  ")
	}

	isStationList := c.QueryParam("stations")
	if isStationList == "true" {
		stationsList, err := ReadStationsList("data/Stations.json")
		if err != nil {
			log.Fatalf("Error reading lines: %v", err)
		}

		return c.JSONPretty(http.StatusOK, stationsList, "  ")
	}
	if err != nil {
		log.Fatalf("Error reading lines: %v", err)
	}

	return c.JSONPretty(http.StatusOK, StationsAndLinesList, "  ")
}

func ReadStationsAndLinesList(filepath string) (types.AllStationsAndLinesList, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return types.AllStationsAndLinesList{}, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return types.AllStationsAndLinesList{}, err
	}

	var AllStationsAndLinesList types.AllStationsAndLinesList
	err = json.Unmarshal(byteValue, &AllStationsAndLinesList)
	if err != nil {
		log.Fatalf("Error reading lines: %v", err)
	}

	return AllStationsAndLinesList, nil
}

func ReadLinesList(filepath string) (map[string][]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return map[string][]string{}, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return map[string][]string{}, err
	}

	var linesList map[string][]string
	err = json.Unmarshal(byteValue, &linesList)
	if err != nil {
		log.Fatalf("Error reading lines: %v", err)
	}

	return linesList, nil
}

func ReadStationsList(filepath string) (map[string]types.StationListEntry, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return map[string]types.StationListEntry{}, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return map[string]types.StationListEntry{}, err
	}

	var linesList map[string]types.StationListEntry
	err = json.Unmarshal(byteValue, &linesList)
	if err != nil {
		log.Fatalf("Error reading lines: %v", err)
	}

	return linesList, nil
}
