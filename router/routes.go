package router

import (
	"ssbb-rms/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func RegisterRoutes(app *fiber.App) {
	app.Get("/", adaptor.HTTPHandlerFunc(auth.Auth))

	app.Get("/auth/{provider}", adaptor.HTTPHandlerFunc(auth.Auth))
	app.Get("/auth/{provider}/callback", adaptor.HTTPHandlerFunc(auth.AuthCallback))
	app.Get("/logout/{provider}", adaptor.HTTPHandlerFunc(auth.Logout))
}
