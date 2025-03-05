package middlewares

import (
	"strconv"
	"time"

	"KFS_Backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// LoggingMiddleware, istek ve yanıt detaylarını loglayan middleware
func LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// İsteğin başlangıç zamanını kaydet
		start := time.Now()

		// Request ID'yi al
		rid, ok := c.Locals("requestID").(string)
		if !ok {
			rid = "unknown"
		}

		// Kullanıcı bilgilerini al
		userID, _ := c.Locals("userID").(int64)
		userRole, _ := c.Locals("role").(string)

		// İstek detaylarını logla
		utils.Log.Info().
			Str("request_id", rid).
			Int64("user_id", userID).
			Str("role", userRole).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("ip", c.IP()).
			Str("user_agent", c.Get("User-Agent")).
			Msg("İstek başladı")

		// Sonraki middleware'e geç
		err := c.Next()

		// İstek tamamlandıktan sonra bilgileri topla
		duration := time.Since(start)
		statusCode := c.Response().StatusCode()
		respSize := len(c.Response().Body())

		// Yanıt detaylarını logla
		logEvent := utils.Log.Info()
		if err != nil || statusCode >= 400 {
			logEvent = utils.Log.Error().Err(err)
		}

		logEvent.
			Str("request_id", rid).
			Int64("user_id", userID).
			Str("role", userRole).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", statusCode).
			Dur("duration", duration).
			Int("req_size", len(c.Body())).
			Int("resp_size", respSize).
			Str("ip", c.IP()).
			Str("user_agent", c.Get("User-Agent")).
			Msg("İstek tamamlandı")

		// Log yap
		utils.LogRequest(utils.RequestLog{
			Method:      c.Method(),
			Path:        c.Path(),
			IP:          c.IP(),
			Duration:    duration,
			StatusCode:  statusCode,
			UserAgent:   c.Get("User-Agent"),
			ReqBodySize: len(c.Body()),
			RespSize:    respSize,
			Error:       err,
			UserID:      strconv.FormatInt(userID, 10), // int64'ü string'e çevir
			Role:        userRole,
		})

		return err
	}
}

/*
LoggingMiddleware Açıklaması:
- Bu middleware, gelen HTTP isteklerinin detaylı loglarını toplar.
- İstek gövdesi içeriğini loglamaz, sadece boyutunu kaydeder - bu güvenlik ve gizlilik için daha uygundur.
- Kullanıcı kimlik bilgilerini (ID, rol), yanıt bilgilerini ve performans ölçümlerini loglar.
- Zerolog ile yüksek performanslı loglama sağlar.

Kullanım:
- app.Use(LoggingMiddleware())
*/
