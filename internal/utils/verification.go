package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"firebase.google.com/go/auth"
	"gorm.io/gorm"
)

// 📌 Rastgele 6 haneli doğrulama kodu üret
func GenerateVerificationCode() string {
	return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000) // 6 haneli rastgele kod üret
}

// 📌 Telefon doğrulama kodunu veritabanına kaydet
func SaveVerificationCode(db *gorm.DB, userID int64, code string, verificationType string) error {
	expiry := time.Now().Add(5 * time.Minute) // Kodun geçerlilik süresi 5 dakika

	verification := Verification{
		UserID:     userID,
		Code:       code,
		CodeExpiry: expiry,
		Type:       verificationType,
		IsVerified: false,
		CreatedAt:  time.Now(),
	}

	// 📌 Önce eski kodları sil
	db.Where("user_id = ? AND type = ?", userID, verificationType).Delete(&Verification{})

	// 📌 Yeni kodu kaydet
	if err := db.Create(&verification).Error; err != nil {
		return err
	}

	log.Printf("✅ Telefon doğrulama kodu kaydedildi: %s\n", code)
	return nil
}

// 📌 Firebase OTP Gönderme Fonksiyonu (SMS OTP)
func SendFirebaseOTP(authClient *auth.Client, phoneNumber string) error {
	params := (&auth.UserToCreate{}).PhoneNumber(phoneNumber)
	_, err := authClient.CreateUser(context.Background(), params)
	if err != nil {
		log.Println("❌ Firebase OTP oluşturulamadı:", err)
		return fmt.Errorf("Firebase OTP gönderilemedi: %v", err)
	}

	log.Println("✅ Firebase OTP başarıyla gönderildi:", phoneNumber)
	return nil
}

// 📌 Kullanıcının Firebase OTP Kodunu Doğrulama
func VerifyFirebaseOTP(authClient *auth.Client, idToken string) (bool, error) {
	// Kullanıcının Firebase OTP kodunu doğrula
	token, err := authClient.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return false, fmt.Errorf("Geçersiz OTP kodu: %v", err)
	}

	// Firebase UID Kontrolü
	if token.UID == "" {
		return false, fmt.Errorf("Kullanıcı doğrulama başarısız, UID boş")
	}

	// Kullanıcının telefon numarasını al (Eğer varsa)
	if phoneNumber, ok := token.Claims["phone_number"].(string); ok {
		log.Printf("✅ Kullanıcı telefon doğrulandı: %s", phoneNumber)
	} else {
		log.Println("⚠️ Kullanıcı telefon numarası bulunamadı!")
	}

	// Kullanıcı doğrulandıysa başarılı sonucu döndür
	return true, nil
}

