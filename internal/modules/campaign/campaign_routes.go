package campaign

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes, campaign modülüne ait tüm HTTP route'larını tanımlar.
// Bu fonksiyonu, ana router'ınızda çağırarak campaign modülü endpoint'lerini sisteme ekleyebilirsiniz.
func RegisterRoutes(router *gin.Engine, controller *Controller) {
	campaignRoutes := router.Group("/campaigns")
	{
		campaignRoutes.POST("/", controller.CreateCampaign)   // Yeni kampanya oluşturur
		campaignRoutes.GET("/:id", controller.GetCampaignByID)  // Belirtilen ID'ye sahip kampanyayı getirir
		campaignRoutes.PUT("/:id", controller.UpdateCampaign)   // Kampanya güncelleme işlemi
		campaignRoutes.DELETE("/:id", controller.DeleteCampaign) // Kampanya silme işlemi
	}
}
