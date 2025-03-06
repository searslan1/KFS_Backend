package internal

import (
	"fmt"
	"log"

	"KFS_Backend/configs"
	"KFS_Backend/internal/modules/campaign"
	"KFS_Backend/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Veritabanı bağlantısı
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

	// Veritabanı bağlantısı için DSN oluştur
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Database.Host,
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
		config.Database.Port,
		config.Database.SSLMode,
	)

	var dbErr error
	DB, dbErr = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if dbErr != nil {
		logger.Error(fmt.Sprintf("⚠️  Veritabanına bağlanılamadı: %v", dbErr))
	} else {
		logger.Info("✅ Veritabanına başarıyla bağlanıldı!")
	}

	// Sağlık kontrol endpoint'i
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Kampanya modülü bileşenlerini oluştur
	campaignRepo := campaign.NewRepository(DB)
	campaignService := campaign.NewService(campaignRepo)
	campaignController := campaign.NewController(campaignService)

	// Kampanya rotalarını kaydet
	campaign.RegisterRoutes(app, campaignController)

	// Sunucuyu çalıştır
	port := ":" + config.Server.Port
	logger.Info("🚀 Sunucu " + port + " portunda çalışıyor...")
	log.Fatal(app.Listen(port))
}
