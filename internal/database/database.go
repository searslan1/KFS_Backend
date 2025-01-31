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

	// 📌 Hazırlanmış ifadeleri devre dışı bırak
	pgConfig := postgres.Config{
		DSN: dsn,
		PreferSimpleProtocol: true, // 🔥 Hazırlanmış ifadeleri kapatıyoruz
	}

	var dbErr error
	for i := 0; i < 3; i++ {
		DB, dbErr = gorm.Open(postgres.New(pgConfig), &gorm.Config{
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

	RunMigrations()
}


func RunMigrations() {
	logger.Info("🚀 Migration işlemi başlatılıyor...")

	if DB == nil {
		logger.Error("❌ Veritabanı bağlantısı boş, migration işlemi iptal edildi!")
		return
	}

	// Ana tablolar (bağımsız olanlar)
	mainTables := map[string]interface{}{
		"users":              &user.User{},
		"auth_users":         &user.AuthUser{},
		"email_verifications": &user.EmailVerification{},
		"user_sessions":      &user.UserSession{},
		"campaigns":          &campaign.Campaign{},
		"investments":        &investment.Investment{}, // Investment modeli zaten campaign_investments tablosuna bağlanıyor
	}

	// Ana tabloları kontrol et ve eksik olanları oluştur
	for tableName, model := range mainTables {
		var tableCount int64
		err := DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = ?", tableName).Scan(&tableCount).Error
		if err != nil {
			logger.Error(fmt.Sprintf("❌ %s tablosu kontrol edilirken hata oluştu: %v", tableName, err))
			log.Fatal(err)
		}

		if tableCount == 0 {
			// Eğer tablo yoksa, oluştur
			logger.Info(fmt.Sprintf("🔹 %s tablosu oluşturuluyor...", tableName))
			err := DB.AutoMigrate(model)
			if err != nil {
				logger.Error(fmt.Sprintf("❌ %s tablosu oluşturulamadı: %v", tableName, err))
				log.Fatal(err)
			}
		} else {
			logger.Info(fmt.Sprintf("✅ %s tablosu zaten mevcut, yeniden oluşturulmayacak.", tableName))
		}
	}

	logger.Info("✅ Veritabanı migrasyonu tamamlandı!")
}


