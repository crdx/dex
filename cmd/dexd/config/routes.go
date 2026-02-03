package config

import (
	"crdx.org/dex/cmd/dexd/controllers/api"
	"crdx.org/dex/cmd/dexd/controllers/deploy"
	"crdx.org/dex/cmd/dexd/controllers/index"
	"github.com/gofiber/fiber/v3"
)

func InitRoutes(app *fiber.App) {
	api.InitRoutes(app)
	deploy.InitRoutes(app)

	app.Get("/robots.txt", func(c fiber.Ctx) error {
		return c.SendString("User-agent: *\nDisallow: /")
	})

	index.InitRoutes(app)
}
