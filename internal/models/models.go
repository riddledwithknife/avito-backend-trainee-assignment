package models

import (
	"time"

	"github.com/lib/pq"
)

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbName"`
	} `json:"database"`
	Redis struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	} `json:"redis"`
}

type Banner struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	FeatureID uint          `gorm:"not null" json:"feature_id"`
	TagIDs    pq.Int64Array `gorm:"type:integer[]" json:"tag_ids"`
	Title     string        `gorm:"not null" json:"title"`
	Text      string        `gorm:"not null" json:"text"`
	URL       string        `gorm:"not null" json:"url"`
	IsActive  bool          `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type Token struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Token   string `gorm:"not null;unique" json:"token"`
	IsAdmin bool   `gorm:"not null" json:"is_admin"`
}
