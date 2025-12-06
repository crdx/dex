package api

import (
	"crypto/sha256"
	"crypto/subtle"
	"errors"

	"crdx.org/dex/cmd/dexd/env"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/keyauth"
)

func InitRoutes(app *fiber.App) {
	api := app.Group("/api", keyauth.New(keyauth.Config{
		Validator: func(c fiber.Ctx, key string) (bool, error) {
			h1 := sha256.Sum256([]byte(env.APIKey()))
			h2 := sha256.Sum256([]byte(key))

			if subtle.ConstantTimeCompare(h1[:], h2[:]) == 1 {
				return true, nil
			}
			return false, errors.New("unauthorised")
		},
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return errUnauthorized(c, err.Error())
		},
	}))

	api.Post("/upload", Upload)
	api.Post("/delete/*", Delete)
	api.Post("/cp", Copy)
	api.Post("/mv", Move)
	api.Post("/set", Set)
	api.Post("/gc", GC)
	api.Post("/expose", Expose)
	api.Post("/urls/delete/:id", DeployURLsDelete)

	api.Get("/cat/*", Cat)
	api.Get("/list", List)
	api.Get("/urls", DeployURLs)
	api.Get("/deployments", Deployments)
	api.Get("/tail", Tail)
}
