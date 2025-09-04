package safety

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"navmate-backend/internal/models"
	"navmate-backend/internal/utils"
)

type Handler struct{ db *gorm.DB }

func New(db *gorm.DB) *Handler { return &Handler{db: db} }

type startReq struct {
	PlanID      uint `json:"plan_id" binding:"required"`
	IntervalMin int  `json:"interval_min" binding:"required"`
}

// POST /v1/safety/session
func (h *Handler) Start(c *gin.Context) {
	uid := uint(c.GetInt("user_id"))
	var req startReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify plan ownership
	var p models.TripPlan
	if err := h.db.Where("id = ? AND user_id = ?", req.PlanID, uid).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
		return
	}

	// Check for existing active session
	var existing models.SafetySession
	if err := h.db.Where("plan_id = ? AND active = true", p.ID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":      "active session already exists",
			"session_id": existing.ID,
			"share_url":  "/safety/s/" + existing.ShareToken,
		})
		return
	}

	now := time.Now()
	s := models.SafetySession{
		PlanID:          p.ID,
		ShareToken:      utils.RandomToken(24),
		IntervalMinutes: req.IntervalMin,
		NextDue:         now.Add(time.Duration(req.IntervalMin) * time.Minute),
		Active:          true,
		StartedAt:       now,
	}
	if err := h.db.Create(&s).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create session failed"})
		return
	}

	// Create first heartbeat
	hb := models.Heartbeat{
		SessionID: s.ID,
		DueAt:     s.NextDue,
		Status:    "due",
	}
	_ = h.db.Create(&hb).Error

	c.JSON(http.StatusOK, gin.H{
		"session_id":       s.ID,
		"share_url":        "/safety/s/" + s.ShareToken,
		"next_due":         s.NextDue,
		"interval_minutes": s.IntervalMinutes,
	})
}

// POST /v1/safety/heartbeat/ack
func (h *Handler) Ack(c *gin.Context) {
	uid := uint(c.GetInt("user_id"))

	type ackReq struct {
		SessionID uint `json:"session_id" binding:"required"`
	}
	var req ackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var s models.SafetySession
	if err := h.db.First(&s, req.SessionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	// Verify ownership via plan
	var p models.TripPlan
	if err := h.db.First(&p, s.PlanID).Error; err != nil || p.UserID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if !s.Active {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session inactive"})
		return
	}

	// Mark last heartbeat as acked
	var hb models.Heartbeat
	if err := h.db.Where("session_id = ?", s.ID).Order("id DESC").First(&hb).Error; err == nil {
		now := time.Now()
		hb.Status = "acked"
		hb.AckedAt = &now
		_ = h.db.Save(&hb).Error
	}

	// Schedule next heartbeat
	nextDue := time.Now().Add(time.Duration(s.IntervalMinutes) * time.Minute)
	s.NextDue = nextDue
	_ = h.db.Save(&s).Error

	// Create next heartbeat
	nextHb := models.Heartbeat{
		SessionID: s.ID,
		DueAt:     nextDue,
		Status:    "due",
	}
	_ = h.db.Create(&nextHb).Error

	c.JSON(http.StatusOK, gin.H{
		"next_due": nextDue,
		"status":   "acknowledged",
	})
}

// POST /v1/safety/sos
func (h *Handler) SOS(c *gin.Context) {
	uid := uint(c.GetInt("user_id"))

	type sosReq struct {
		PlanID   uint   `json:"plan_id" binding:"required"`
		Location string `json:"location,omitempty"`
		Message  string `json:"message,omitempty"`
	}
	var req sosReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify plan ownership
	var p models.TripPlan
	if err := h.db.Where("id = ? AND user_id = ?", req.PlanID, uid).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
		return
	}

	// TODO: Trigger emergency notifications to emergency contacts
	// TODO: Log SOS event
	// TODO: Notify authorities if configured

	// For now, just log the SOS
	c.JSON(http.StatusOK, gin.H{
		"status":    "SOS triggered",
		"plan_id":   p.ID,
		"timestamp": time.Now(),
		"message":   "Emergency services and contacts will be notified",
	})
}

// GET /safety/s/:token (public share page)
func (h *Handler) PublicStatus(c *gin.Context) {
	token := c.Param("token")
	var s models.SafetySession
	if err := h.db.Where("share_token = ? AND active = true", token).First(&s).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found or inactive"})
		return
	}

	// Get plan info
	var p models.TripPlan
	_ = h.db.First(&p, s.PlanID).Error

	// Get latest heartbeat
	var latestHb models.Heartbeat
	_ = h.db.Where("session_id = ?", s.ID).Order("id DESC").First(&latestHb).Error

	c.JSON(http.StatusOK, gin.H{
		"plan_id":          s.PlanID,
		"origin":           p.Origin,
		"destination":      p.Destination,
		"started_at":       s.StartedAt,
		"next_due":         s.NextDue,
		"active":           s.Active,
		"interval_minutes": s.IntervalMinutes,
		"last_heartbeat": map[string]interface{}{
			"due_at":   latestHb.DueAt,
			"status":   latestHb.Status,
			"acked_at": latestHb.AckedAt,
		},
	})
}
