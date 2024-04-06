package ws

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/michelprogram/htmx-go/internal/server/jwt"
	"log"
)

func Upgrade(c *fiber.Ctx) error {

	if websocket.IsWebSocketUpgrade(c) {

		id, err := jwt.ParseID(c)

		if err != nil {
			log.Println(err.Error())
			return err
		}

		c.Locals("id", id)

		c.Locals("allowed", true)

		return c.Next()
	}

	return fiber.ErrUpgradeRequired
}
