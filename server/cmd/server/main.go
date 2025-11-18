// ============================================
// cmd/server/main.go
// ============================================
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"server/internal/config"
	"server/internal/routes"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	cfg := config.Load()

	app := fiber.New(fiber.Config{
		AppName:      "Certificados Firma API",
		ErrorHandler: errorHandler,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
	}))

	routes.RegisterRoutes(app)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		log.Println("Apagando servidor...")
		_ = app.Shutdown()
	}()

	log.Printf("🚀 Servidor ejecutándose en http://localhost%s", cfg.Port)
	if err := app.Listen(cfg.Port); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}

func errorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Error interno del servidor"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"error":   message,
		"success": false,
	})
}