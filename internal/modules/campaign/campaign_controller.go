package campaign

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Controller, campaign modülü için HTTP isteklerini işleyen yapıdır.
type Controller struct {
	service Service
}

// NewController, campaign modülü için yeni bir controller örneği oluşturur.
func NewController(service Service) *Controller {
	return &Controller{service: service}
}

// CreateCampaign, POST /campaigns endpoint'ini işleyerek yeni kampanya oluşturur.
func (c *Controller) CreateCampaign(ctx *gin.Context) {
	var campaign CampaignProfile
	if err := ctx.ShouldBindJSON(&campaign); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.service.CreateCampaign(&campaign); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, campaign)
}

// GetCampaignByID, GET /campaigns/:id endpoint'iyle belirli bir kampanyayı getirir.
func (c *Controller) GetCampaignByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}
	campaign, err := c.service.GetCampaignByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, campaign)
}

// UpdateCampaign, PUT /campaigns/:id endpoint'iyle kampanya verilerini günceller.
func (c *Controller) UpdateCampaign(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	var campaign CampaignProfile
	if err := ctx.ShouldBindJSON(&campaign); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	campaign.ID = uint(id)

	if err := c.service.UpdateCampaign(&campaign); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, campaign)
}

// DeleteCampaign, DELETE /campaigns/:id endpoint'iyle kampanyayı siler.
func (c *Controller) DeleteCampaign(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	if err := c.service.DeleteCampaign(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Campaign deleted successfully"})
}
