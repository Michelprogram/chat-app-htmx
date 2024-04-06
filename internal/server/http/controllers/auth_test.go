package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/michelprogram/htmx-go/internal/database"
	"github.com/michelprogram/htmx-go/internal/models"
	"github.com/michelprogram/htmx-go/internal/repository"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"
	"os"
	"runtime/pprof"
	"testing"
)

func init() {
	_ = os.Setenv("URL_MONGODB", "mongodb://user:password@localhost:27018/?authMechanism=SCRAM-SHA-1")

	conn, err := database.NewMongo(false, "21_chat")

	if err != nil {
		panic(err)
	}

	_ = conn.Database.Drop(context.TODO())
}

func setupAuthController() *AuthController {

	conn, err := database.NewMongo(false, "21_chat")

	if err != nil {
		panic(err)
	}

	userRepository := repository.NewUserRepository("users", conn)

	return NewAuthController(userRepository)
}

func BenchmarkAuthController_Register(b *testing.B) {

	f, err := os.Create("register_cpu.prof")
	if err != nil {
		b.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	authController := setupAuthController()

	app := fiber.New()

	for i := 0; i < b.N; i++ {
		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
		ctx.Locals("parser", models.User{Username: fmt.Sprintf("user-%d", i), Password: fmt.Sprintf("password-%d", i)})
		_ = authController.Register(ctx)
	}
}

func BenchmarkAuthController_Login(b *testing.B) {
	f, err := os.Create("login_cpu.prof")
	if err != nil {
		b.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	authController := setupAuthController()

	app := fiber.New()

	for i := 0; i < b.N; i++ {
		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
		ctx.Locals("parser", models.User{Username: fmt.Sprintf("user-%d", i), Password: fmt.Sprintf("password-%d", i)})
		_ = authController.Login(ctx)
	}
}
