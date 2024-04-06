package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/michelprogram/htmx-go/internal/database"
	"github.com/michelprogram/htmx-go/internal/server/http"
	"github.com/michelprogram/htmx-go/internal/server/room"
	"github.com/michelprogram/htmx-go/internal/server/ws"
)

type WebServer struct {
	Port  string
	Debug bool

	App  *fiber.App
	Conn *database.Mongo
}

func NewServer(port string, debug bool, conn *database.Mongo) WebServer {

	app := fiber.New()

	chat := room.GetInstance()

	if debug {
		app.Use(logger.New())
	}

	app.Use(cors.New())

	//Html controllers
	http.HandlersHttp(app, conn)

	//Websocket controllers
	ws.HandlersWebsocket(app, conn)

	go chat.HandleEvent()

	return WebServer{
		Port:  fmt.Sprintf(":%s", port),
		Debug: debug,
		App:   app,
		Conn:  conn,
	}
}

func (s WebServer) Start() error {

	err := s.App.Listen(s.Port)

	if err != nil {
		return err
	}

	return nil
}
