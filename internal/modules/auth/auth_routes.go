package auth

import "github.com/gofiber/fiber/v2"

func RegisterUserRoutes(app *fiber.App, controller *AuthController) {
	app.Post("/users", controller.RegisterHandler)
	app.Post("/auth/login", controller.LoginHandler)
	app.Get("/auth/users", controller.GetAllUsersHandler)
	app.Get("/auth/user/:id", controller.GetUserByIDHandler)
	app.Post("/auth/send-email-verification", controller.SendEmailVerificationHandler)
	app.Post("/auth/verify-email", controller.VerifyEmailHandler)

}
