package models

import (
	"time"

	"gorm.io/gorm"
)

// Entrepreneur girişimci modeli
type Entrepreneur struct {
	EntrepreneurID   uint      `gorm:"column:entrepreneur_id;primaryKey;autoIncrement"`
	UserID           uint      `gorm:"column:user_id;unique;not null"`
	StartupName      string    `gorm:"column:startup_name;type:varchar(255);not null"`
	Industry         string    `gorm:"column:industry;type:varchar(100);not null"`
	FundingNeeded    float64   `gorm:"column:funding_needed;type:decimal(15,2);not null"`
	BusinessModel    string    `gorm:"column:business_model;type:text"`
	PitchDeckURL     string    `gorm:"column:pitch_deck_url;type:text"`
	Status           string    `gorm:"column:status;type:varchar(20);default:'pending'"`
	CreatedAt        time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime"`
	IsEDevletApproved bool     `gorm:"column:is_edevlet_approved;type:boolean;default:false"`
	IsAdminApproved   bool     `gorm:"column:is_admin_approved;type:boolean;default:false"`
}
