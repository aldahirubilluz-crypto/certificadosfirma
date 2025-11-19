// ============================================
// internal/handlers/certificate_handlers.go
// ============================================
package handlers

import (
	"server/internal/services"
	"github.com/gofiber/fiber/v3"
)

var certService = services.NewCertificateService()

func GetAllCertificates(c fiber.Ctx) error {
	response, err := certService.GetCertificatesGrouped()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.JSON(response)
}
