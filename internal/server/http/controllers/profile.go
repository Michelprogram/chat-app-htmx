package controllers

import (
	"bytes"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/michelprogram/htmx-go/internal/models"
	"github.com/michelprogram/htmx-go/internal/repository"
	"github.com/michelprogram/htmx-go/web/htmx/account"
)

type UserController struct {
	repository.IUserRepository
}

func NewUserController(user repository.IUserRepository) *UserController {
	return &UserController{
		IUserRepository: user,
	}
}

func (uc UserController) Update(c *fiber.Ctx) error {

	var parser models.User

	var buffer bytes.Buffer

	id := c.Locals("id").(string)

	err := c.BodyParser(&parser)

	if err != nil {
		return err
	}

	err = uc.FindAndUpdate(id, parser)
	if err != nil {
		return err
	}

	_ = account.Index(parser.Username).Render(context.TODO(), &buffer)

	return c.Send(buffer.Bytes())
}

func (uc UserController) Edit(c *fiber.Ctx) error {

	var buffer bytes.Buffer

	id := c.Locals("id").(string)

	user := uc.FindById(id)

	_ = account.EditAccount(*user).Render(context.TODO(), &buffer)

	return c.Send(buffer.Bytes())
}
