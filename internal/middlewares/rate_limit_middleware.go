package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"
)

func RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Expiration: 30 * time.Second,
		Max:        10,
	})
}
