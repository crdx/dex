package main

import (
	"embed"
	"log"
	"net/http"
	"os"

	"crdx.org/dex/cmd/dexd/config"
	"crdx.org/dex/cmd/dexd/env"
	"crdx.org/dex/db"
	"crdx.org/dex/pkg/util"
	"crdx.org/duckopt/v2"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/samber/lo"
)

//go:embed views
var views embed.FS

func getUsage() string {
	return `
		Usage:
			$0 [options] [--env PATH]

		Options:
			--env PATH    Read environment file
	`
}

type Opts struct {
	EnvFile string `docopt:"--env"`
}

func main() {
	log.SetFlags(0)
	checkHealth()

	opts := duckopt.MustBind[Opts](getUsage(), "$0")

	initEnvironment(opts.EnvFile)
	initState()

	app := fiber.New(config.GetFiberConfig(views, "views"))
	app.Get("/health", healthcheck.New())

	config.InitMiddleware(app)
	config.InitRoutes(app)

	panic(app.Listen(env.Host() + ":" + env.Port()))
}

func checkHealth() {
	if len(os.Args) != 2 || os.Args[1] != "--health" {
		return
	}

	response, err := http.Get("http://localhost:" + os.Getenv("PORT") + "/health")
	if err != nil || response == nil {
		os.Exit(1)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		os.Exit(1)
	}

	os.Exit(0)
}

func initEnvironment(envFile string) {
	if envFile != "" {
		lo.Must0(env.InitFrom(envFile))
	} else if util.PathExists(".env") {
		lo.Must0(env.InitFrom(".env"))
	} else {
		env.Init()
	}

	if err := env.Validate(); err != nil {
		log.Fatal(err)
	}
}

func initState() *db.Config {
	dbConfig := config.GetDbConfig()
	lo.Must0(db.Init(dbConfig))
	return dbConfig
}
