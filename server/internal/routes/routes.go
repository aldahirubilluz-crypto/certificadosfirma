// ============================================
// internal/routes/routes.go
// ============================================
package routes

import (
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(app *fiber.App) {
	RegisterCertificateRoutes(app)
}
