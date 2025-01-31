package user

import (
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	service *UserService
}

func NewUserController(service *UserService) *UserController {
	return &UserController{service}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz giriş"})
	}

	err := c.service.RegisterUser(req.Email, req.Password, req.Role)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"message": "Kullanıcı kaydedildi"})
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz giriş"})
	}

	ip := ctx.IP()
	deviceID := ctx.Get("User-Agent")

	accessToken, refreshToken, err := c.service.LoginUser(req.Email, req.Password, ip, deviceID)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Hatalı giriş"})
	}

	return ctx.JSON(fiber.Map{"access_token": accessToken, "refresh_token": refreshToken})
}
