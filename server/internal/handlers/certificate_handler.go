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

func GetCertificateByThumbprint(c fiber.Ctx) error {
	thumbprint := c.Params("thumbprint")

	cert, err := certService.GetCertificateByThumbprint(thumbprint)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Certificado no encontrado",
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"data":    cert,
		"success": true,
	})
}