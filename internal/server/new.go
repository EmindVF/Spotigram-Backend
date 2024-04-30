package server

import (
	serviceAbstractions "spotigram/internal/service/abstractions"

	"github.com/gofiber/fiber/v2"
)

// Return a fiber server instance.
// requestBodyLimit sets the maximum size of a request (in bytes).
func NewFiberServer(requestBodyLimit int) serviceAbstractions.Server {
	return &FiberServer{
		app: fiber.New(fiber.Config{
			BodyLimit: requestBodyLimit,
		}),
	}
}
