package server

import (
	serviceAbstractions "spotigram/internal/service/abstractions"

	"github.com/gofiber/fiber/v2"
)

// Return a fiber server instance.
func NewFiberServer() serviceAbstractions.Server {
	return &FiberServer{
		app: fiber.New(),
	}
}
