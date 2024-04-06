package pkg

import (
	"bytes"
	"embed"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

//go:embed dist/*
var Dist embed.FS

func ParseTemplToBuffer(c *fiber.Ctx, template templ.Component) ([]byte, error) {
	var buffer bytes.Buffer

	c.Set("Content-type", "text/html")

	err := template.Render(c.Context(), &buffer)

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
