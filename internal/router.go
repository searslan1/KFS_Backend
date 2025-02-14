package internal

import (
	"KFS_Backend/internal/middlewares"
	"KFS_Backend/internal/modules/user"

	"github.com/gofiber/fiber/v2"
)

// SetupRouter uygulamanın tüm route'larını tanımlar
func SetupRouter(app *fiber.App, userController *user.UserController) {
	// Global Middleware'ler
	app.Use(middlewares.RateLimiter())

	// Sağlık kontrol endpointi
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Kullanıcı modülü rotaları
	user.SetupUserRoutes(app, userController)  // 

	// Diğer modülleri buraya ekleyebiliriz (örneğin kampanya, yatırım, admin)
}
