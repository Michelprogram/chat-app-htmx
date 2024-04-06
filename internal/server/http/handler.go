package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/michelprogram/htmx-go/internal/database"
	"github.com/michelprogram/htmx-go/internal/repository"
	. "github.com/michelprogram/htmx-go/internal/server/http/controllers"
)

func HandlersHttp(app *fiber.App, conn *database.Mongo) {

	//Repositories
	userRepository := repository.NewUserRepository("users", conn)
	messageRepository := repository.NewMessageRepository("messages", conn)
	emoteRepository := repository.NewEmoteRepository("messages", conn)

	//Base & assets
	app.Use("/assets", AssetsFiles())

	app.Get("/", Index)

	//Aut
	authController := NewAuthController(userRepository)
	auth := app.Group("/auth")
	auth.Use(ParseUserBody)

	auth.Post("/login", authController.Login)
	auth.Post("/register", authController.Register)

	//Htmx
	info := app.Group("/htmx")

	info.Use(ParseUserIdFromJwt)

	userController := NewUserController(userRepository)
	user := info.Group("/user")

	user.Get("/", userController.Edit)
	user.Post("/", userController.Update)

	messageController := NewMessageController(messageRepository, userRepository)
	message := info.Group("/message")

	message.Get("/last/:id", messageController.LastMessages)
	message.Delete("/:id", messageController.Remove)

	emoteController := NewEmoteController(emoteRepository, userRepository)
	emote := info.Group("/emote")

	emote.Post("/:messageid", emoteController.Add)
	emote.Put("/:emoteid", emoteController.Counter)

}
