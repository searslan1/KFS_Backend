package internal

import (
	"KFS_Backend/internal/middlewares"
	"KFS_Backend/internal/modules/auth"
	"github.com/gofiber/fiber/v2"
	"KFS_Backend/internal/database"
)

// SetupRouter uygulamanın tüm route'larını tanımlar
func SetupRouter(app *fiber.App, userController *auth.AuthController) {
	// Global Middleware'ler
	app.Use(middlewares.RateLimitMiddleware(middlewares.RateLimitConfig{}))

	// Sağlık kontrol endpointi
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Veritabanı bağlantısını al

	// Auth modülü bağımlılıklarını oluştur
	authRepo := &auth.AuthRepository{DB: database.GetDB()}
	authService := &auth.AuthService{Repo: authRepo}
	authController := &auth.AuthController{Service: authService}

	auth.RegisterUserRoutes(app, authController)

	// Authentication rotaları
	app.Post("/register", authController.RegisterHandler)
	app.Post("/login", authController.LoginHandler)
}