package booking

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"navmate-backend/internal/adapters/ride"
	"navmate-backend/internal/models"
)

type Handler struct{ db *gorm.DB }

func New(db *gorm.DB) *Handler { return &Handler{db: db} }

type createReq struct {
	PlanID uint `json:"plan_id" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	uid := uint(c.GetInt("user_id"))
	var req createReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// load plan (owner only) + selected itinerary
	var p models.TripPlan
	if err := h.db.Where("id = ? AND user_id = ?", req.PlanID, uid).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
		return
	}
	if p.SelectedItineraryID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plan not selected"})
		return
	}
	var it models.Itinerary
	if err := h.db.Where("id = ?", *p.SelectedItineraryID).First(&it).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "itinerary missing"})
		return
	}

	// ถ้ามี RIDE leg ให้จอง, ถ้าไม่มีถือว่า transit-only (ไม่ต้อง book)
	var rideProvider *string
	var rideFound bool
	var fare = it.RoughCostCents
	var legs []models.Leg
	_ = h.db.Where("itinerary_id = ?", it.ID).Order("index ASC").Find(&legs).Error
	for _, l := range legs {
		if l.Mode == "RIDE" && l.Provider != nil {
			rideProvider = l.Provider
			rideFound = true
			break
		}
	}
	if !rideFound {
		c.JSON(http.StatusOK, gin.H{"message": "no ride legs; nothing to book", "plan_id": p.ID, "itinerary_id": it.ID})
		return
	}

	br := ride.Book(*rideProvider, fare)

	b := models.RideBooking{
		PlanID: p.ID, ItineraryID: it.ID, Provider: br.Provider,
		Status: br.Status, EtaMinutes: br.EtaMinutes, FareCents: br.FareCents,
	}
	if err := h.db.Create(&b).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create booking failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"booking_id": b.ID, "status": b.Status, "eta_minutes": b.EtaMinutes, "fare_cents": b.FareCents,
	})
}

func (h *Handler) Get(c *gin.Context) {
	uid := uint(c.GetInt("user_id"))
	id := c.Param("id")
	var b models.RideBooking

	// Verify ownership and fetch booking in a single query
	err := h.db.Joins("JOIN trip_plans ON trip_plans.id = ride_bookings.plan_id").
		Where("ride_bookings.id = ? AND trip_plans.user_id = ?", id, uid).
		First(&b).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, b)
}
