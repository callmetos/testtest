package models

import "time"

// TripPlan = คำขอวางแผนการเดินทางของผู้ใช้ (ประตู-ถึง-ประตู)
type TripPlan struct {
	ID                  uint       `gorm:"primaryKey" json:"id"`
	UserID              uint       `gorm:"index;not null" json:"user_id"`
	Origin              string     `gorm:"not null" json:"origin"`
	Destination         string     `gorm:"not null" json:"destination"`
	DepartAt            *time.Time `json:"depart_at,omitempty"`
	Status              string     `gorm:"not null;default:planned" json:"status"` // planned|selected|active|completed|cancelled
	SelectedItineraryID *uint      `json:"selected_itinerary_id,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`

	Itineraries []Itinerary `gorm:"foreignKey:PlanID;constraint:OnDelete:CASCADE;" json:"-"`
}

type Itinerary struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	PlanID         uint      `gorm:"index;not null" json:"plan_id"`
	ModeMix        string    `gorm:"not null" json:"mode_mix"` // เช่น WALK+TRANSIT+RIDE
	TotalMinutes   int       `gorm:"not null" json:"total_minutes"`
	RoughCostCents int       `gorm:"not null" json:"rough_cost_cents"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Legs []Leg `gorm:"foreignKey:ItineraryID;constraint:OnDelete:CASCADE;" json:"-"`
}

type Leg struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	ItineraryID uint    `gorm:"index;not null" json:"itinerary_id"`
	Index       int     `gorm:"not null" json:"index"`
	Mode        string  `gorm:"not null" json:"mode"` // WALK|TRANSIT|RIDE
	FromName    string  `gorm:"not null" json:"from_name"`
	ToName      string  `gorm:"not null" json:"to_name"`
	Minutes     int     `gorm:"not null" json:"minutes"`
	DistanceM   int64   `gorm:"not null" json:"distance_m"`
	Provider    *string `json:"provider,omitempty"` // สำหรับ RIDE
}

// การจองรถ (สำหรับ RIDE legs)
type RideBooking struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PlanID      uint      `gorm:"index;not null" json:"plan_id"`
	ItineraryID uint      `gorm:"index;not null" json:"itinerary_id"`
	Provider    string    `gorm:"not null" json:"provider"`
	Status      string    `gorm:"not null;default:pending" json:"status"` // pending|confirmed|cancelled|driver_assigned|in_progress|completed
	EtaMinutes  int       `gorm:"not null" json:"eta_minutes"`
	FareCents   int       `gorm:"not null" json:"fare_cents"`
	PaymentID   *uint     `json:"payment_id,omitempty"`
	QuoteID     string    `gorm:"index" json:"quote_id,omitempty"` // for idempotency
	ExternalRef string    `json:"external_ref,omitempty"`          // provider's booking reference
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Payment struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	AmountCents int       `gorm:"not null" json:"amount_cents"`
	Currency    string    `gorm:"not null;default:THB" json:"currency"`
	Status      string    `gorm:"not null;default:authorized" json:"status"` // authorized|captured|refunded|voided|declined
	ExternalRef string    `gorm:"not null;default:'stub'" json:"external_ref"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SafetySession struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	PlanID          uint       `gorm:"index;not null" json:"plan_id"`
	ShareToken      string     `gorm:"uniqueIndex;not null" json:"share_token"`
	IntervalMinutes int        `gorm:"not null" json:"interval_minutes"`
	NextDue         time.Time  `gorm:"not null" json:"next_due"`
	Active          bool       `gorm:"not null;default:true" json:"active"`
	StartedAt       time.Time  `gorm:"not null" json:"started_at"`
	EndedAt         *time.Time `json:"ended_at,omitempty"`
}

type Heartbeat struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	SessionID uint       `gorm:"index;not null" json:"session_id"`
	DueAt     time.Time  `gorm:"not null" json:"due_at"`
	AckedAt   *time.Time `json:"acked_at,omitempty"`
	Status    string     `gorm:"not null;default:due" json:"status"` // due|acked|escalated
	CreatedAt time.Time  `json:"created_at"`
}

// EmergencyContact for future use
type EmergencyContact struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Name      string    `gorm:"not null" json:"name"`
	Phone     string    `gorm:"not null" json:"phone"`
	Email     string    `json:"email,omitempty"`
	Priority  int       `gorm:"not null;default:1" json:"priority"` // 1 = highest
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
