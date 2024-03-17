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
	var StationsAndLinesList types.Data
	StationsAndLinesList, err := ReadStationsAndLinesList("data/StationsAndLinesList.json")

	if err != nil {
		log.Fatalf("Error reading lines: %v", err)
	}
	return c.JSONPretty(http.StatusOK, StationsAndLinesList, "  ")
}

func ReadStationsAndLinesList(filepath string) (types.Data, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return types.Data{}, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return types.Data{}, err
	}

	var data types.Data
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		log.Fatalf("Error reading lines: %v", err)
	}

	return data, nil
}
