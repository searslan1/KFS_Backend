package middlewares

import (
	"KFS_Backend/internal/database"
	"KFS_Backend/internal/modules/user"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	RoleAdmin      = "admin"
	RoleModerator  = "moderator"
	RoleCorporate  = "corporate"
	RoleIndividual = "individual"
)

const (
	STRICT = iota
	ANY
)

func checkUserRoles(userRole string, checkType int, roles []string) bool {
	roleMap := make(map[string]bool, len(roles))
	for _, role := range roles {
		roleMap[role] = true
	}

	switch checkType {
	case STRICT:
		return roleMap[userRole]
	case ANY:
		if userRole == RoleAdmin {
			return true
		}
		return roleMap[userRole]
	default:
		return false
	}
}

func Authorize(checkType int, roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("userID").(int64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Kullanıcı kimliği bulunamadı, lütfen önce giriş yapın",
			})
		}

		cachedRole, roleExists := c.Locals("role").(string)

		if !roleExists {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Rol bilgisi bulunamadı, lütfen tekrar giriş yapın",
			})
		}

		db := database.DB
		var currentUser user.User
		if err := db.Select("user_type").First(&currentUser, userID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Kullanıcı bulunamadı",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Rol doğrulaması sırasında hata oluştu",
			})
		}

		if cachedRole != currentUser.UserType {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Kullanıcı tipi tutarsız, lütfen tekrar giriş yapın",
			})
		}

		if !checkUserRoles(currentUser.UserType, checkType, roles) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Bu sayfaya erişim yetkiniz bulunmuyor",
			})
		}

		return c.Next()
	}
}

func RequireAdmin() fiber.Handler {
	return Authorize(STRICT, RoleAdmin)
}

func RequireModerator() fiber.Handler {
	return Authorize(ANY, RoleAdmin, RoleModerator)
}

func RequireCorporate() fiber.Handler {
	return Authorize(ANY, RoleAdmin, RoleModerator, RoleCorporate)
}

func RequireIndividual() fiber.Handler {
	return Authorize(ANY, RoleAdmin, RoleModerator, RoleIndividual)
}

func RequireAllUsers() fiber.Handler {
	return Authorize(ANY, RoleAdmin, RoleModerator, RoleCorporate, RoleIndividual)
}

/*
Authorize Middleware - Kullanıcı yetkilendirme middleware'i

Kullanım Durumu (Use Case):
Bu middleware, belirli rollerin erişimine izin verilen API endpoint'lerini korumak için kullanılır.
Yalnızca belirtilen rollerden birine sahip kullanıcıların korumalı kaynaklara erişmesine izin verir.

İşlevi:
1. Kullanıcı kimliğini ve rolünü context'ten alır
2. Kullanıcı kimliği veya rol bilgisi bulunamazsa yetkisiz erişim hatası döner
3. Veritabanından kullanıcının rolünü doğrular
4. Kullanıcının rolü belirtilen rollerden biriyle eşleşmezse erişim izni vermez
5. Doğrulama başarılıysa, sonraki handler'lara geçiş yapar

Nasıl Kullanılır:
- Tüm rotalar için: app.Use(middlewares.Authorize(STRICT, "admin"))
- Belirli bir grup için: app.Group("/admin").Use(middlewares.RequireAdmin())
- Tek bir endpoint için: app.Get("/admin/dashboard", middlewares.RequireAdmin(), adminDashboardHandler)

Zincirleme Çalışma Mantığı:
- AuthMiddleware'den sonra kullanılmalıdır,çünkü önce AuthMiddleware çalışır ve kullanıcı kimliğini doğrular.
- Ardından Authorize middleware çalışır ve kullanıcının rolünü kontrol eder.
- Her iki doğrulama da başarılı olursa, istek bir sonraki handler'a geçer.
- Başarısız doğrulama durumunda uygun HTTP durum kodu ve hata mesajı ile yanıt döner.
*/
