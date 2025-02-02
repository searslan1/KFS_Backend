package main

import (
	"c:/Users/gulay/OneDrive/Masaüstü/KFS/KFS_Backend/internal/db"
	"log"
	
)

func main() {
	// Veritabanı bağlantı dizesi
	dataSourceName := "user=yourusername dbname=yourdbname sslmode=disable password=yourpassword"

	// Veritabanı bağlantısını başlat
	db.InitDB(dataSourceName)

	// Migrasyon dosyasının yolunu belirtin
	migrationFilePath := "c:/Users/gulay/OneDrive/Masaüstü/KFS/KFS_Backend/migrations/001_create_tables.sql"

	// Migrasyonları çalıştır
	db.RunMigrations(migrationFilePath)

	// ...diğer kodlar...
	log.Println("Application started successfully")
}
