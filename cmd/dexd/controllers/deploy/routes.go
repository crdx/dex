package deploy

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"crdx.org/dex/cmd/dexd/controllers/api"
	"crdx.org/dex/cmd/dexd/env"
	"crdx.org/dex/db"
	"crdx.org/dex/pkg/mail"
	"crdx.org/dex/pkg/util"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

//go:embed usage.md
var usageTemplate string

func InitRoutes(app *fiber.App) {
	deployLimiter := limiter.New(limiter.Config{
		Max:        5,
		Expiration: 60 * time.Second,
		KeyGenerator: func(c fiber.Ctx) string {
			return c.Params("token")
		},
		LimitReached: func(c fiber.Ctx) error {
			return c.SendStatus(http.StatusTooManyRequests)
		},
	})

	app.Get("/deploy/:token", Usage)
	app.Post("/deploy/:token", deployLimiter, Deploy)
}

type usageData struct {
	PublicURL string
	DeployURL string
	ExpiresAt string
}

func Usage(c fiber.Ctx) error {
	token, err := validateToken(c)
	if err != nil {
		return err
	}

	deployURL, _ := url.JoinPath(env.BaseURL(), "deploy", token.Token)

	item, found := db.FindItem(token.ItemID)
	if !found {
		return api.FailureResponse(c, http.StatusNotFound, "item not found")
	}

	data := usageData{
		PublicURL: item.URL(),
		DeployURL: deployURL,
	}

	if token.ExpiresAt.Valid {
		data.ExpiresAt = token.ExpiresAt.V.Format(time.RFC3339)
	}

	tmpl := template.Must(template.New("usage").Parse(usageTemplate))
	var buf bytes.Buffer
	_ = tmpl.Execute(&buf, data)

	c.Type("md", "utf-8")
	return c.Send(buf.Bytes())
}

func Deploy(c fiber.Ctx) error {
	token, err := validateToken(c)
	if err != nil {
		return err
	}

	item, found := db.FindItem(token.ItemID)
	if !found {
		return api.FailureResponse(c, http.StatusNotFound, "item not found")
	}

	note := c.FormValue("note")
	if note == "" {
		return api.FailureResponse(c, http.StatusBadRequest, "note is required")
	}

	deployer := c.FormValue("deployer")
	if deployer == "" {
		return api.FailureResponse(c, http.StatusBadRequest, "deployer is required")
	}

	fileHeader, err := c.FormFile("content")
	if err != nil {
		return api.FailureResponse(c, http.StatusBadRequest, "content file is required")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return api.FailureResponse(c, http.StatusInternalServerError, "failed to open file")
	}
	defer func() { _ = file.Close() }()

	content, err := io.ReadAll(file)
	if err != nil {
		return api.FailureResponse(c, http.StatusInternalServerError, "failed to read file")
	}

	blob := db.FindOrCreateBlob(content)

	item.UpdateBlobHash(blob.Hash)

	db.CreateDeployment(&db.Deployment{
		TokenID:   token.ID,
		Note:      note,
		Deployer:  util.Truncate(deployer, 20),
		IPAddress: util.Truncate(c.IP(), 45),
		UserAgent: util.Truncate(c.Get("user-agent"), 200),
	})

	go notify(item.Label, note, deployer, c.IP(), c.Get("user-agent"), item.URL())

	output := fmt.Sprintf("URL: %s\nNote: %s\n", item.URL(), note)
	if token.ExpiresAt.Valid {
		output += fmt.Sprintf("Info: You can continue to deploy to this endpoint for another %s. The public URL will always be accessible.\n", util.FormatDuration(time.Until(token.ExpiresAt.V), true, 2, ""))
	}
	return c.Type("txt").SendString(output)
}

func validateToken(c fiber.Ctx) (*db.Token, error) {
	tokenStr := c.Params("token")

	token, found := db.FindTokenByToken(tokenStr)
	if !found {
		return nil, api.FailureResponse(c, http.StatusNotFound, "deploy url not found")
	}

	if token.ExpiresAt.Valid && token.ExpiresAt.V.Before(time.Now()) {
		return nil, api.FailureResponse(c, http.StatusGone, "deploy url expired")
	}

	return token, nil
}

func notify(label, note, ipAddress, userAgent, publicURL string) {
	subject := fmt.Sprintf("deployed %s", label)

	body := fmt.Sprintf(`%s was just deployed!

%s

Note:   %s
IP:     %s
Device: %s
`,
		label,
		publicURL,
		deployer,
		note,
		ipAddress,
		userAgent,
	)

	mail.Send(subject, body)
}
