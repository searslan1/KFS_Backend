package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

// Config yapısı
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

// Server yapılandırması
type ServerConfig struct {
	Port string `yaml:"port"`
}

// Database yapılandırması
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"sslmode"`
}

// Config verisini yükler
func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env dosyası yüklenemedi, varsayılan ayarlar kullanılıyor.")
	}

	// YAML dosyasını oku
	file, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	// Çevresel değişkenleri oku
	if port := os.Getenv("PORT"); port != "" {
		config.Server.Port = port
	}

	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		config.Database.User = dbUser
	}

	return &config, nil
}
