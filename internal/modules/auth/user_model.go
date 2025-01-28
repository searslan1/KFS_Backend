package auth

import "time"

// auth temel kullanıcı modeli
type auth struct {
	authID       uint      `gorm:"primaryKey;autoIncrement" json:"auth_id"`       // Primary key
	Email        string    `gorm:"type:varchar(255);unique;not null" json:"email"` // Unique ve not null
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`            // Şifre hashlenmiş olarak saklanır
	authType     string    `gorm:"type:varchar(20);not null" json:"auth_type"`     // Kullanıcı türü
	CreatedAt    time.Time `gorm:"default:current_timestamp" json:"created_at"`    // Hesap oluşturma tarihi
}

type AuthUser struct {
	UserID              uint       `gorm:"primaryKey;not null" json:"user_id"`            // Foreign key, users tablosuna referans
	User                User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"` // Foreign key ile ilişkilendirme
	EmailVerified       bool       `gorm:"default:false" json:"email_verified"`           // Email doğrulama
	PhoneVerified       bool       `gorm:"default:false" json:"phone_verified"`           // Telefon doğrulama
	PasswordResetToken  string     `gorm:"type:varchar(255)" json:"password_reset_token"` // Şifre sıfırlama tokeni
	PasswordResetExpires time.Time `json:"password_reset_expires"`                        // Token geçerlilik süresi
	FailedLoginAttempts int        `gorm:"default:0" json:"failed_login_attempts"`        // Hatalı giriş sayısı
	AccountLockedUntil  time.Time  `json:"account_locked_until"`                          // Hesap kilitli olma süresi
	CreatedAt           time.Time  `gorm:"default:current_timestamp" json:"created_at"`   // Kayıt oluşturma tarihi
}