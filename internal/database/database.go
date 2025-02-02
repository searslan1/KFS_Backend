package database

import (
	"fmt"
	"log"
	"time"

	"KFS_Backend/configs"
	"KFS_Backend/internal/modules/campaign"
	"KFS_Backend/internal/modules/investment"
	"KFS_Backend/internal/modules/auth"
	"KFS_Backend/internal/utils"
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

	// 🔥 Önce User tablosunu kontrol edip oluşturuyoruz
	var tableExists bool
	tableExists = DB.Migrator().HasTable(&auth.User{})
	
	if !tableExists {
		logger.Info("🔹 users tablosu oluşturuluyor...")
		err := DB.AutoMigrate(&auth.User{})
		if err != nil {
			logger.Error(fmt.Sprintf("❌ users tablosu oluşturulamadı: %v", err))
			log.Fatal(err)
		}
	} else {
		logger.Info("✅ users tablosu zaten mevcut, yeniden oluşturulmayacak.")
	}
	

	// Diğer tabloları kontrol et
	mainTables := map[string]interface{}{
		"auth_users":         &auth.AuthUser{},
		// "email_verifications": &user.EmailVerification{},
		// "user_sessions":      &user.UserSession{},
		"campaigns":          &campaign.Campaign{},
		"investments":        &investment.Investment{},
		"verification_codes": &utils.Verification{},
	}

	for tableName, model := range mainTables {
		var tableCount int64
		DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE LOWER(table_name) = LOWER(?)", tableName).Scan(&tableCount)

		if tableCount == 0 {
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


