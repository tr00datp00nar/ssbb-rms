package router

import (
	"ssbb-rms/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func RegisterRoutes(app *fiber.App) {
	app.Get("/", adaptor.HTTPHandlerFunc(handlers.Auth))

	app.Get("/auth/{provider}", adaptor.HTTPHandlerFunc(handlers.Auth))
	app.Get("/auth/{provider}/callback", adaptor.HTTPHandlerFunc(handlers.AuthCallback))
	app.Get("/logout/{provider}", adaptor.HTTPHandlerFunc(handlers.Logout))
}
