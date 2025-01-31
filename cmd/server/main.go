package main

import (
	"KFS_Backend/internal"
	"KFS_Backend/internal/database"
	"KFS_Backend/pkg/logger"
)

func main() {
	logger.InitLogger()

	// Veritabanı bağlantısı ve migrasyon
	database.ConnectDatabase()

	// Sunucuyu başlat
	internal.StartServer()
}
