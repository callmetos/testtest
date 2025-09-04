package ride

import (
	"math/rand"
	"time"
)

type BookingResult struct {
	Provider   string
	Status     string // confirmed, surge_too_high, failed
	EtaMinutes int
	FareCents  int
}

func Book(provider string, roughFareCents int) BookingResult {
	// Simulate random scenarios
	rand.Seed(time.Now().UnixNano())
	scenario := rand.Intn(10)

	switch {
	case scenario < 7: // 70% success
		return BookingResult{
			Provider:   provider,
			Status:     "confirmed",
			EtaMinutes: rand.Intn(15) + 3,                     // 3-18 minutes
			FareCents:  roughFareCents + rand.Intn(200) - 100, // ±1 THB variance
		}
	case scenario < 9: // 20% surge pricing
		surgeMultiplier := 1.5 + rand.Float64() // 1.5x - 2.5x
		return BookingResult{
			Provider:   provider,
			Status:     "surge_too_high",
			EtaMinutes: rand.Intn(10) + 2,
			FareCents:  int(float64(roughFareCents) * surgeMultiplier),
		}
	default: // 10% failure
		return BookingResult{
			Provider:   provider,
			Status:     "failed",
			EtaMinutes: 0,
			FareCents:  roughFareCents,
		}
	}
}

// Cancel cancels an existing booking
func Cancel(provider, bookingRef string) bool {
	// Simulate cancellation success/failure
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(10) < 9 // 90% success rate
}

// GetStatus gets real-time status of a booking
func GetStatus(provider, bookingRef string) BookingResult {
	statuses := []string{"confirmed", "driver_assigned", "driver_arriving", "in_progress", "completed", "cancelled"}
	rand.Seed(time.Now().UnixNano())

	return BookingResult{
		Provider:   provider,
		Status:     statuses[rand.Intn(len(statuses))],
		EtaMinutes: rand.Intn(20),
		FareCents:  0, // Status check doesn't return fare
	}
}
