package middlewares

import (
	"time"

	"KFS_Backend/internal/database"
	"KFS_Backend/internal/modules/user"
	"KFS_Backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Cookies("access_token")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token bulunamadı"})
		}

		db := database.DB
		var blacklistedToken user.BlacklistedToken
		if err := db.Where("token = ? AND expiry > ?", tokenString, time.Now()).First(&blacklistedToken).Error; err == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token blackliste alınmış"})
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Geçersiz veya süresi dolmuş token"})
		}

		var userAuth user.AuthUser
		if err := db.Where("user_id = ?", claims.UserID).First(&userAuth).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Kullanıcı bulunamadı"})
		}

		if !userAuth.AccountLockedUntil.IsZero() && userAuth.AccountLockedUntil.After(time.Now()) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Hesap kilitli"})
		}

		var currentUser user.User
		if err := db.Where("user_id = ?", userAuth.UserID).First(&currentUser).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Kullanıcı bilgileri alınamadı"})
		}

		if currentUser.UserType != claims.Role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Kullanıcı tipi uyuşmuyor"})
		}

		if !userAuth.EmailVerified {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Email doğrulanmamış"})
		}

		if !userAuth.PhoneVerified {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Telefon doğrulanmamış"})
		}

		c.Locals("userID", userAuth.UserID)
		c.Locals("role", currentUser.UserType)

		return c.Next()
	}
}

/*
AuthMiddleware - Kullanıcı kimlik doğrulama middleware'i

Kullanım Durumu (Use Case):
Bu middleware, kimlik doğrulama gerektiren API endpoint'lerini korumak için kullanılır.
Yalnızca geçerli access_token'a sahip ve tüm doğrulama adımlarını geçen kullanıcıların
korumalı kaynaklara erişmesine izin verir.

İşlevi:
1. İstekteki 'access_token' çerezinden JWT token'ı alır
2. Token'ın blacklist'te olup olmadığını kontrol eder
3. Token'ın geçerliliğini doğrular (imza ve süre)
4. Kullanıcının veritabanında var olduğunu doğrular
5. Hesap kilidi durumunu kontrol eder
6. Token'daki rolün kullanıcı tipiyle eşleşip eşleşmediğini doğrular
7. E-posta ve telefon doğrulamasını kontrol eder
8. Doğrulama başarılıysa, sonraki handler'lara kullanıcı kimliği ve rol bilgisini aktarır

Nasıl Kullanılır:
- Tüm rotalar için: app.Use(middlewares.AuthMiddleware())
- Belirli bir grup için: app.Group("/api").Use(middlewares.AuthMiddleware())
- Tek bir endpoint için: app.Get("/protected", middlewares.AuthMiddleware(), protectedHandler)

Başarısız doğrulama durumunda uygun HTTP durum kodu ve hata mesajı ile yanıt döner.
*/
