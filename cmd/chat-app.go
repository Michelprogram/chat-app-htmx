package main

import (
	"flag"
	"github.com/joho/godotenv"
	"github.com/michelprogram/htmx-go/internal/database"
	_ "github.com/michelprogram/htmx-go/internal/database"
	"github.com/michelprogram/htmx-go/internal/server"
	"log"
	"os"
)

var debug bool

func main() {

	_ = godotenv.Load(".env")

	flag.BoolVar(&debug, "debug", false, "Debug mode")

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatalf("Error there is no PORT variable set in environnement.")
	}

	conn, err := database.NewMongo(debug, "21_chat")

	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	app := server.NewServer(port, debug, conn)

	err = app.Start()

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

}
