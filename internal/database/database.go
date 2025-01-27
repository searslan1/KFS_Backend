package database

import (
	"fmt"
	"log"
	"os/user"
	"time"

	"KFS_Backend/configs"
	"KFS_Backend/pkg/logger"

	// Modülleri import et
	"KFS_Backend/internal/modules/campaign"
	"KFS_Backend/internal/modules/investment"

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

	// DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
		config.Database.Host, config.Database.User, config.Database.Password,
		config.Database.Name, config.Database.Port,
	)

	// 3 kez tekrar deneme mekanizması
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
		time.Sleep(3 * time.Second) // 3 saniye bekleyerek tekrar dene
	}

	if dbErr != nil {
		logger.Error("❌ Veritabanına bağlanılamadı, uygulama durduruluyor!")
		log.Fatal(dbErr)
	}

	// Bağlantı Testi
	sqlDB, _ := DB.DB()
	if err := sqlDB.Ping(); err != nil {
		logger.Error("⚠️  Veritabanı bağlantısı başarısız!")
		log.Fatal(err)
	}

	// **Migration Çalıştır**
	RunMigrations()
}

func RunMigrations() {
	err := DB.AutoMigrate(
		&user.User{},             // Kullanıcı Modeli
		&campaign.Campaign{},     // Kampanya Modeli
		&investment.Investment{}, // Yatırım Modeli
	)
	if err != nil {
		logger.Error("❌ Veritabanı migrasyonu başarısız!")
	} else {
		logger.Info("✅ Veritabanı migrasyonu tamamlandı!")
	}
}
/*database: Veritabanı işlemleri ile ilgili tüm işlevleri yönetir. 
Yani, veritabanı bağlantılarını açma, migrasyonları çalıştırma ve veritabanı sorguları gibi işlemleri burada yapabilirsiniz. 
Bu dosya, veritabanı ile ilgili tüm mantığı ve işlemleri içerdiğinden dolayı, 
uygulamanın veritabanı ile ilgili kodlarını daha modüler hale getirir.*/