package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"KFS_Backend/configs"

	"github.com/golang-jwt/jwt"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	keysLoaded bool // Anahtarların yüklenip yüklenmediğini kontrol etmek için
)

func LoadJWTKeys() error {
	if keysLoaded {
		return nil // Anahtarlar zaten yüklüyse tekrar yükleme
	}

	privateKeyPEM, err := os.ReadFile("configs/jwtRS256.key")
	if err != nil {
		return fmt.Errorf("özel anahtar okuma hatası: %w", err)
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

	if privateKey == nil || publicKey == nil {
		return errors.New("anahtarlar başarıyla yüklenemedi")
	}

	keysLoaded = true
	return nil
}

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// Access Token oluşturma
func GenerateAccessToken(userID uint, role string) (string, error) {
	config := configs.LoadJWTConfig()

	claims := JWTClaims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(config.AccessTokenExp).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

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

// ValidateToken: Gelen access token'ın imzasını ve claims bilgisini doğrular.
func ValidateToken(tokenString string) (*JWTClaims, error) {
	// Eğer RSA anahtarları yüklenmemişse hata döndür
	if !keysLoaded {
		return nil, errors.New("JWT anahtarları yüklenmedi") // Anahtarların yüklenmediğini bildirir.
	}

	// JWT parser'ı oluşturulur ve yalnızca RS256 algoritmasını geçerli kabul eder.
	parser := &jwt.Parser{
		ValidMethods: []string{jwt.SigningMethodRS256.Alg()}, // Sadece RS256 algoritmasını kabul et.
	}

	// tokenString, özel claims yapısına (JWTClaims) göre parse edilip doğrulanır.
	token, err := parser.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// İmza algoritmasının RSA olup olmadığını kontrol et.
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("beklenmeyen imza algoritması: %v", token.Header["alg"])
		}
		return publicKey, nil // Doğrulama için publicKey döndürülür.
	})

	// Parse sırasında hata oluşursa, hatayı wrap ederek döndür.
	if err != nil {
		return nil, fmt.Errorf("token doğrulama hatası: %w", err)
	}

	// Token geçerli değilse hata döndür.
	if !token.Valid {
		return nil, errors.New("geçersiz token")
	}

	// Token'ın claims kısmı JWTClaims tipine dönüştürülmeye çalışılır.
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("token claims parse edilemedi")
	}

	// Doğrulanmış claims döndürülür.
	return claims, nil
}

// ValidateRefreshToken: Gelen refresh token'ın imzasını ve standart claims bilgisini doğrular.
func ValidateRefreshToken(tokenString string) (*jwt.StandardClaims, error) {
	// Eğer RSA anahtarları yüklenmemişse hata döndür.
	if !keysLoaded {
		return nil, errors.New("JWT anahtarları yüklenmedi")
	}

	// JWT parser'ı oluşturulur, geçerli yöntem olarak RS256 kabul edilir.
	parser := &jwt.Parser{
		ValidMethods: []string{jwt.SigningMethodRS256.Alg()},
	}

	// tokenString, standart JWT claims yapısına (jwt.StandardClaims) göre parse edilir ve doğrulanır.
	token, err := parser.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		// İmza algoritmasının RSA olup olmadığını kontrol et.
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("beklenmeyen imza algoritması: %v", token.Header["alg"])
		}
		return publicKey, nil // Doğrulama için publicKey döndürülür.
	})

	// Parse sırasında hata oluşursa, hatayı wrap ederek döndür.
	if err != nil {
		return nil, fmt.Errorf("token doğrulama hatası: %w", err)
	}

	// Token geçerli değilse hata döndür.
	if !token.Valid {
		return nil, errors.New("geçersiz token")
	}

	// Token'ın claims kısmı jwt.StandardClaims tipine dönüştürülmeye çalışılır.
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, errors.New("token claims parse edilemedi")
	}

	// Doğrulanmış standart claims döndürülür.
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
LoadJWTKeys: RSA anahtarlarını dosyalardan yükler ve doğrulama için hazır hale getirir.
JWTClaims: Access token içerisine eklenmek üzere kullanıcıya ait bilgileri ve standart JWT alanlarını tutar.
GenerateAccessToken: Kullanıcı bilgileri ve geçerlilik süresiyle bir access token oluşturur ve privateKey ile imzalar.
GenerateRefreshToken: Oturum bilgilerini içeren bir refresh token oluşturur, daha uzun ömürlüdür.
ValidateToken: Gelen token'ı RS256 algoritmasıyla doğrulayıp, özel claims yapısını parse eder.
ValidateRefreshToken: Refresh token'ı standart JWT claims yapısı üzerinden doğrular.
GenerateSecureRefreshToken: Rastgele güvenli bir refresh token stringi üretir.
*/
