package internal

import (
	"fmt"
	"log"
	//"os"

	"KFS_Backend/configs"
	"KFS_Backend/pkg/logger"
	"KFS_Backend/internal/modules/auth"
	"KFS_Backend/internal/database" // database paketini ekliyoruz

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	//"gorm.io/driver/postgres"
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

	// Veritabanını bağla
	database.ConnectDatabase()

	// Migration işlemini başlat
	database.RunMigrations() // Burada RunMigrations fonksiyonunu çağırıyoruz

	// AuthRepository ve AuthService oluştur
	authRepo := &auth.AuthRepository{DB: database.DB}
	authService := &auth.AuthService{Repo: authRepo}
	authController := &auth.AuthController{Service: authService}

	// Router'ı yükle
	auth.RegisterUserRoutes(app, authController)

	// Sunucuyu çalıştır
	port := ":" + config.Server.Port
	logger.Info("🚀 Sunucu " + port + " portunda çalışıyor...")
	log.Fatal(app.Listen(port))
}
