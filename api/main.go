package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}

	app := fiber.New()

	app.Use(logger.New())
	setUpRoutes(app)

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))

}

func setUpRoutes(app *fiber.App) {
	// app.Get("/:url", func(ctx *app) {})
	// app.Post("/api/v1/")

}
