package middlewares

import "github.com/gofiber/fiber/v2"

func SecurityHeadersMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// XSS koruması
		c.Set("X-XSS-Protection", "1; mode=block")
		// Clickjacking koruması
		c.Set("X-Frame-Options", "DENY")
		// MIME type sniffing koruması
		c.Set("X-Content-Type-Options", "nosniff")
		// Referrer Policy
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		// Content Security Policy
		c.Set("Content-Security-Policy", "default-src 'self'")
		// HSTS (HTTPS zorunluluğu)
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		// Cache kontrolü
		c.Set("Cache-Control", "no-store, no-cache, must-revalidate")

		return c.Next()
	}
}
