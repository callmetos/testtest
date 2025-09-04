package jwtauth

import (
	"os"
	"time"
)

func getSecret() string {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		s = "dev-secret-change-me" // ใส่ค่า production เองใน .env
	}
	return s
}

func getTTL() time.Duration {
	ttlStr := os.Getenv("JWT_TTL")
	if ttlStr == "" {
		ttlStr = "168h" // 7 วัน
	}
	d, err := time.ParseDuration(ttlStr)
	if err != nil {
		return 168 * time.Hour
	}
	return d
}
