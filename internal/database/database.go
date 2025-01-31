package database

import (
	"fmt"
	"log"
	"time"

	"KFS_Backend/configs"
	"KFS_Backend/internal/modules/campaign"
	"KFS_Backend/internal/modules/investment"
	"KFS_Backend/internal/modules/user"
	"KFS_Backend/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Config yüklenirken hata: %v", err)
	}

	sslMode := "require"
	if config.Database.SSLMode == "disable" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Database.Host, config.Database.User, config.Database.Password,
		config.Database.Name, config.Database.Port, sslMode,
	)

	var dbErr error
	for i := 0; i < 3; i++ {
		DB, dbErr = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Info),
		})

		if dbErr == nil {
			logger.Info("✅ Veritabanına başarıyla bağlandı!")
			break
		}

		logger.Warn(fmt.Sprintf("⚠️  Veritabanı bağlantı hatası: %v", dbErr))
		time.Sleep(3 * time.Second)
	}

	if dbErr != nil {
		logger.Error("❌ Veritabanına bağlanılamadı, uygulama durduruluyor!")
		log.Fatal(dbErr)
	}

	sqlDB, err := DB.DB()
	if err != nil || sqlDB.Ping() != nil {
		logger.Error("⚠️  Veritabanı bağlantısı başarısız!")
		log.Fatal(err)
	}

	//RunMigrations()
}

func RunMigrations() {
	logger.Info("🚀 Migration işlemi başlatılıyor...")

	if DB == nil {
		logger.Error("❌ Veritabanı bağlantısı boş, migration işlemi iptal edildi!")
		return
	}

	// 1️⃣ Önce bağımsız tablolar
	err := DB.AutoMigrate(
		&user.User{}, // Kullanıcı Modeli
	)
	if err != nil {
		logger.Error(fmt.Sprintf("❌ User tablosu oluşturulamadı: %v", err))
		log.Fatal(err)
	}

	// 2️⃣ Bağımlı tabloları migrate et
	err = DB.AutoMigrate(
		&user.AuthUser{},          // Kullanıcı Yetkilendirme
		&user.EmailVerification{}, // E-posta doğrulama
		&user.UserSession{},       // Kullanıcı Oturum Yönetimi
		&campaign.Campaign{},      // Kampanya Modeli
		&investment.Investment{},  // Yatırım Modeli
	)
	if err != nil {
		logger.Error(fmt.Sprintf("❌ Diğer tablolar oluşturulamadı: %v", err))
		log.Fatal(err)
	}

	logger.Info("✅ Veritabanı migrasyonu tamamlandı!")
}
