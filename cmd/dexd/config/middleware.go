package config

import (
	"time"

	"crdx.org/dex/cmd/dexd/env"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func InitMiddleware(app *fiber.App) {
	if !env.Production() {
		app.Use(logger.New())
	}

	app.Use(limiter.New(limiter.Config{Max: 300, Expiration: 60 * time.Second}))
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	app.Use(compress.New())
}
