package profile

import (
	"time"
)

// UserProfile modeli (User ile ilişkili profil bilgileri)
type UserProfile struct {
	ID                int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            int64     `gorm:"not null;index" json:"user_id"`
	ProfilePictureURL string    `json:"profile_picture_url"`
	Bio               string    `json:"bio"`
	WebsiteURL        string    `json:"website_url"`
	SocialMedia       string    `gorm:"type:JSONB" json:"social_media"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// Address tablosu (Kullanıcının adres bilgileri)
type Address struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"not null;index" json:"user_id"`
	Country   string    `gorm:"not null" json:"country"`
	City      string    `gorm:"not null" json:"city"`
	District  string    `json:"district"`
	Street    string    `json:"street"`
	ZipCode   string    `json:"zip_code"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// MediaFile tablosu (Kullanıcıya ait medya dosyaları)
type MediaFile struct {
	ID         string    `gorm:"primaryKey" json:"id"`
	UserID     int64     `gorm:"not null;index" json:"user_id"`
	FileURL    string    `gorm:"not null" json:"file_url"`
	FileType   string    `gorm:"not null" json:"file_type"`
	FileSize   int       `gorm:"not null" json:"file_size"`
	MimeType   string    `gorm:"not null" json:"mime_type"`
	UploadedAt time.Time `gorm:"autoCreateTime" json:"uploaded_at"`
}

// RoleProfile tablosu (Kullanıcının rol profili)
type RoleProfile struct {
	ID                 int64  `gorm:"primaryKey" json:"id"`
	UserID             int64  `gorm:"not null;index" json:"user_id"`
	CompanyName        string `gorm:"not null" json:"company_name"`
	TaxID              string `gorm:"unique;not null" json:"tax_id"`
	RegistrationNumber string `gorm:"unique" json:"registration_number"`
	TaxOffice          string `gorm:"not null" json:"tax_office"`
	PhoneNumber        string `json:"phone_number"`
}

// CompanyRepresentative tablosu (Kullanıcının şirket temsilciliği)
type CompanyRepresentative struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	CompanyID int64     `gorm:"not null;index" json:"company_id"`
	UserID    int64     `gorm:"not null;index" json:"user_id"`
	Role      string    `gorm:"not null" json:"role"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
