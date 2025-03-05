package middlewares

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/time/rate"
)

// RateLimitConfig yapılandırma parametrelerini saklamak için.
type RateLimitConfig struct {
	Requests        int           // Saniyede izin verilen istek sayısı
	Burst           int           // Maksimum burst limiti
	CleanupDuration time.Duration // Temizlik süresi; süre boyunca aktif olmayan IP'ler silinecek
	ErrorMessage    string        // Hata durumunda dönecek mesaj
	TrustProxy      bool          // Proxy güvenirliği ayarı
	TokenBased      bool          // Token bazlı limitleme aktif mi?
	TokenHeader     string        // Token'ın bulunduğu header (TokenBased true ise)
	EnableLogging   bool          // Rate limit aşımı loglanacak mı?
}

// visitor yapısı, ilgili IP için rate limiter (limitleyici) ve son erişim zamanını saklar.
// Limiter ve token sayacı için iyileştirilmiş yapı
type visitor struct {
	limiter  *rate.Limiter // IP bazında isteklere izin verme kontrolleri için kullanılır.
	lastSeen time.Time     // Bu IP'nin en son ne zaman istek gönderdiğini saklar.
	mu       sync.Mutex    // Eklendi
}

// visitors global olarak tüm IP'ler ve onların visitor yapısı tutuluyor.
var visitors = sync.Map{}

// getVisitor fonksiyonu:
// Verilen IP adresine ait limitter'ı getirir ya da eğer yoksa oluşturur.
func getVisitor(ip string, config RateLimitConfig) (*visitor, bool) {
	// visitors map'inde IP'ye ait visitor var mı diye kontrol ediyoruz.
	v, exists := visitors.Load(ip)
	if !exists {
		// Eğer yoksa, limiter'ı config parametresine göre oluşturuyoruz.
		l := rate.NewLimiter(rate.Limit(config.Requests), config.Burst)
		// Yeni visitor oluşturup kaydediyoruz.
		vis := &visitor{
			limiter:  l,
			lastSeen: time.Now(),
		}
		visitors.Store(ip, vis)
		return vis, true // Yeni visitor oluşturuldu
	}
	// Eğer visitor varsa, son erişim zamanını güncelliyoruz.
	vis := v.(*visitor)
	vis.lastSeen = time.Now()
	// Mevcut visitor'ı döndürüyoruz.
	return vis, false // Mevcut visitor güncellendi
}

// cleanupVisitors fonksiyonu:
// Belirtilen inactiveDuration süresi boyunca hiçbir etkinlik göstermeyen IP'leri visitors map'inden temizler.
func cleanupVisitors(inactiveDuration time.Duration) {
	for {
		// cleanup işlemleri için inactiveDuration kadar bekler.
		time.Sleep(inactiveDuration)
		now := time.Now()
		// visitors map'inde ki her öğeyi kontrol ediyoruz.
		visitors.Range(func(key, value any) bool {
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
	if config.TokenBased && config.TokenHeader == "" {
		config.TokenHeader = "Authorization"
	}

	// Uygulama başlatıldığında, arka planda inactive IP'leri temizleyen goroutine başlatılır.
	go cleanupVisitors(config.CleanupDuration)

	// Middleware fonksiyonu, her istek için çalışır.
	return func(c *fiber.Ctx) error {
		var identifier string

		// Token bazlı limitleme aktif mi kontrol et
		if config.TokenBased {
			// Token header'dan al
			token := c.Get(config.TokenHeader)
			if token != "" {
				identifier = token
			} else {
				// Token bulunamadıysa IP kullan
				identifier = getClientIP(c, config.TrustProxy)
			}
		} else {
			// IP bazlı limitleme
			identifier = getClientIP(c, config.TrustProxy)
		}

		// IP veya token'a ait visitor'ı alıyoruz.
		vis, _ := getVisitor(identifier, config)

		vis.mu.Lock()
		allowed := vis.limiter.Allow()
		remaining := int(float64(vis.limiter.Burst()) - vis.limiter.Tokens())
		vis.mu.Unlock()

		c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", config.Burst))
		c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Second).Unix()))

		// Eğer limiter istek izni vermiyorsa, HTTP 429 (Too Many Requests) döndürüyoruz.
		if !allowed {
			// OWASP tavsiyesi: Retry-After header'ı ekle
			c.Set("Retry-After", "60")

			// İsteği logla
			if config.EnableLogging {
				log.Printf("Rate limit aşıldı: %s", identifier)
			}

			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": config.ErrorMessage,
			})
		}

		// Eğer izin verildiyse, zincirdeki sonraki middleware veya handler'a geçiyoruz.
		return c.Next()
	}
}

// getClientIP özellikle proxy arkasında gerçek IP adresini almak için güvenilir bir yöntem sunar
func getClientIP(c *fiber.Ctx, trustProxy bool) string {
	if trustProxy {
		// X-Forwarded-For header'ı birden fazla IP içerebilir, ilk geçerli olanı al
		if xff := c.Get("X-Forwarded-For"); xff != "" {
			// Virgülle ayrılmış IP'lerden ilkini al (gerçek istemci IP'si)
			ips := strings.Split(xff, ",")
			if len(ips) > 0 {
				return strings.TrimSpace(ips[0])
			}
		}

		// Diğer güvenilir proxy header'larını kontrol et
		if xrip := c.Get("X-Real-IP"); xrip != "" {
			return xrip
		}

		if cfip := c.Get("CF-Connecting-IP"); cfip != "" { // Cloudflare
			return cfip
		}
	}

	// Hiçbir header bulunamazsa doğrudan IP'yi al
	return c.IP()
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
