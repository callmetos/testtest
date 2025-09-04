package maps

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"googlemaps.github.io/maps"
)

// ItinOpt and LegOpt structs remain the same as they define the output format.
type LegOpt struct {
	Mode      string
	From, To  string
	Minutes   int
	DistanceM int64
	Provider  *string
}
type ItinOpt struct {
	ModeMix        string
	TotalMinutes   int
	RoughCostCents int
	Legs           []LegOpt
}

// GoogleMapsAdapter handles communication with Google Maps APIs
type GoogleMapsAdapter struct {
	client *maps.Client
}

// NewGoogleMapsAdapter creates a new adapter instance
func NewGoogleMapsAdapter(apiKey string) (*GoogleMapsAdapter, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("Google Maps API key is missing")
	}
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create google maps client: %w", err)
	}
	return &GoogleMapsAdapter{client: client}, nil
}

// EstimateItineraries calculates route options using Google Directions API
func (a *GoogleMapsAdapter) EstimateItineraries(origin, destination string, departAt *time.Time) []ItinOpt {
	var options []ItinOpt

	// 1. Get directions for DRIVING (to simulate RIDE)
	drivingReq := &maps.DirectionsRequest{
		Origin:      origin,
		Destination: destination,
		Mode:        maps.TravelModeDriving,
	}
	drivingRoute, _, err := a.client.Directions(context.Background(), drivingReq)
	if err != nil {
		log.Printf("Warning: Error getting driving directions: %v", err)
	}

	// 2. Get directions for TRANSIT
	transitReq := &maps.DirectionsRequest{
		Origin:      origin,
		Destination: destination,
		Mode:        maps.TravelModeTransit,
	}
	transitRoute, _, err := a.client.Directions(context.Background(), transitReq)
	if err != nil {
		log.Printf("Warning: Error getting transit directions: %v", err)
	}

	// Process DRIVING route to create a RIDE option
	if len(drivingRoute) > 0 && len(drivingRoute[0].Legs) > 0 {
		leg := drivingRoute[0].Legs[0]
		rideProvider := "RideNow"
		options = append(options, ItinOpt{
			ModeMix:        "RIDE",
			TotalMinutes:   int(math.Round(leg.Duration.Minutes())),
			RoughCostCents: calculateRideFare(leg.Distance.Meters, int(leg.Duration.Seconds())),
			Legs: []LegOpt{{
				Mode:      "RIDE",
				From:      leg.StartAddress,
				To:        leg.EndAddress,
				Minutes:   int(math.Round(leg.Duration.Minutes())),
				DistanceM: int64(leg.Distance.Meters),
				Provider:  &rideProvider,
			}},
		})
	}

	// Process TRANSIT route to create a WALK+TRANSIT option
	if len(transitRoute) > 0 && len(transitRoute[0].Legs) > 0 {
		leg := transitRoute[0].Legs[0]
		options = append(options, ItinOpt{
			ModeMix:        "WALK+TRANSIT",
			TotalMinutes:   int(math.Round(leg.Duration.Minutes())),
			RoughCostCents: 3500, // Assume a flat fee for transit for simplicity
			Legs: []LegOpt{{
				Mode:      "TRANSIT", // Simplified to one leg for now
				From:      leg.StartAddress,
				To:        leg.EndAddress,
				Minutes:   int(math.Round(leg.Duration.Minutes())),
				DistanceM: int64(leg.Distance.Meters),
			}},
		})
	}

	// If no options were found from Google, return the original stub data as a fallback
	if len(options) == 0 {
		log.Println("No routes found from Google, falling back to stub data.")
		return getStubData(origin, destination)
	}

	return options
}

// calculateRideFare is a simple fare calculation model.
// Example: 40 THB base fare + 8 THB/km + 2 THB/min
func calculateRideFare(distanceMeters int, durationSeconds int) int {
	baseFareCents := 4000
	perKmCents := 800
	perMinCents := 200
	km := float64(distanceMeters) / 1000.0
	minutes := float64(durationSeconds) / 60.0
	cost := float64(baseFareCents) + (km * float64(perKmCents)) + (minutes * float64(perMinCents))
	return int(math.Round(cost))
}

// getStubData provides fallback data if Google API fails
func getStubData(origin, destination string) []ItinOpt {
	ride := "RideNow"
	return []ItinOpt{
		{ModeMix: "RIDE", TotalMinutes: 18, RoughCostCents: 12000, Legs: []LegOpt{{Mode: "RIDE", From: origin, To: destination, Minutes: 18, DistanceM: 9000, Provider: &ride}}},
		{ModeMix: "WALK+TRANSIT", TotalMinutes: 42, RoughCostCents: 3000, Legs: []LegOpt{{Mode: "WALK", From: origin, To: "Station A", Minutes: 8, DistanceM: 600}, {Mode: "TRANSIT", From: "Station A", To: "Station B", Minutes: 30, DistanceM: 12000}, {Mode: "WALK", From: "Station B", To: destination, Minutes: 4, DistanceM: 300}}},
	}
}
