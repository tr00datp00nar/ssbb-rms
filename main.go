package main

import (
	"ssbb-rms/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	router.RegisterRoutes(app)

	app.Listen(":3000")
}
