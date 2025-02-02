package utils

import "time"

type Verification struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     int64     `json:"user_id" gorm:"uniqueIndex"`
	Code       string    `json:"code" gorm:"type:VARCHAR(6)"`
	CodeExpiry time.Time `json:"code_expiry"`
	Type       string    `json:"type" gorm:"type:VARCHAR(10)"` // "email" veya "phone"
	IsVerified bool      `json:"is_verified" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}
