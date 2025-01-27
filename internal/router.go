package internal

import (
	"KFS_Backend/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

// Router oluşturma
func SetupRouter(app *fiber.App) {
	// Middleware ekleme
	app.Use(middlewares.RateLimiter())

	// Sağlık kontrol endpointi
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Buraya modül bazlı rotalar eklenecek
}
