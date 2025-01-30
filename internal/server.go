package internal

import (
	"fmt"
	"log"
	"os"

	"KFS_Backend/configs"
	"KFS_Backend/pkg/logger"
	"KFS_Backend/internal/modules/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
)

// Veritabanı bağlantısı
var DB *gorm.DB

func StartServer() {
	// .env dosyasını yükle
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

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

	// Middleware ekle
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))
	app.Use(func(c *fiber.Ctx) error {
		logger.Info(fmt.Sprintf("Request: %s %s", c.Method(), c.Path()))
		return c.Next()
	})

	// Veritabanı bağlantı bilgilerini oku
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"), os.Getenv("DB_SSLMODE"),
	)

	var dbErr error
	DB, dbErr = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if dbErr != nil {
		logger.Error(fmt.Sprintf("⚠️  Supabase veritabanına bağlanılamadı: %v", dbErr))
		log.Fatal("Uygulama durduruluyor...")
	}

	// Migration
	MigrateDatabase()

	// AuthRepository ve AuthService oluştur
	authRepo := &auth.AuthRepository{DB: DB}
	authService := &auth.AuthService{Repo: authRepo}
	authController := &auth.AuthController{Service: authService}

	// Router'ı yükle
	auth.RegisterUserRoutes(app, authController)

	// Sunucuyu çalıştır
	port := ":" + config.Server.Port
	logger.Info("🚀 Sunucu " + port + " portunda çalışıyor...")
	log.Fatal(app.Listen(port))
}

// Tabloları migrate etmek için
func MigrateDatabase() {
	logger.Info("⚡ Veritabanı migrasyonu başlatılıyor...")
	if err := DB.AutoMigrate(&auth.User{}, &auth.AuthUser{}); err != nil {
		log.Fatalf("⚠️  Migration sırasında hata: %v", err)
	} else {
		logger.Info("✅ Veritabanı migrasyonu tamamlandı!")
	}
}

