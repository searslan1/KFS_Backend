package configs

import (
	"log"
	"os"

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
