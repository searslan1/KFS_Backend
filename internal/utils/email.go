package utils

import (
	"bytes"            // HTTP isteklerinde kullanılan veri tamponu için.
	"encoding/json"    // JSON formatında veri işlemek için.
	"fmt"              // Formatlı çıktı ve hata mesajları oluşturmak için.
	"io"               // HTTP yanıtlarının gövdesini okumak için.
	"log"              // Hata ve bilgi mesajlarını loglamak için.
	"net/http"         // HTTP isteklerini göndermek için.
	"os"               // Ortam değişkenlerini okumak için.
)

// SendGrid API'ye gönderilecek e-posta isteği yapısını tanımlayan struct
type SendGridEmailRequest struct {
	Personalizations []Personalization `json:"personalizations"` // Alıcı bilgilerini tutar
	From             EmailAddress      `json:"from"`             // Gönderen e-posta adresini tutar
	Subject          string            `json:"subject"`          // E-posta konusu
	Content          []Content         `json:"content"`          // E-posta içeriği
}

// E-posta alıcılarının bilgilerini tutan struct
type Personalization struct {
	To []EmailAddress `json:"to"` // Alıcı e-posta adresleri
}

// E-posta adresini tutan struct
type EmailAddress struct {
	Email string `json:"email"` // E-posta adresi
}

// E-posta içeriğini tutan struct
type Content struct {
	Type  string `json:"type"`  // İçerik türü (örneğin, "text/plain" veya "text/html")
	Value string `json:"value"` // İçerik metni
}

// SendEmailWithSendGrid: SendGrid API kullanarak e-posta gönderir
func SendEmailWithSendGrid(toEmail, subject, body string) error {
	// 1. API anahtarını kontrol et
	sendGridAPIKey := os.Getenv("SENDGRID_API_KEY") // API anahtarı ortam değişkeninden okunur
	if sendGridAPIKey == "" {
		// Eğer API anahtarı eksikse hata loglanır ve hata döndürülür
		log.Println("[ERROR] SendGrid API key is missing")
		return fmt.Errorf("SendGrid API key is missing")
	}

	// 2. Gönderen e-posta adresini kontrol et
	senderEmail := os.Getenv("EMAIL_SENDER") // Gönderen e-posta adresi ortam değişkeninden okunur
	if senderEmail == "" {
		// Eğer gönderen adres eksikse hata loglanır ve hata döndürülür
		log.Println("[ERROR] Email sender address is missing")
		return fmt.Errorf("Email sender address is missing")
	}

	// 3. SendGrid e-posta isteği oluştur
	emailRequest := SendGridEmailRequest{
		Personalizations: []Personalization{
			{
				To: []EmailAddress{{Email: toEmail}}, // Alıcı adresi
			},
		},
		From:    EmailAddress{Email: senderEmail}, // Gönderen adresi
		Subject: subject,                         // E-posta konusu
		Content: []Content{
			{
				Type:  "text/plain", // E-posta türü (düz metin)
				Value: body,         // E-posta içeriği
			},
		},
	}

	// 4. JSON formatına dönüştür
	requestBody, err := json.Marshal(emailRequest)
	if err != nil {
		// JSON formatına dönüştürme sırasında bir hata oluşursa loglanır ve döndürülür
		log.Printf("[ERROR] Failed to create request body: %v\n", err)
		return fmt.Errorf("failed to create request body: %w", err)
	}

	// 5. HTTP POST isteği oluştur
	req, err := http.NewRequest("POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(requestBody))
	if err != nil {
		// HTTP isteği oluşturma sırasında bir hata oluşursa loglanır ve döndürülür
		log.Printf("[ERROR] Failed to create HTTP request: %v\n", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 6. HTTP başlıklarını ayarla
	req.Header.Set("Authorization", "Bearer "+sendGridAPIKey) // API anahtarı
	req.Header.Set("Content-Type", "application/json")       // İçerik türü JSON olarak ayarlanır

	// 7. HTTP isteğini gönder
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// İstek gönderme sırasında bir hata oluşursa loglanır ve döndürülür
		log.Printf("[ERROR] Failed to send email: %v\n", err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	defer resp.Body.Close() // İstek tamamlandıktan sonra yanıt gövdesini kapat

	// 8. Yanıt kontrolü
	if resp.StatusCode != http.StatusAccepted {
		// Eğer yanıt kodu beklenen değerden farklıysa loglanır ve hata döndürülür
		body, _ := io.ReadAll(resp.Body) // Yanıt gövdesi okunur
		log.Printf("[ERROR] Failed to send email. Status code: %d, Response: %s\n", resp.StatusCode, string(body))
		return fmt.Errorf("failed to send email, status code: %d, response: %s", resp.StatusCode, string(body))
	}

	// 9. Başarılı durumda bilgi loglanır
	log.Printf("[INFO] Email sent successfully to %s\n", toEmail)
	return nil
}
