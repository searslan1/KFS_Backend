package internal

import (
	"KFS_Backend/internal/middlewares"
	"github.com/gofiber/fiber/v2"
	"KFS_Backend/internal/modules/auth"

)

// Router oluşturma
func SetupRouter1(app *fiber.App) {
	// Middleware ekleme
	app.Use(middlewares.RateLimiter())

	// Sağlık kontrol endpointi
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	authRepo := &auth.AuthRepository{DB: DB}
	authService := &auth.AuthService{Repo: authRepo}
	authController := &auth.AuthController{Service: authService}

	app.Post("/register", authController.RegisterHandler)
	app.Post("/login", authController.LoginHandler)

}
