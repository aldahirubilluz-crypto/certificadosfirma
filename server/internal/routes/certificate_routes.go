// ============================================
// internal/routes/certificate_routes.go
// ============================================
package routes

import (
	"server/internal/handlers"
	"github.com/gofiber/fiber/v3"
)

func RegisterCertificateRoutes(router fiber.Router) {
	certs := router.Group("certificates")
	
	certs.Get("/", handlers.GetAllCertificates)
}