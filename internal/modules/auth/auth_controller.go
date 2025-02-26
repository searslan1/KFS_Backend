package auth

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
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
		Email    string `json:"email"`     // Kullanıcının e-posta adresi.
		Password string `json:"password"`  // Kullanıcının şifresi.
		UserType string `json:"user_type"` // Kullanıcının türü (örneğin: admin, user).
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
func (c *AuthController) LoginHandler(ctx *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	accessToken, refreshToken, err := c.Service.AuthenticateUser(req.Email, req.Password, ctx.IP(), ctx.Get("User-Agent"))
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Expires:  time.Now().Add(time.Minute * 15),
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Expires:  time.Now().Add(time.Hour * 24 * 7),
	})

	return ctx.JSON(fiber.Map{"message": "Login successful"})
}

func (c *AuthController) LogoutHandler(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Refresh token is required"})
	}

	err := c.Service.LogoutUser(refreshToken)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Expires:  time.Now().Add(-time.Hour),
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Expires:  time.Now().Add(-time.Hour),
	})

	return ctx.JSON(fiber.Map{"message": "User logged out successfully"})
}

func (c *AuthController) RefreshTokenHandler(ctx *fiber.Ctx) error {
	// ✅ Refresh Token’i Cookie’den al
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Refresh token is required"})
	}

	// ✅ IP ve User-Agent bilgilerini al, boş olup olmadığını kontrol et
	ip := ctx.IP()
	userAgent := ctx.Get("User-Agent")
	if userAgent == "" {
		userAgent = "Unknown-Device"
	}

	// ✅ Yeni Token oluşturma işlemi
	accessToken, newRefreshToken, err := c.Service.RefreshAccessToken(refreshToken, ip, userAgent)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// ✅ Yeni Access Token’ı Cookie olarak ekleyelim
	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax", // ✅ Çerezlerin tarayıcı tarafından engellenmemesi için "Lax"
		Expires:  time.Now().Add(time.Minute * 15),
	})

	// ✅ Yeni Refresh Token’ı Cookie olarak ekleyelim
	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax", // ✅ Daha güvenli ve tarayıcı uyumluluğunu artırır
		Expires:  time.Now().Add(time.Hour * 24 * 7),
	})

	// ✅ Yanıtı döndür
	return ctx.JSON(fiber.Map{"message": "Token refreshed successfully"})
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
		Email  string `json:"email"`   // Kullanıcının e-posta adresi.
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
		OTP    string `json:"otp"`     // Kullanıcıya gönderilen doğrulama kodu (OTP).
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
