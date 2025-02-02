package auth

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// AuthController yapısı, AuthService'i çağırarak istekleri yöneten bir kontrolör.
type AuthController struct {
	Service *AuthService // İş mantığını barındıran AuthService'e bir referans.
}

// Kullanıcı kayıt işlemini gerçekleştiren handler.
// İstek gövdesinden email, şifre ve kullanıcı türü alır.
// AuthService'deki RegisterUser fonksiyonunu çağırır.
func (c *AuthController) RegisterHandler(ctx *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`    // Kullanıcının e-posta adresi.
		Password string `json:"password"` // Kullanıcının şifresi.
		UserType string `json:"user_type"`// Kullanıcının türü (örneğin: admin, user).
	}

	// İstek gövdesini struct'a çevir.
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Kullanıcıyı kaydetmek için servis fonksiyonunu çağır.
	err := c.Service.RegisterUser(req.Email, req.Password, req.UserType)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Başarılı yanıt döner.
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User successfully registered"})
}

// Kullanıcı giriş işlemini gerçekleştiren handler.
// E-posta ve şifre alır, AuthService ile doğrular ve kullanıcı bilgilerini döner.
func (c *AuthController) LoginHandler(ctx *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`    // Kullanıcının e-posta adresi.
		Password string `json:"password"` // Kullanıcının şifresi.
	}

	// İstek gövdesini struct'a çevir.
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Kullanıcıyı doğrula.
	user, err := c.Service.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// Başarılı yanıt döner (kullanıcı bilgileri ile birlikte).
	return ctx.JSON(user)
}

// Kullanıcıyı ID'ye göre getiren handler.
// URL parametresinden user_id alır, AuthService'ten kullanıcıyı getirir.
func (c *AuthController) GetUserByIDHandler(ctx *fiber.Ctx) error {
	// URL'den user_id'yi al.
	userIDParam := ctx.Params("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64) // ID'yi int64'e çevir.
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Kullanıcıyı ID ile getir.
	user, err := c.Service.GetUserByID(userID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	// Başarılı yanıt döner (kullanıcı bilgileri ile birlikte).
	return ctx.JSON(user)
}

// Kullanıcıya e-posta doğrulama kodu gönderen handler.
// İstek gövdesinden user_id ve email alır, doğrulama kodunu üretir ve e-posta gönderir.
func (c *AuthController) SendEmailVerificationHandler(ctx *fiber.Ctx) error {
	var req struct {
		UserID int64  `json:"user_id"` // Kullanıcının ID'si.
		Email  string `json:"email"`  // Kullanıcının e-posta adresi.
	}

	// İstek gövdesini struct'a çevir.
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Doğrulama kodunu üret ve e-posta gönder.
	err := c.Service.SendEmailVerification(req.UserID, req.Email)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Başarılı yanıt döner.
	return ctx.JSON(fiber.Map{"message": "Verification email sent"})
}

// Kullanıcıyı e-posta doğrulama kodu ile doğrulayan handler.
// İstek gövdesinden user_id ve OTP alır, doğrulama işlemini yapar.
func (c *AuthController) VerifyEmailHandler(ctx *fiber.Ctx) error {
	var req struct {
		UserID int64  `json:"user_id"` // Kullanıcının ID'si.
		OTP    string `json:"otp"`    // Kullanıcıya gönderilen doğrulama kodu (OTP).
	}

	// İstek gövdesini struct'a çevir.
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// OTP'yi doğrula.
	err := c.Service.VerifyEmailOTP(req.UserID, req.OTP)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// Başarılı yanıt döner.
	return ctx.JSON(fiber.Map{"message": "Email verified successfully"})
}
func (c *AuthController) GetAllUsersHandler(ctx *fiber.Ctx) error {
	// Service katmanından tüm kullanıcıları al
	users, err := c.Service.GetAllUsers()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(users)
}
