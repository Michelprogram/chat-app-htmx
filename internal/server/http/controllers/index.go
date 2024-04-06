package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/michelprogram/htmx-go/internal/pkg"
)

func Index(c *fiber.Ctx) error {

	c.Set("Content-Type", "text/html")

	content, err := pkg.Dist.ReadFile("dist/index.html")

	if err != nil {
		return err
	}

	return c.Send(content)
}
