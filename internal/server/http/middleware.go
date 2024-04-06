package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/michelprogram/htmx-go/internal/models"
	"github.com/michelprogram/htmx-go/internal/pkg"
	"github.com/michelprogram/htmx-go/internal/server/jwt"
	"net/http"
)

func ParseUserBody(c *fiber.Ctx) error {
	var parser models.User

	err := c.BodyParser(&parser)

	if err != nil {
		return err
	}

	c.Locals("parser", parser)

	return c.Next()
}

func ParseUserIdFromJwt(c *fiber.Ctx) error {
	id, err := jwt.ParseID(c)

	if err != nil {
		return err
	}

	c.Locals("id", id)

	return c.Next()
}

func AssetsFiles() fiber.Handler {
	return filesystem.New(filesystem.Config{
		Root:       http.FS(pkg.Dist),
		PathPrefix: "dist/assets",
		Browse:     true,
	})
}
