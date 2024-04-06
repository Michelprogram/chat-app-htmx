package controllers

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/michelprogram/htmx-go/internal/models"
	"github.com/michelprogram/htmx-go/internal/pkg"
	"github.com/michelprogram/htmx-go/internal/repository"
	"github.com/michelprogram/htmx-go/internal/server/jwt"
	"github.com/michelprogram/htmx-go/web/htmx/chat"
	"github.com/michelprogram/htmx-go/web/htmx/errors"
)

type AuthController struct {
	repository.IUserRepository
}

func NewAuthController(user repository.IUserRepository) *AuthController {
	return &AuthController{
		IUserRepository: user,
	}
}

func (auth AuthController) Register(c *fiber.Ctx) error {

	parser := c.Locals("parser").(models.User)

	user := auth.FindByUsername(parser.Username)

	if user != nil {

		c.Status(302)

		tmpl, _ := pkg.ParseTemplToBuffer(c, errors.UsernameAlreadyExists(user.Username))

		return c.Send(tmpl)
	}

	user, err := auth.Insert(parser.Username, parser.Password)

	//TODO: replace with redirection to /auth/login
	if err != nil {
		return err
	}

	token, err := jwt.Generate(*user)

	if err != nil {
		return err
	}

	err = auth.UpdateToken(user.ID, token)

	tmpl, err := pkg.ParseTemplToBuffer(c, chat.Websocket(user))

	if err != nil {
		return err
	}

	data, err := json.Marshal(map[string]string{"jwt": token})

	c.Set("HX-Trigger", string(data))

	return c.Send(tmpl)

}

func (auth AuthController) Login(c *fiber.Ctx) error {

	parser := c.Locals("parser").(models.User)

	user := auth.Find(parser.Username, parser.Password)

	if user == nil {
		return c.SendStatus(404)
	}

	token, err := jwt.Generate(*user)

	if err != nil {
		return err
	}

	err = auth.UpdateToken(user.ID, token)

	tmpl, err := pkg.ParseTemplToBuffer(c, chat.Websocket(user))

	if err != nil {
		return err
	}

	data, err := json.Marshal(map[string]string{"jwt": token})

	c.Set("HX-Trigger", string(data))

	return c.Send(tmpl)
}
