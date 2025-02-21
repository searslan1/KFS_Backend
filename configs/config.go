package configs

import (
	"log"
	"os"
	"time"
	"github.com/joho/godotenv"
)

// Config yapısı
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// Server yapılandırması
type ServerConfig struct {
	Port string
}

// Database yapılandırması
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// Config verisini yükler
func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env dosyası yüklenemedi, varsayılan ayarlar kullanılıyor.")
	}

	config := &Config{
		Server: ServerConfig{
			Port: os.Getenv("PORT"),
		},
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
	}

	return config, nil
}
/*config: Uygulamanın genel yapılandırmasını ve çevresel parametrelerini tutar. 
Bu dosya genellikle global ayarları içerir, yani API anahtarları, 
veritabanı bağlantı bilgileri ve uygulamanın genel çalışma koşullarını burada belirleyebilirsiniz.*/


type JWTConfig struct {
    SecretKey       string
    AccessTokenExp  time.Duration
    RefreshTokenExp time.Duration
}

func LoadJWTConfig() JWTConfig {
    return JWTConfig{
        SecretKey:       os.Getenv("JWT_SECRET"),
        AccessTokenExp:  15 * time.Minute,
        RefreshTokenExp: 7 * 24 * time.Hour, // 7 gün
    }
}
