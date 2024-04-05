package server

import (
	"fmt"
	"os"
	"os/signal"
	"spotigram/internal/server/config"
	"spotigram/internal/server/controllers"
	"spotigram/internal/server/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type FiberServer struct {
	app *fiber.App
}

// Starts fiber server.
func (s *FiberServer) Start() {
	// Gracefull shutdown.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		_ = s.app.Shutdown()
	}()

	// Logger
	s.app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	s.app.Get("/about", controllers.AboutHandler)

	auth := s.app.Group("/auth")
	auth.Post("/register", controllers.SignUpHandler)
	auth.Post("/login", controllers.SignInHandler)
	auth.Get("/logout", middleware.DeserializeTokenHandler, controllers.LogoutHandler)
	auth.Get("/refresh", controllers.RefreshAccessTokenHandler)

	me := s.app.Group("/me")
	me.Get("/info", middleware.DeserializeTokenHandler, controllers.MyInfoHandler)

	s.app.All("*", controllers.NotFoundHandler)

	// Listening
	serverUrl := fmt.Sprintf(":%v", config.Cfg.App.Port)
	s.app.Listen(serverUrl)
}
