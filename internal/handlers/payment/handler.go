package payment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"navmate-backend/internal/adapters/payment"
	"navmate-backend/internal/models"
)

type Handler struct{ db *gorm.DB }

func New(db *gorm.DB) *Handler { return &Handler{db: db} }

type authorizeReq struct {
	BookingID   uint `json:"booking_id" binding:"required"`
	AmountCents int  `json:"amount_cents" binding:"required"`
}

// POST /v1/payments/authorize
func (h *Handler) Authorize(c *gin.Context) {
	uid := uint(c.GetInt("user_id"))
	var req authorizeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify booking ownership
	var b models.RideBooking
	if err := h.db.First(&b, req.BookingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

	// Verify plan ownership
	var p models.TripPlan
	if err := h.db.First(&p, b.PlanID).Error; err != nil || p.UserID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// Call payment adapter
	idStr, status := payment.Authorize(req.AmountCents)

	pay := models.Payment{
		AmountCents: req.AmountCents,
		Currency:    "THB",
		Status:      status,
		ExternalRef: idStr,
	}
	if err := h.db.Create(&pay).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create payment failed"})
		return
	}

	// Link payment to booking
	b.PaymentID = &pay.ID
	_ = h.db.Save(&b).Error

	if status == "declined" {
		c.JSON(http.StatusPaymentRequired, gin.H{
			"error":      "payment declined",
			"payment_id": pay.ID,
			"status":     pay.Status,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payment_id":   pay.ID,
		"status":       pay.Status,
		"external_ref": pay.ExternalRef,
	})
}

// POST /v1/payments/:id/capture
func (h *Handler) Capture(c *gin.Context) {
	uid := uint(c.GetInt("user_id"))
	id := c.Param("id")

	var p models.Payment
	if err := h.db.First(&p, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	// Verify ownership via booking
	var b models.RideBooking
	if err := h.db.Where("payment_id = ?", p.ID).First(&b).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

	var plan models.TripPlan
	if err := h.db.First(&plan, b.PlanID).Error; err != nil || plan.UserID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if p.Status != "authorized" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment not authorized"})
		return
	}

	p.Status = "captured"
	if err := h.db.Save(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "capture failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payment_id": p.ID,
		"status":     p.Status,
	})
}

// POST /v1/payments/:id/refund
func (h *Handler) Refund(c *gin.Context) {
	uid := uint(c.GetInt("user_id"))
	id := c.Param("id")

	var p models.Payment
	if err := h.db.First(&p, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	// Verify ownership via booking
	var b models.RideBooking
	if err := h.db.Where("payment_id = ?", p.ID).First(&b).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

	var plan models.TripPlan
	if err := h.db.First(&plan, b.PlanID).Error; err != nil || plan.UserID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if p.Status != "captured" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment not captured"})
		return
	}

	p.Status = "refunded"
	if err := h.db.Save(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "refund failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payment_id": p.ID,
		"status":     p.Status,
	})
}

// POST /v1/payments/webhook (public endpoint)
func (h *Handler) Webhook(c *gin.Context) {
	// TODO: Verify webhook signature
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	// TODO: Process webhook events (payment status updates)
	// For now, just acknowledge receipt
	c.JSON(http.StatusOK, gin.H{"received": true})
}
