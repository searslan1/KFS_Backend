package internal

import (
	"fmt"
	"log"

	// Paketler
	"KFS_Backend/configs"
	"KFS_Backend/internal/database"
	"KFS_Backend/internal/modules/auth"
	"KFS_Backend/pkg/logger"
	"KFS_Backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Veritabanı bağlantısı
var DB *gorm.DB

func StartServer() {
	// .env dosyasını yükle
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("❌ Error loading .env file")
	}

	// Config yükle
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Config yüklenirken hata: %v", err)
	}

	// Logger başlat
	logger.InitLogger()
	logger.Info("⚡ API başlatılıyor...")

	// ✅ JWT Anahtarlarını Yükle
	if err := utils.LoadJWTKeys(); err != nil {
		log.Fatalf("❌ JWT Anahtarları yüklenemedi: %v", err)
	}
	// Fiber başlat
	app := fiber.New()

	// Middleware ekle
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))
	app.Use(func(c *fiber.Ctx) error {
		logger.Info(fmt.Sprintf("📩 Request: %s %s", c.Method(), c.Path()))
		return c.Next()
	})

	// Veritabanını bağla
	database.ConnectDatabase()

	// Migration işlemini başlat
	database.RunMigrations()

	// AuthRepository ve AuthService oluştur
	authRepo := &auth.AuthRepository{DB: database.DB}
	authService := &auth.AuthService{Repo: authRepo}
	authController := &auth.AuthController{Service: authService}

	// Router'ı yükle
	SetupRouter(app, authController)

	// ✅ Supabase için SSL bağlantısını ayarla
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
		config.Database.Host, config.Database.User, config.Database.Password,
		config.Database.Name, config.Database.Port,
	)

	var dbErr error
	DB, dbErr = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if dbErr != nil {
		logger.Error(fmt.Sprintf("⚠️  Supabase veritabanına bağlanılamadı: %v", dbErr))
	} else {
		logger.Info("✅ Supabase veritabanına başarıyla bağlandı!")
	}

	// Sunucuyu çalıştır
	port := ":" + config.Server.Port
	logger.Info("🚀 Sunucu " + port + " portunda çalışıyor...")
	log.Fatal(app.Listen(port))
}
