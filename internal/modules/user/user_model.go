package user

import (
	"time"
)

type User struct {
	UserID       int64     `json:"user_id" gorm:"primaryKey"` // Primary Key ve Auto Increment
	Email        string    `json:"email" gorm:"unique;not null"`                 // Unique ve Not Null
	PasswordHash string    `json:"password_hash" gorm:"not null"`                // Not Null
	UserType     string    `json:"user_type" gorm:"type:VARCHAR(20);not null"`   // UserType VARCHAR(20)
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`             // Otomatik oluşturulma zamanı
}

// AuthUser tablosu
type AuthUser struct {
	UserID               int64     `json:"user_id" gorm:"primaryKey"`    // Foreign Key ve Not Null
	EmailVerified        bool      `json:"email_verified" gorm:"default:false"`   // Default False
	PhoneVerified        bool      `json:"phone_verified" gorm:"default:false"`   // Default False
	PasswordResetToken   string    `json:"password_reset_token,omitempty"`        // Opsiyonel
	PasswordResetExpires time.Time `json:"password_reset_expires,omitempty"`      // Opsiyonel
	FailedLoginAttempts  int       `json:"failed_login_attempts" gorm:"default:0"`// Default 0
	AccountLockedUntil   time.Time `json:"account_locked_until,omitempty"`        // Opsiyonel
	CreatedAt            time.Time `json:"created_at" gorm:"autoCreateTime"`      // Otomatik oluşturulma zamanı
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

type BlacklistedToken struct {
    ID          uint      `gorm:"primaryKey"`
    Token       string    `gorm:"unique;not null"`
    Expiry      time.Time `gorm:"not null"`
    CreatedAt   time.Time
}