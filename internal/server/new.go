package server

import (
	serviceAbstractions "spotigram/internal/service/abstractions"

	"github.com/gofiber/fiber/v2"
)

func NewFiberServer() serviceAbstractions.Server {
	return &FiberServer{
		app: fiber.New(),
	}
}
