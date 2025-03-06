package campaign

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes, campaign modülüne ait tüm HTTP route'larını tanımlar.
func RegisterRoutes(app *fiber.App, controller *Controller) {
	campaignRoutes := app.Group("/campaigns")
	{
		campaignRoutes.Post("/", controller.CreateCampaign)
		campaignRoutes.Get("/:id", controller.GetCampaignByID)
		campaignRoutes.Put("/:id", controller.UpdateCampaign)
		campaignRoutes.Delete("/:id", controller.DeleteCampaign)
	}
}
