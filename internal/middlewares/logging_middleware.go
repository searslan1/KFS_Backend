package middlewares

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// LoggingMiddleware her isteği loglar
func LoggingMiddleware(c *fiber.Ctx) error {
	start := time.Now()
	err := c.Next()
	duration := time.Since(start)

	log.Printf("[%s] %s %s %v", c.Method(), c.Path(), c.IP(), duration)
	return err
}
