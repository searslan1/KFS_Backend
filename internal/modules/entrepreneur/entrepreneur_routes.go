package routes

import (
	"entrepreneur/controller"
	"github.com/gofiber/fiber/v2"
)

// RegisterEntrepreneurRoutes girişimci profilleri için API endpointlerini tanımlar
func RegisterEntrepreneurRoutes(app *fiber.App, entrepreneurController *controller.EntrepreneurController) {
	entrepreneurRoutes := app.Group("/entrepreneurs")
	{
		entrepreneurRoutes.Post("/", entrepreneurController.CreateEntrepreneur)           // Yeni girişimci profili oluştur
		entrepreneurRoutes.Put("/:id/approve", entrepreneurController.AdminApproveEntrepreneur
