package jwt

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/michelprogram/htmx-go/internal/models"
	"os"
	"time"
)

var secret = os.Getenv("JWT_SECRET")

func Generate(user models.User) (string, error) {

	claims := jwt.MapClaims{
		"id":  user.ID.Hex(),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ParseID(c *fiber.Ctx) (string, error) {

	//TODO: prendre dans le header
	token := c.Cookies("token")

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return "", fiber.NewError(401, "invalid jwt")
	}

	return fmt.Sprintf("%s", claims["id"]), nil
}
