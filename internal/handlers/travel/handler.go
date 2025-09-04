package travel

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"navmate-backend/config" // Import config to get API Key
	"navmate-backend/internal/adapters/maps"
	"navmate-backend/internal/models"
)

type Handler struct {
	db          *gorm.DB
	mapsAdapter *maps.GoogleMapsAdapter // NEW: Use the new adapter
}

// New handler now needs the app config to get the API Key
func New(db *gorm.DB, cfg *config.Config) *Handler {
	// NEW: Initialize the Google Maps Adapter using the key from config
	adapter, err := maps.NewGoogleMapsAdapter(cfg.Google.GoogleMapsAPIKey)
	if err != nil {
		// Log a fatal error if the adapter can't be created, as it's critical.
		log.Fatalf("Failed to create maps adapter: %v", err)
	}

	return &Handler{
		db:          db,
		mapsAdapter: adapter,
	}
}

type planReq struct {
	Origin      string  `json:"origin" binding:"required"`
	Destination string  `json:"destination" binding:"required"`
	DepartAt    *string `json:"depart_at"` // RFC3339 optional
}

func (h *Handler) Plan(c *gin.Context) {
	uid := uint(c.GetInt("user_id"))

	var req planReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var tptr *time.Time
	if req.DepartAt != nil && *req.DepartAt != "" {
		if t, err := time.Parse(time.RFC3339, *req.DepartAt); err == nil {
			tptr = &t
		}
	}

	// NEW: generate options using the real adapter
	opts := h.mapsAdapter.EstimateItineraries(req.Origin, req.Destination, tptr)

	plan := models.TripPlan{UserID: uid, Origin: req.Origin, Destination: req.Destination, DepartAt: tptr}
	if err := h.db.Create(&plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create plan failed"})
		return
	}

	type optResp struct {
		ItineraryID    uint   `json:"itinerary_id"`
		ModeMix        string `json:"mode_mix"`
		TotalMinutes   int    `json:"total_minutes"`
		RoughCostCents int    `json:"rough_cost_cents"`
	}
	resp := make([]optResp, 0, len(opts))

	for _, o := range opts {
		it := models.Itinerary{
			PlanID: plan.ID, ModeMix: o.ModeMix, TotalMinutes: o.TotalMinutes, RoughCostCents: o.RoughCostCents,
		}
		if err := h.db.Create(&it).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "create itinerary failed"})
			return
		}
		for i, l := range o.Legs {
			leg := models.Leg{
				ItineraryID: it.ID, Index: i, Mode: l.Mode, FromName: l.From, ToName: l.To,
				Minutes: l.Minutes, DistanceM: l.DistanceM, Provider: l.Provider,
			}
			if err := h.db.Create(&leg).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "create leg failed"})
				return
			}
		}
		resp = append(resp, optResp{
			ItineraryID: it.ID, ModeMix: it.ModeMix, TotalMinutes: it.TotalMinutes, RoughCostCents: it.RoughCostCents,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"plan_id": plan.ID,
		"options": resp,
	})
}

func (h *Handler) GetPlan(c *gin.Context) {
	uid := uint(c.GetInt("user_id"))
	id := c.Param("id")

	var p models.TripPlan
	if err := h.db.Where("id = ? AND user_id = ?", id, uid).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var itins []models.Itinerary
	_ = h.db.Where("plan_id = ?", p.ID).Order("id ASC").Find(&itins).Error

	c.JSON(http.StatusOK, gin.H{
		"id": p.ID, "origin": p.Origin, "destination": p.Destination, "status": p.Status,
		"selected_itinerary_id": p.SelectedItineraryID, "itinerary_count": len(itins),
	})
}

type selectReq struct {
	ItineraryID uint `json:"itinerary_id" binding:"required"`
}

func (h *Handler) SelectItinerary(c *gin.Context) {
	uid := uint(c.GetInt("user_id"))
	id := c.Param("id")

	var req selectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var p models.TripPlan
	if err := h.db.Where("id = ? AND user_id = ?", id, uid).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
		return
	}

	var cnt int64
	if err := h.db.Model(&models.Itinerary{}).
		Where("id = ? AND plan_id = ?", req.ItineraryID, p.ID).
		Count(&cnt).Error; err != nil || cnt == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid itinerary_id"})
		return
	}

	p.SelectedItineraryID = &req.ItineraryID
	p.Status = "selected"
	if err := h.db.Save(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}
	c.Status(http.StatusNoContent)
}
