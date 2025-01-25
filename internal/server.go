package internal

import (
	"fmt"
	"log"

	"github.com/alfonso/KFS_Backend/configs"
	"github.com/alfonso/KFS_Backend/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Veritabanı bağlantısı (isteğe bağlı)
var DB *gorm.DB

// Sunucuyu başlat
func StartServer() {
	// Config yükle
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Config yüklenirken hata: %v", err)
	}

	// Logger başlat
	logger.InitLogger()
	logger.Info("⚡ API başlatılıyor...")

	// Fiber başlat
	app := fiber.New()

	// Router'ı yükle
	SetupRouter(app)

	// Veritabanına bağlanmayı DENE, ama başarısız olursa API yine de çalışsın
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Database.Host, config.Database.User, config.Database.Password,
		config.Database.Name, config.Database.Port, config.Database.SSLMode,
	)

	var dbErr error
	DB, dbErr = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if dbErr != nil {
		logger.Error("⚠️  Veritabanına bağlanılamadı, ancak API çalışmaya devam ediyor...")
	} else {
		logger.Info("✅ Veritabanına başarıyla bağlandı!")
	}

	// Sunucuyu çalıştır
	port := ":" + config.Server.Port
	logger.Info("🚀 Sunucu " + port + " portunda çalışıyor...")
	log.Fatal(app.Listen(port))
}
