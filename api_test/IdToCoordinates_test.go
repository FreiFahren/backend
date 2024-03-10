package api_test

import (
	"testing"
	"github.com/FreiFahren/backend/api"
	"math"
)

func TestIdToCoordinates(t *testing.T) {

	tests := []struct {
		id       string
		expected [2]float64 // Latitude, Longitude
	}{
		{"U-Ado", [2]float64{52.4998948, 13.3071423}},
		{"S-Ad", [2]float64{52.4348328, 13.5414261}},
		{"S-Ba", [2]float64{52.3914002, 13.0928906}},
		{"U-Kt", [2]float64{52.490724, 13.3601535}},
		{"S-OrS", [2]float64{52.5249776, 13.3929084}},
		{"U-Scha", [2]float64{52.5669332, 13.3121764}},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			latitude, longitude, err := api.IdToCoordinates(tt.id)
			if err != nil {
				t.Fatalf("IdToCoordinates(%s) returned an error: %v", tt.id, err)
			}
			if math.Abs(latitude-tt.expected[0]) > 0.000001 || math.Abs(longitude-tt.expected[1]) > 0.000001 {
				t.Errorf("IdToCoordinates(%s) = %v, %v; expected %v, %v", tt.id, latitude, longitude, tt.expected[0], tt.expected[1])
			}
		})
	}
}
