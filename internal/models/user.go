package models

import "time"

type User struct {
	ID           uint    `gorm:"primaryKey"`
	Email        string  `gorm:"uniqueIndex;not null"`
	PasswordHash string  `gorm:"not null"`               // สำหรับ local; ของ Google จะใส่ค่า dummy hash
	Provider     string  `gorm:"default:local;not null"` // local|google
	GoogleID     *string `gorm:"uniqueIndex"`            // อาจเป็น NULL หากเป็น local
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
