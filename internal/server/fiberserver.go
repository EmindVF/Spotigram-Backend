package server

import (
	"fmt"
	"os"
	"os/signal"
	"spotigram/internal/server/config"
	"spotigram/internal/server/controllers"
	"spotigram/internal/server/middleware"
	ws "spotigram/internal/server/websocket"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
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
	auth.Get("/logout", middleware.DeserializeTokenHandler, controllers.LogoutHandler)
	auth.Get("/refresh", controllers.RefreshAccessTokenHandler)
	auth.Post("/register", controllers.SignUpHandler)
	auth.Post("/login", controllers.SignInHandler)

	me := s.app.Group("/me")
	me.Get("/friends", middleware.DeserializeTokenHandler, controllers.FriendsHandler)
	me.Get("/friend-requests-sent", middleware.DeserializeTokenHandler, controllers.FriendRequestSentHandler)
	me.Get("/friend-requests-received", middleware.DeserializeTokenHandler, controllers.FriendRequestReceivedHandler)
	me.Get("/info", middleware.DeserializeTokenHandler, controllers.MyInfoHandler)
	me.Get("/public-key", middleware.DeserializeTokenHandler, controllers.MyPublicKeyHandler)
	me.Get("/picture", middleware.DeserializeTokenHandler, controllers.MyPictureHandler)
	me.Post("/change-name", middleware.DeserializeTokenHandler, controllers.ChangeNameHandler)
	me.Post("/change-password", middleware.DeserializeTokenHandler, controllers.ChangePasswordHandler)
	me.Post("/change-public-key", middleware.DeserializeTokenHandler, controllers.ChangePublicKeyHandler)
	me.Post("/change-picture", middleware.DeserializeTokenHandler, controllers.ChangePictureHandler)

	user := s.app.Group("/user")
	user.Get("/all", middleware.DeserializeTokenHandler, controllers.UsersHandler)
	user.Get("/info", middleware.DeserializeTokenHandler, controllers.UserInfoHandler)
	user.Get("/public-key", middleware.DeserializeTokenHandler, controllers.UserPublicKeyHandler)
	user.Get("/picture", middleware.DeserializeTokenHandler, controllers.UserPictureHandler)

	playlist := s.app.Group("/playlist")
	playlist.Get("/all", middleware.DeserializeTokenHandler, controllers.PlaylistsHandler)
	playlist.Get("/songs", middleware.DeserializeTokenHandler, controllers.PlaylistSongsHandler)
	playlist.Post("/create", middleware.DeserializeTokenHandler, controllers.AddPlaylistHandler)
	playlist.Post("/rename", middleware.DeserializeTokenHandler, controllers.RenamePlaylistHandler)
	playlist.Delete("/delete", middleware.DeserializeTokenHandler, controllers.DeletePlaylistHandler)
	playlist.Post("/add-song", middleware.DeserializeTokenHandler, controllers.AddPlaylistSongHandler)
	playlist.Delete("/delete-song", middleware.DeserializeTokenHandler, controllers.DeletePlaylistSongHandler)

	song := s.app.Group("/song")
	song.Get("/all", middleware.DeserializeTokenHandler, controllers.SongsHandler)
	song.Get("/info", middleware.DeserializeTokenHandler, controllers.SongInfoHandler)
	song.Post("/rename", middleware.DeserializeTokenHandler, controllers.RenameSongHandler)
	song.Get("/picture", middleware.DeserializeTokenHandler, controllers.SongPictureHandler)
	song.Get("/download", middleware.DeserializeTokenHandler, controllers.DownloadSongHandler)
	song.Delete("/delete", middleware.DeserializeTokenHandler, controllers.DeleteSongHandler)
	song.Get("/stream/:filename", middleware.DeserializeTokenHandler, controllers.GetSongChunk)
	song.Post("/upload/:songname", middleware.DeserializeTokenHandler, controllers.UploadSongHandler)

	chat := s.app.Group("/chat")
	chat.Get("/messages", middleware.DeserializeTokenHandler, controllers.ChatMessagesHandler)
	chat.Get("/unread-messages", middleware.DeserializeTokenHandler, controllers.ChatUnreadMessagesHandler)

	go ws.RunChatHub()
	s.app.Get("/connect", middleware.DeserializeTokenHandler, ws.WebsocketChatUpgradeHandler,
		websocket.New(ws.WebsocketChatLoop))

	s.app.All("*", controllers.NotFoundHandler)

	// Listening
	serverUrl := fmt.Sprintf(":%v", config.Cfg.App.Port)
	s.app.Listen(serverUrl)
}
