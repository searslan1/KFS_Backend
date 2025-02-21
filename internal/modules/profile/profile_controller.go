package profile

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type ProfileController struct {
	Service *ProfileService
}

func NewProfileController(service *ProfileService) *ProfileController {
	return &ProfileController{Service: service}
}

// Kullanıcının profil bilgilerini getir
func (c *ProfileController) GetUserProfileHandler(ctx *fiber.Ctx) error {
	userIDParam := ctx.Params("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	profile, err := c.Service.GetUserProfile(userID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(profile)
}

// Kullanıcının adres bilgilerini getir
func (c *ProfileController) GetUserAddressHandler(ctx *fiber.Ctx) error {
	userIDParam := ctx.Params("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	address, err := c.Service.GetUserAddress(userID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(address)
}

// Kullanıcının medya dosyalarını getir
func (c *ProfileController) GetUserMediaFilesHandler(ctx *fiber.Ctx) error {
	userIDParam := ctx.Params("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	mediaFiles, err := c.Service.GetUserMediaFiles(userID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(mediaFiles)
}

// Kullanıcının rol bilgilerini getir
func (c *ProfileController) GetUserRoleHandler(ctx *fiber.Ctx) error {
	userIDParam := ctx.Params("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	role, err := c.Service.GetUserRole(userID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(role)
}

// Kullanıcının profilini güncelle
func (c *ProfileController) UpdateUserProfileHandler(ctx *fiber.Ctx) error {
	userIDParam := ctx.Params("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var updatedProfile UserProfile
	if err := ctx.BodyParser(&updatedProfile); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	updatedProfile.UserID = userID // UserID'yi set et

	if err := c.Service.UpdateUserProfile(userID, &updatedProfile); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"message": "Profile updated successfully"})
}

func (c *ProfileController) UpdateUserAddressHandler(ctx *fiber.Ctx) error {
	userIDParam := ctx.Params("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var address Address
	if err := ctx.BodyParser(&address); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	address.UserID = userID // UserID'yi set et

	// Adres güncelleme veya ekleme işlemi
	if err := c.Service.UpdateUserAddress(userID, &address); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"message": "Address updated successfully"})
}


func (c *ProfileController) UpdateUserRoleHandler(ctx *fiber.Ctx) error {
	userIDParam := ctx.Params("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var role RoleProfile
	if err := ctx.BodyParser(&role); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Gelen role'deki user_id yerine API'den gelen user_id set ediliyor
	role.UserID = userID

	// Role güncelleme veya ekleme işlemi
	if err := c.Service.UpdateUserRole(userID, &role); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"message": "Role updated successfully"})
}

