package auth

import "time"

// User tablosu
type User struct {
	UserID       int64     `json:"user_id" gorm:"primaryKey"` // Primary Key ve Auto Increment
	Email        string    `json:"email" gorm:"unique;not null"`                 // Unique ve Not Null
	PasswordHash string    `json:"password_hash" gorm:"not null"`                // Not Null
	UserType     string    `json:"user_type" gorm:"type:VARCHAR(20);not null"`   // UserType VARCHAR(20)
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`             // Otomatik oluşturulma zamanı
	//OTP'nin geçerlilik süresini utmak için bunu kullandım
	OTP          string    `json:"otp" gorm:"type:VARCHAR(6)"`
	OTPExpiresAt time.Time `json:"otp_expires_at"`
	PhoneVerified bool     `json:"phone_verified" gorm:"default:false"`          // Telefon doğrulama
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
//structlar ile çalıştığımız için sql yerine gorm ile çalıştım