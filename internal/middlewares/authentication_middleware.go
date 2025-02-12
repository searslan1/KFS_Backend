package middlewares

import (
    "strings"

    "KFS_Backend/internal/utils"
    "github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token bulunamadı"})
        }

        tokenParts := strings.Split(authHeader, "Bearer ")
        if len(tokenParts) != 2 {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Geçersiz token formatı"})
        }

        tokenString := tokenParts[1]

        claims, err := utils.ValidateToken(tokenString)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Geçersiz veya süresi dolmuş token"})
        }

        c.Locals("userID", claims.UserID)
        c.Locals("role", claims.Role)

        return c.Next()
    }
}

/*
Authentication Middleware Açıklaması:
- Bu middleware, gelen isteklerin kimlik doğrulamasını yapar.
- "Authorization" header'ını kontrol eder ve Bearer şemasına uygun bir token arar.
- Token bulunamazsa veya formatı geçersizse, HTTP 401 (Unauthorized) hatası döndürür.
- Token geçerliyse, utils.ValidateToken fonksiyonu ile doğrulanır.
- Doğrulama başarılı olursa, token içerisindeki kullanıcı ID'si ve rol bilgisi c.Locals'e eklenir.
- Son olarak, c.Next() çağrılarak isteğin sonraki handler'a geçmesi sağlanır.

Örnek Kullanım:
- Bu middleware, kimlik doğrulaması gerektiren endpoint'lerde kullanılabilir.
- Örneğin, Fiber uygulamanızda kimlik doğrulaması yapılacak bir route tanımlamasında şu şekilde kullanılabilir:
    app.Get("/profile", AuthMiddleware(), profileHandler)
- Bu durumda, /profile endpoint'ine gelen istekler önce AuthMiddleware'den geçecek ve geçerli bir token varsa profileHandler çalışacaktır.
*/