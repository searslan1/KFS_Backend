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
