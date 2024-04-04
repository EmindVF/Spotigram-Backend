package server

import (
	"fmt"
	"spotigram/internal/server/config"
	"spotigram/internal/server/controllers"
	"spotigram/internal/server/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type FiberServer struct {
	app *fiber.App
}

func (s *FiberServer) Start() {
	s.app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	s.app.Get("/about", controllers.AboutHandler)

	auth := s.app.Group("/auth")
	auth.Post("/register", controllers.SignUpHandler)
	auth.Post("/login", controllers.SignInHandler)
	auth.Get("/logout", middleware.DeserializeUser, controllers.LogoutHandler)
	auth.Get("/refresh", controllers.RefreshAccessTokenHandler)

	me := s.app.Group("/me")
	me.Get("/info", middleware.DeserializeUser, controllers.MyInfoHandler)

	s.app.All("*", controllers.NotFoundHandler)

	serverUrl := fmt.Sprintf(":%v", config.Cfg.App.Port)
	s.app.Listen(serverUrl)
}
