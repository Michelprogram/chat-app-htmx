package ws

import (
	"github.com/gofiber/fiber/v2"
	"github.com/michelprogram/htmx-go/internal/database"
	"github.com/michelprogram/htmx-go/internal/repository"
	"github.com/michelprogram/htmx-go/internal/server/ws/sockets"
)

func HandlersWebsocket(app *fiber.App, conn *database.Mongo) {

	//Repositories
	userRepository := repository.NewUserRepository("users", conn)
	messageRepository := repository.NewMessageRepository("messages", conn)

	app.Use("/ws", Upgrade)

	wss := sockets.NewChat(userRepository, messageRepository)

	app.Get("/ws/chatroom", wss.Chats())
}
