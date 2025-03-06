package campaign

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
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
func (c *Controller) CreateCampaign(ctx *fiber.Ctx) error {
	var campaign CampaignProfile
	if err := ctx.BodyParser(&campaign); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := c.service.CreateCampaign(&campaign); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(http.StatusCreated).JSON(campaign)
}

// GetCampaignByID, GET /campaigns/:id endpoint'iyle belirli bir kampanyayı getirir.
func (c *Controller) GetCampaignByID(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid campaign ID"})
	}
	campaign, err := c.service.GetCampaignByID(uint(id))
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(http.StatusOK).JSON(campaign)
}

// UpdateCampaign, PUT /campaigns/:id endpoint'iyle kampanya verilerini günceller.
func (c *Controller) UpdateCampaign(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid campaign ID"})
	}

	var campaign CampaignProfile
	if err := ctx.BodyParser(&campaign); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	campaign.ID = uint(id)

	if err := c.service.UpdateCampaign(&campaign); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(http.StatusOK).JSON(campaign)
}

// DeleteCampaign, DELETE /campaigns/:id endpoint'iyle kampanyayı siler.
func (c *Controller) DeleteCampaign(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid campaign ID"})
	}

	if err := c.service.DeleteCampaign(uint(id)); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "Campaign deleted successfully"})
}
