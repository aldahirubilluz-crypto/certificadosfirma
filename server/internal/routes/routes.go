// ============================================
// internal/routes/routes.go
// ============================================
package routes

import (
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(app *fiber.App) {
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	api := app.Group("/api/v1")
	RegisterCertificateRoutes(api)
}