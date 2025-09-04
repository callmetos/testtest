package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"navmate-backend/config"
	"navmate-backend/db"
	"navmate-backend/internal/routes"
)

func main() {
	// Logger
	// เปลี่ยนจาก zap.NewProduction() เป็น zap.NewDevelopment() เพื่อให้ Log สวยงามขึ้น
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()

	// Config
	cfg := config.Load()
	sugar.Infow("config loaded", "port", cfg.Server.Port)

	// DB
	_, err := db.InitDB(cfg)
	if err != nil {
		sugar.Fatalf("db init: %v", err)
	}
	sugar.Infow("db initialized")

	// Router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Basic CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-control-allow-methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Routes
	routes.SetupRouter(router, db.DB, cfg)

	// Start
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	sugar.Infow("server starting", "addr", addr)
	if err := router.Run(addr); err != nil {
		sugar.Fatalf("server run: %v", err)
	}
}
