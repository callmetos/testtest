package maps

import "time"

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

func EstimateItineraries(origin, destination string, departAt *time.Time) []ItinOpt {
	_ = departAt
	ride := "RideNow"
	return []ItinOpt{
		{ // 1) Transit only
			ModeMix:      "WALK+TRANSIT",
			TotalMinutes: 42, RoughCostCents: 3000,
			Legs: []LegOpt{
				{Mode: "WALK", From: origin, To: "Station A", Minutes: 8, DistanceM: 600},
				{Mode: "TRANSIT", From: "Station A", To: "Station B", Minutes: 30, DistanceM: 12000},
				{Mode: "WALK", From: "Station B", To: destination, Minutes: 4, DistanceM: 300},
			},
		},
		{ // 2) Ride only
			ModeMix:      "RIDE",
			TotalMinutes: 18, RoughCostCents: 12000,
			Legs: []LegOpt{
				{Mode: "RIDE", From: origin, To: destination, Minutes: 18, DistanceM: 9000, Provider: &ride},
			},
		},
		{ // 3) Mixed
			ModeMix:      "WALK+TRANSIT+RIDE",
			TotalMinutes: 28, RoughCostCents: 9000,
			Legs: []LegOpt{
				{Mode: "WALK", From: origin, To: "Stop C", Minutes: 5, DistanceM: 400},
				{Mode: "TRANSIT", From: "Stop C", To: "Hub D", Minutes: 15, DistanceM: 7000},
				{Mode: "RIDE", From: "Hub D", To: destination, Minutes: 8, DistanceM: 3000, Provider: &ride},
			},
		},
	}
}
