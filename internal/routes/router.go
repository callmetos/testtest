package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"navmate-backend/config"
	"navmate-backend/internal/handlers/auth"
	"navmate-backend/internal/handlers/booking"
	"navmate-backend/internal/handlers/payment"
	"navmate-backend/internal/handlers/safety"
	"navmate-backend/internal/handlers/travel"
	"navmate-backend/internal/middleware"
	"navmate-backend/pkg/jwtauth"
)

func SetupRouter(router *gin.Engine, DB *gorm.DB, cfg *config.Config) {
	router.Use(gin.Recovery())

	// Health checks
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "timestamp": time.Now().UTC().Format(time.RFC3339)})
	})
	router.GET("/api/health", func(c *gin.Context) { // alias
		c.JSON(http.StatusOK, gin.H{"status": "ok", "timestamp": time.Now().UTC().Format(time.RFC3339)})
	})

	jwtSvc := jwtauth.NewFromEnv()

	if cfg != nil {
		gh := auth.NewGoogleHandler(DB, jwtSvc, cfg)
		router.GET("/auth/google/login", gh.Login)
		router.GET("/auth/google/callback", gh.Callback)
	}

	safeH := safety.New(DB)
	router.GET("/safety/s/:token", safeH.PublicStatus)

	// API v1 routes
	v1 := router.Group("/v1")
	{
		// Auth routes (BE-2)
		a := auth.New(DB, jwtSvc)
		v1.POST("/auth/signup", a.Register)
		v1.POST("/auth/login", a.Login)
		v1.GET("/me", middleware.AuthJWT(jwtSvc), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"user_id": c.GetInt("user_id"), "email": c.GetString("email")})
		})

		// Trip planning routes (BE-5)
		travH := travel.New(DB)
		v1.POST("/trips/plan", middleware.AuthJWT(jwtSvc), travH.Plan)
		v1.GET("/trips/plans/:id", middleware.AuthJWT(jwtSvc), travH.GetPlan)
		v1.POST("/trips/plans/:id/select", middleware.AuthJWT(jwtSvc), travH.SelectItinerary)

		// Booking routes (BE-6)
		bookH := booking.New(DB)
		v1.POST("/bookings", middleware.AuthJWT(jwtSvc), bookH.Create)
		v1.GET("/bookings/:id", middleware.AuthJWT(jwtSvc), bookH.Get)
		//v1.DELETE("/bookings/:id", middleware.AuthJWT(jwtSvc), bookH.Cancel)

		// Payment routes (BE-7)
		payH := payment.New(DB)
		v1.POST("/payments/authorize", middleware.AuthJWT(jwtSvc), payH.Authorize)
		v1.POST("/payments/:id/capture", middleware.AuthJWT(jwtSvc), payH.Capture)
		v1.POST("/payments/:id/refund", middleware.AuthJWT(jwtSvc), payH.Refund)
		v1.POST("/payments/webhook", payH.Webhook)

		// Safety routes (BE-9)
		v1.POST("/safety/session", middleware.AuthJWT(jwtSvc), safeH.Start)
		v1.POST("/safety/heartbeat/ack", middleware.AuthJWT(jwtSvc), safeH.Ack)
		v1.POST("/safety/sos", middleware.AuthJWT(jwtSvc), safeH.SOS)
	}
}
