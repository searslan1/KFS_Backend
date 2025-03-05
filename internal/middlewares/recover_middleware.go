package middlewares

import (
	"runtime/debug"

	"KFS_Backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// RecoveryMiddleware, panic durumlarını yakalayıp işleyerek uygulamanın çökmesini engeller
func RecoveryMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// recover fonksiyonunu defer ile çağırarak panic durumunu yakalıyoruz
		defer func() {
			if r := recover(); r != nil {
				// Stack trace bilgisini al
				stackTrace := debug.Stack()

				// Hatayı logla
				utils.Log.Error().
					Interface("panic", r).
					Str("stack_trace", string(stackTrace)).
					Str("path", c.Path()).
					Str("method", c.Method()).
					Str("ip", c.IP()).
					Msg("Recovered from panic")

				// Kullanıcıya kontrollü bir hata yanıtı dön
				err := c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Sunucu hatası oluştu",
					"message": "İşleminiz gerçekleştirilemedi, lütfen daha sonra tekrar deneyin",
				})

				if err != nil {
					// JSON yanıtı bile dönemiyorsak, basit bir metin yanıtı dene
					_ = c.Status(500).SendString("Internal Server Error")
				}
			}
		}()

		// Sonraki middleware'e geç
		return c.Next()
	}
}

/*
RecoveryMiddleware Kullanımı:
- Bu middleware'i her zaman ilk sırada eklemek önemlidir.
- Böylece diğer middleware'lerde veya rota işleyicilerinde oluşabilecek panic'leri yakalayabilir.

app := fiber.New()
app.Use(middlewares.RecoveryMiddleware())  // ←  Daima ilk olarak ekleyin
app.Use(middlewares.AuthMiddleware())
app.Use(middlewares.LoggingMiddleware())
// ...diğer middleware'ler
*/
