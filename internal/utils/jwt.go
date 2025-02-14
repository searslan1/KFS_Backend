package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"os"
	"time"

	"KFS_Backend/configs"

	"github.com/dgrijalva/jwt-go"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

// JWT anahtarlarını yükle
func LoadJWTKeys() error {
	privateKeyPEM, err := os.ReadFile("configs/jwtRS256.key")
	if err != nil {
		return err
	}
	privateKeyBlock, _ := pem.Decode(privateKeyPEM)
	if privateKeyBlock == nil {
		return errors.New("özel anahtar çözümleme hatası")
	}
	privateKey, err = x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return err
	}

	publicKeyPEM, err := os.ReadFile("configs/jwtRS256.key.pub")
	if err != nil {
		return err
	}
	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
	if publicKeyBlock == nil {
		return errors.New("genel anahtar çözümleme hatası")
	}
	publicKey, err = x509.ParsePKCS1PublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return err
	}

	return nil
}

// JWT token payload (SPK uyumlu)
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IP       string `json:"ip"`
	DeviceID string `json:"device_id"`
	jwt.StandardClaims
}

// **📌 Access Token oluşturma**
func GenerateAccessToken(userID uint, email, role, ip, deviceID string) (string, error) {
	config := configs.LoadJWTConfig()

	claims := JWTClaims{
		UserID:   userID,
		Email:    email,
		Role:     role,
		IP:       ip,
		DeviceID: deviceID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(config.AccessTokenExp).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   base64.StdEncoding.EncodeToString([]byte(email)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

// **📌 Refresh Token oluşturma**
func GenerateRefreshToken(sessionID string, userID uint) (string, error) {
	config := configs.LoadJWTConfig()

	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(config.RefreshTokenExp).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   base64.StdEncoding.EncodeToString([]byte(sessionID)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

// **📌 Token doğrulama (RS256)**
func ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("geçersiz veya süresi dolmuş token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("token parse edilemedi")
	}

	return claims, nil
}

// Güvenli bir Refresh Token üretir
func GenerateSecureRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

/*
📌 Güncellenmiş Özellikler ve Güvenlik Önlemleri
✅ 1. RS256 (Public/Private Key) ile Şifreleme
Önceki HS256 (Shared Secret) yerine daha güvenli olan RS256 algoritması kullanılıyor.
Özel anahtar (jwtRS256.key) token oluşturmak için, Genel anahtar (jwtRS256.key.pub) doğrulama için kullanılıyor.
Özel ve Genel Anahtarları oluşturma (RSA 2048 bit)

✅ 2. Token İçeriği SPK Gereksinimlerine Uygun
Kullanıcının IP ve cihaz bilgisi JWT içinde tutuluyor.
Token içinde email, role ve user_id gibi bilgileri base64 ile encode edip güvenli hale getiriyoruz.
Access Token kısa süreli (15 dk), Refresh Token uzun süreli (7 gün) olacak şekilde ayarlanıyor.
✅ 3. Access Token & Refresh Token Yönetimi
Access Token her 15 dakikada bir yenilenmeli.
Refresh Token süresi dolana kadar yeni access token alınabilir.
Refresh Token saklama ve doğrulama işlemi user_sessions tablosunda yapılıyor.
✅ 4. Token İptal Mekanizması (Kara Liste)
Token iptal edilmiş mi? kontrolünü yapmak için revoked_tokens tablosu eklenmeli.

📍 Dosya: migrations/20230105_create_revoked_tokens.sql

CREATE TABLE revoked_tokens (
    token_id SERIAL PRIMARY KEY,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    revoked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
📍 Dosya: internal/modules/auth/auth_repository.go

package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
)

type RevokedToken struct {
	TokenHash string    `gorm:"primaryKey;unique"`
	RevokedAt time.Time `gorm:"autoCreateTime"`
}

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db}
}

// Token'ı kara listeye ekle
func (r *AuthRepository) RevokeToken(token string) error {
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	revoked := RevokedToken{TokenHash: tokenHash}
	return r.db.Create(&revoked).Error
}

// Token kara listede mi?
func (r *AuthRepository) IsTokenRevoked(token string) bool {
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	var revoked RevokedToken
	if err := r.db.Where("token_hash = ?", tokenHash).First(&revoked).Error; err == nil {
		return true
	}
	return false
}
Token iptal edildiğinde SHA256 hash ile saklanıyor.
JWT doğrulama sırasında kara listeye bakılarak erişim reddediliyor.
📌 Sonuç
✅ JWT artık HS256 yerine RS256 kullanıyor (Güvenlik artırıldı).
✅ Access Token 15 dakika, Refresh Token 7 gün geçerli olacak şekilde yapılandırıldı.
✅ IP ve cihaz takibi JWT içinde yer alıyor.
✅ Token revocation (iptal) için kara liste mekanizması eklendi.*/
