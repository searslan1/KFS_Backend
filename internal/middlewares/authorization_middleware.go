package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func Authorize(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Rol bulunamadı"})
		}
		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Yetki yok"})
	}
}

/*
Authorization Middleware Açıklaması:

- Bu middleware, endpoint'lere erişimi belirli roller ile sınırlar.
- Fiber context'indeki kullanıcı rolünü alarak ona göre yetkilendirme kontrolü yapar.
-c.Locals("role") kullanılarak, daha önce eklenmiş olan kullanıcı rolü elde edilir.

- c.Locals("role") kullanılarak, daha önce eklenmiş olan kullanıcı rolü elde edilir.
- Parametre olarak verilen rollerden herhangi biri kullanıcı rolü ile eşleşirse, c.Next() çağrılır ve istek devam eder.
- Eşleşme bulunamazsa, uygun hata mesajı ile HTTP 401 (Rol bulunamadı) veya 403 (Yetki yok) yanıtı döndürülür.

Kullanım:
- Bu middleware, fonksiyonel olarak yetkilendirme kontrolü yapmak istediğiniz endpoint'lerde kullanılabilir.
- Örneğin, Fiber uygulamanızda yetki kontrolü uygulanacak bir route tanımlamasında şu şekilde kullanılabilir:
    app.Get("/protected", Authorize("admin", "manager"), protectedHandler)
- Böylece, sadece "admin" veya "manager" rolüne sahip kullanıcılar bu rota üzerinden erişim sağlayabilir.
*/
