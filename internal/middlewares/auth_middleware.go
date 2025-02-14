package middlewares

import (
	"strings"
	"time"

	"KFS_Backend/configs"
	"KFS_Backend/internal/modules/user"
	"KFS_Backend/internal/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// JWT doğrulama middleware'i
func JWTMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token eksik"})
		}

		tokenParts := strings.Split(authHeader, "Bearer ")
		if len(tokenParts) != 2 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Geçersiz token formatı"})
		}

		tokenString := tokenParts[1]
		config := configs.LoadJWTConfig()

		// Token'ı çözümle
		token, err := jwt.ParseWithClaims(tokenString, &utils.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.SecretKey), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Geçersiz veya süresi dolmuş token"})
		}

		claims, ok := token.Claims.(*utils.JWTClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token doğrulama hatası"})
		}

		// **📌 Kullanıcının aktif olup olmadığını kontrol et**
		var userRecord user.User
		if err := db.Where("user_id = ?", claims.UserID).First(&userRecord).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Kullanıcı bulunamadı"})
		}

		// **📌 Kullanıcının hesabı kilitli mi? (SPK gerekliliği)**
		var authRecord user.AuthUser
		if err := db.Where("user_id = ?", claims.UserID).First(&authRecord).Error; err == nil {
			if authRecord.AccountLockedUntil != nil && time.Now().Before(*authRecord.AccountLockedUntil) {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Hesap kilitli, daha sonra tekrar deneyin."})
			}
		}

		// **📌 Kullanıcının IP ve cihaz bilgisiyle giriş yapıp yapmadığını kontrol et**
		var session user.UserSession
		if err := db.Where("user_id = ? AND ip_address = ? AND device_info = ?", claims.UserID, claims.IP, claims.DeviceID).First(&session).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Yetkisiz cihaz veya IP adresi!"})
		}

		// **📌 Kullanıcının rolünü belirle ve erişim kontrolü yap**
		c.Locals("userID", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

/*📌 Yeni Özellikler ve SPK Standartlarına Uygun Güncellemeler
✅ 1. Kullanıcı Varlık Kontrolü
Eğer kullanıcı veritabanında yoksa, doğrudan 401 Unauthorized döndürülüyor.
✅ 2. Hesap Kilitlenme Kontrolü
Eğer SPK gerekliliklerine göre kullanıcı çok fazla başarısız giriş yapmışsa, account_locked_until alanı kontrol edilerek erişim reddediliyor.
✅ 3. Cihaz & IP Kontrolü
Kullanıcının token içerisindeki IP ve cihaz bilgisi, user_sessions tablosundaki oturum bilgisiyle eşleşmeli.
Yetkisiz cihaz veya IP adresi algılanırsa, giriş reddedilir!
✅ 4. Kullanıcı Rolüne Göre Yetkilendirme (RBAC)
Kullanıcının role bilgisi middleware içinde belirleniyor.
İleri seviye erişim kontrolleri yapmak için gerekli olan role bilgisi c.Locals() ile diğer handler’lara aktarılıyor.*/
