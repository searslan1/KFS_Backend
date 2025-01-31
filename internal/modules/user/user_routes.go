package user

import "github.com/gofiber/fiber/v2"

func SetupUserRoutes(app *fiber.App, userController *UserController) {
	userGroup := app.Group("/api/auth")

	userGroup.Post("/register", userController.Register)
	userGroup.Post("/login", userController.Login)
}
