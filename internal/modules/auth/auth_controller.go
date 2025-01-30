package auth

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type AuthController struct {
	Service *AuthService
}

func (c *AuthController) RegisterHandler(ctx *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		UserType string `json:"user_type"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := c.Service.RegisterUser(req.Email, req.Password, req.UserType)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User successfully registered"})
}

func (c *AuthController) LoginHandler(ctx *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	user, err := c.Service.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(user)
}
// **Yeni: Kullanıcıyı ID'ye göre getir**
func (c *AuthController) GetUserByIDHandler(ctx *fiber.Ctx) error {
	// URL'den user_id'yi al
	userIDParam := ctx.Params("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Kullanıcıyı getir
	user, err := c.Service.GetUserByID(userID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(user)
}

func (c *AuthController) SendEmailOTPHandler(ctx *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := c.Service.SendEmailOTP(req.Email)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"message": "OTP sent to email"})
}

func (c *AuthController) SendPhoneOTPHandler(ctx *fiber.Ctx) error {
	var req struct {
		PhoneNumber string `json:"phone_number"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := c.Service.SendPhoneOTP(req.PhoneNumber)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"message": "OTP sent to phone"})
}
