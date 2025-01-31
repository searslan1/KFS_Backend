package user

import (
	"time"
)

// Kullanıcı Modeli
type User struct {
	UserID       uint      `gorm:"primaryKey"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string    `gorm:"type:varchar(255);not null"`
	Role         string    `gorm:"type:varchar(20);not null;check:role IN ('individual', 'corporate', 'admin', 'moderator')"`
	Verified   	 bool      `gorm:"default:false"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

// Kullanıcı Hesap Güvenliği ve Oturum Yönetimi
type AuthUser struct {
	UserID             uint       `gorm:"not null;uniqueIndex"`
	FailedAttempts     int        `gorm:"default:0"`
	AccountLockedUntil *time.Time `gorm:"default:null"`
	CreatedAt          time.Time  `gorm:"autoCreateTime"`
	UpdatedAt          time.Time  `gorm:"autoUpdateTime"`
	User               User       `gorm:"foreignKey:UserID;references:UserID;constraint:OnDelete:CASCADE"`
}

// E-posta Doğrulama Modeli
type EmailVerification struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"not null;index"`
	Email      string    `gorm:"type:varchar(255);not null"`
	CodeHash   string    `gorm:"type:varchar(255);not null"`
	CodeExpiry time.Time `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	User       User      `gorm:"foreignKey:UserID;references:UserID;constraint:OnDelete:CASCADE"`
}

// Kullanıcı Oturum Yönetimi
type UserSession struct {
	SessionID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID             uint      `gorm:"not null"`
	IPAddress          string    `gorm:"type:inet;not null"`
	UserAgent          string    `gorm:"type:text"`
	DeviceInfo         string    `gorm:"type:text"`
	LoginTime          time.Time `gorm:"autoCreateTime"`
	LogoutTime         *time.Time
	IsActive           bool      `gorm:"default:true"`
	LastActivity       time.Time `gorm:"autoCreateTime;autoUpdateTime"`
	RefreshToken       string    `gorm:"type:varchar(255);uniqueIndex"`
	RefreshTokenExpiry time.Time `gorm:"not null"`
}
