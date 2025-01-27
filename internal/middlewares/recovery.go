package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

// Panik hatalarını yakalayan middleware
func RecoveryMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if err := recover(); err != nil {
				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Sunucu hatası! 🛑",
				})
			}
		}()
		return c.Next()
	}
}
