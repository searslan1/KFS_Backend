package profile

import "github.com/gofiber/fiber/v2"

func RegisterProfileRoutes(app *fiber.App, controller *ProfileController) {
	app.Get("/profile/:id", controller.GetUserProfileHandler)
	app.Get("/profile/:id/address", controller.GetUserAddressHandler)
	app.Get("/profile/:id/media", controller.GetUserMediaFilesHandler)
	app.Get("/profile/:id/role", controller.GetUserRoleHandler)
	app.Put("/profile/:id", controller.UpdateUserProfileHandler)
	app.Put("/profile/:id/address", controller.UpdateUserAddressHandler)
	app.Put("/profile/:id/role", controller.UpdateUserRoleHandler)
}
