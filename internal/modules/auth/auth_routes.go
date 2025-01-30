package auth

import "github.com/gofiber/fiber/v2"

func RegisterUserRoutes(app *fiber.App, controller *AuthController) {
	app.Post("/users", controller.RegisterHandler)
	app.Post("/auth/login", controller.LoginHandler)
	app.Get("/auth/user/:id", controller.GetUserByIDHandler)
	app.Post("/auth/email/send-otp", controller.SendEmailOTPHandler)
	app.Post("/auth/phone/send-otp", controller.SendPhoneOTPHandler)
}
