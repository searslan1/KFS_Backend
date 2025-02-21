package utils

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// 📌 Firebase'i başlatan fonksiyon
func InitFirebase() (*firebase.App, error) {
	opt := option.WithCredentialsFile("config/serviceAccountKey.json") // JSON dosyanın yolunu doğru ayarla
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("🔥 Firebase başlatılamadı: %v", err)
		return nil, err
	}
	log.Println("✅ Firebase başarıyla başlatıldı")
	return app, nil
}
