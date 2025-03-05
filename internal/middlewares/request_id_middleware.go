package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)


func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Eğer istemci request ID göndermişse onu kullan
		rid := c.Get("X-Request-ID")
		if rid == "" {
			// Yoksa yeni bir request ID oluştur
			rid = uuid.New().String()
		}

		// Request ID'yi context'e kaydet
		c.Locals("requestID", rid)
		// Response header'a ekle
		c.Set("X-Request-ID", rid)

		return c.Next()
	}
}

// RequestIDMiddleware, her HTTP isteği için benzersiz bir tanımlayıcı (request ID) oluşturan ve yöneten middleware fonksiyonudur.
// Bu middleware şunları yapar:
// - Gelen istekte X-Request-ID header'ı varsa, bu değeri kullanır
// - Yoksa, yeni bir UUID oluşturur
// - Request ID'yi context içinde saklar ve response header'ına ekler
// - Bu sayede her isteğin takibi ve loglama işlemleri daha kolay hale gelir