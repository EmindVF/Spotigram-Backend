package server

import (
	"database/sql"
	"fmt"

	"spotigram/config"

	"github.com/gofiber/fiber/v2"
)

type fiberServer struct {
	app *fiber.App
	db  *sql.DB
	cfg *config.Config
}

func NewFiberServer(cfg *config.Config, db *sql.DB) Server {
	return &fiberServer{
		app: fiber.New(),
		db:  db,
		cfg: cfg,
	}
}

func (s *fiberServer) Start() {
	s.app.Get("/about", func(ctx *fiber.Ctx) error {
		ctx.SendString("You have reached a test version of Spotigram!")
		return nil
	})
	serverUrl := fmt.Sprintf(":%v", s.cfg.App.Port)
	s.app.Listen(serverUrl)
}
