package middlewares

import "github.com/gofiber/fiber/v2"

// Örnek domain listesi //Biz buraya kendi domainlerimizi ekleyeceğiz
var allowedDomains = []string{
	"https://example.com",
	"https://sub.example.com",
}

// CorsMiddleware, Cross-Origin Resource Sharing (CORS) politikalarını yöneten bir middleware fonksiyonudur.
// Bu middleware şunları yapar:
// - Önceden tanımlanmış güvenli domainlerden gelen isteklere izin verir
// - İzin verilen HTTP metodlarını (GET,POST,PUT,DELETE,OPTIONS) belirler
// - İzin verilen HTTP başlıklarını tanımlar
// - Credentials (kimlik bilgileri) kullanımına izin verir
// - OPTIONS isteklerini (preflight requests) uygun şekilde yanıtlar
// - İzinsiz domainlerden gelen istekleri güvenli bir şekilde reddeder
//
// Not: Bu middleware, frontend uygulamanızın API'ye güvenli bir şekilde erişebilmesini sağlar.
// Frontend ve backend farklı domainlerde çalıştığında, tarayıcı güvenlik nedeniyle istekleri
// engeller. Bu middleware, güvenilir frontend domainlerinden gelen isteklere izin vererek
// bu sorunu çözer.
func CorsMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		origin := c.Get("Origin")
		allowed := false

		for _, d := range allowedDomains {
			if d == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Set("Access-Control-Allow-Origin", origin)
		} else {
			c.Set("Access-Control-Allow-Origin", "null") // İzinli değilse
		}

		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Set("Access-Control-Allow-Credentials", "true")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	}
}
