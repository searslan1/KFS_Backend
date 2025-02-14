package middlewares

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/time/rate"
)

// Yeni: RateLimitConfig yapılandırma parametrelerini saklamak için.
type RateLimitConfig struct {
	Requests        int           // Saniyede izin verilen istek sayısı
	Burst           int           // Maksimum burst limiti
	CleanupDuration time.Duration // Temizlik süresi; süre boyunca aktif olmayan IP'ler silinecek
	ErrorMessage    string        // Hata durumunda dönecek mesaj
}

// visitor yapısı, ilgili IP için rate limiter (limitleyici) ve son erişim zamanını saklar.
type visitor struct {
	limiter  *rate.Limiter // IP bazında isteklere izin verme kontrolleri için kullanılır.
	lastSeen time.Time     // Bu IP'nin en son ne zaman istek gönderdiğini saklar.
}

// visitors global olarak tüm IP'ler ve onların visitor yapısı tutuluyor.
var visitors = sync.Map{}

// getVisitor fonksiyonu:
// Verilen IP adresine ait limitter'ı getirir ya da eğer yoksa oluşturur.
func getVisitor(ip string, config RateLimitConfig) *rate.Limiter {
	// visitors map'inde IP'ye ait visitor var mı diye kontrol ediyoruz.
	v, ok := visitors.Load(ip)
	if !ok {
		// Eğer yoksa, limiter'ı config parametresine göre oluşturuyoruz.
		l := rate.NewLimiter(rate.Limit(config.Requests), config.Burst)
		// Yeni visitor oluşturup kaydediyoruz.
		visitors.Store(ip, &visitor{
			limiter:  l,
			lastSeen: time.Now(),
		})
		return l
	}
	// Eğer visitor varsa, son erişim zamanını güncelliyoruz.
	vis := v.(*visitor)
	vis.lastSeen = time.Now()
	// Mevcut limiter'ı döndürüyoruz.
	return vis.limiter
}

// cleanupVisitors fonksiyonu:
// Belirtilen inactiveDuration süresi boyunca hiçbir etkinlik göstermeyen IP'leri visitors map'inden temizler.
func cleanupVisitors(inactiveDuration time.Duration) {
	for {
		// cleanup işlemleri için inactiveDuration kadar bekler.
		time.Sleep(inactiveDuration)
		now := time.Now()
		// visitors map'inde ki her öğeyi kontrol ediyoruz.
		visitors.Range(func(key, value interface{}) bool {
			v := value.(*visitor)
			// Eğer IP'nin son erişim zamanı, inactiveDuration süresinden eskiyse, siliniyor.
			if now.Sub(v.lastSeen) > inactiveDuration {
				visitors.Delete(key)
			}
			return true
		})
	}
}

// RateLimitMiddleware fonksiyonu:
// Fiber middleware olarak her istek için IP bazında rate limit kontrolü yapar.
func RateLimitMiddleware(config RateLimitConfig) fiber.Handler {
	// Varsayılan değerler atama
	if config.Requests == 0 {
		config.Requests = 1
	}
	if config.Burst == 0 {
		config.Burst = 5
	}
	if config.CleanupDuration == 0 {
		config.CleanupDuration = 5 * time.Minute
	}
	if config.ErrorMessage == "" {
		config.ErrorMessage = "Çok fazla istek. Lütfen bir süre sonra tekrar deneyin."
	}

	// Uygulama başlatıldığında, arka planda inactive IP'leri temizleyen goroutine başlatılır.
	go cleanupVisitors(config.CleanupDuration)

	// Middleware fonksiyonu, her istek için çalışır.
	return func(c *fiber.Ctx) error {
		// İstek yapan IP adresini alıyoruz.
		ip := c.IP()
		// IP'ye ait limiter'ı alıyoruz.
		limiter := getVisitor(ip, config)

		// Eğer limiter istek izni vermiyorsa, HTTP 429 (Too Many Requests) döndürüyoruz.
		if !limiter.Allow() {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": config.ErrorMessage,
			})
		}

		// Eğer izin verildiyse, zincirdeki sonraki middleware veya handler'a geçiyoruz.
		return c.Next()
	}
}

/*
RateLimitMiddleware Açıklaması:
- Bu middleware, gelen her istek için istemcinin IP adresine göre istek limiti kontrolü yapar.
- RateLimitConfig yapısı ile:
  • Saniyede izin verilen istek sayısı (Requests),
  • Maksimum burst limiti (Burst),
  • Temizlik süresi (CleanupDuration) ve
  • Hata mesajı (ErrorMessage)
  ayarlanabilir.
- Her istek geldiğinde, getVisitor fonksiyonu ilgili IP’ye ait bir rate limiter oluşturur veya mevcut olanı günceller.
- Eğer rate limiter, belirlenen limiti aşmışsa HTTP 429 (Too Many Requests) yanıtı döndürülür.
- Arka planda çalışan cleanupVisitors fonksiyonu, belirlenen süre boyunca etkinlik göstermeyen IP'leri otomatik olarak temizler.
Kullanım:
- Fiber uygulamanızda bu middleware’i global olarak eklemek için, örneğin main.go dosyasında app.Use(RateLimitMiddleware(config)) şeklinde çağırabilirsiniz.
- Böylece, tüm rotalar için otomatik IP bazlı istek sınırlandırması sağlanır.
- Bu middleware DDOS saldırılarına karşı koruma sağlar ve sunucunun aşırı yüklenmesini engeller.
*/