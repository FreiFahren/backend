package api_test

import (
	"log"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/FreiFahren/backend/database"
	"github.com/joho/godotenv"
)

func setup() {
	// Load .env file
	os.Chdir("..")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.CreatePool()

	database.CreatePoolTestTable()
}

func teardown() {
	database.ClosePool()
}

func TestGetLatestStationCoordinatesConcurrency(t *testing.T) {
	setup()
	defer teardown()

	errs := make(chan error, 1000) // Buffer the channel to prevent goroutines from blocking

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			t.Logf("Running test %d", i)

			// Call the function with the pool (on the real database)
			_, err := database.GetLatestStationCoordinates()
			if err != nil {
				errs <- err
			}

			_, err = database.GetHistoricStations(time.Now())
			if err != nil {
				errs <- err
			}

			// Seed the random number generator with the current time
			rand.New(rand.NewSource(20))

			// Enter a random ticket info into the database (on the pool test table)

			now := time.Now()
			message := "Platz der Lufbrücke"
			author := rand.Int63()
			line := "U6"
			stationId := "U-PL"
			stationName := "Platz der Lufbrücke"
			directionName := "Alt-Tegel"
			directionId := "U-ATg"

			err = database.InsertPoolInfo(&now, &message, &author, &line, &stationName, &stationId, &directionName, &directionId)
			if err != nil {
				log.Fatalf("Failed to insert ticket info: %v", err)
			}

		}(i)
	}
	wg.Wait()

	close(errs) // Close the channel to signal that no more errors will be sent

	// Check if any errors were sent on the channel
	if err, ok := <-errs; ok {
		t.Fatalf("Failed to get latest station coordinates: %v", err)
	}
}
