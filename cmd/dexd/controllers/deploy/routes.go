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

	"crdx.org/dex/cmd/dexd/env"
	"crdx.org/dex/db"
	"crdx.org/dex/pkg/mail"
	"crdx.org/dex/pkg/util"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

//go:embed usage.md
var usageTemplate string

const errItemNotFound = "Item not found. The item associated with this deploy URL may have been deleted."

func InitRoutes(app *fiber.App) {
	deployLimiter := limiter.New(limiter.Config{
		Max:        5,
		Expiration: 60 * time.Second,
		KeyGenerator: func(c fiber.Ctx) string {
			return c.Params("token")
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
	token, ok, err := validateToken(c)
	if !ok || err != nil {
		return err
	}

	deployURL, _ := url.JoinPath(env.BaseURL(), "deploy", token.Token)

	item, found := db.FindItem(token.ItemID)
	if !found {
		return textError(c, http.StatusNotFound, errItemNotFound)
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

	c.Type("md")
	return c.Send(buf.Bytes())
}

func Deploy(c fiber.Ctx) error {
	token, ok, err := validateToken(c)

	if !ok || err != nil {
		return err
	}

	item, found := db.FindItem(token.ItemID)
	if !found {
		return textError(c, http.StatusNotFound, errItemNotFound)
	}

	change := c.FormValue("change")
	if change == "" {
		return textError(c, http.StatusBadRequest, "Missing required form field 'change'. Provide a short description of what changed in this deployment.")
	}
	if len(change) > 80 {
		return textError(c, http.StatusBadRequest, "Field 'change' exceeds 80 characters. Provide a shorter description.")
	}

	deployer := c.FormValue("deployer")
	if deployer == "" {
		return textError(c, http.StatusBadRequest, "Missing required form field 'deployer'. Provide the name or identifier of who is deploying.")
	}

	fileHeader, err := c.FormFile("content")
	if err != nil {
		return textError(c, http.StatusBadRequest, "Missing required form field 'content'. Provide the file to deploy as multipart form data.")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return textError(c, http.StatusInternalServerError, "Failed to open uploaded file. Try again or check the file is valid.")
	}
	defer func() { _ = file.Close() }()

	content, err := io.ReadAll(file)
	if err != nil {
		return textError(c, http.StatusInternalServerError, "Failed to read uploaded file. Try again or check the file is valid.")
	}

	blob := db.FindOrCreateBlob(content)

	item.UpdateBlobHash(blob.Hash)

	util.AssertPublicIP(c.IP())

	db.CreateDeployment(&db.Deployment{
		TokenID:   token.ID,
		Note:      change,
		Deployer:  util.Truncate(deployer, 20),
		IPAddress: util.Truncate(c.IP(), 45),
		UserAgent: util.Truncate(c.Get("user-agent"), 200),
	})

	go notify(item.Label, change, deployer, c.IP(), c.Get("user-agent"), item.URL())

	output := fmt.Sprintf("Congratulations! You have successfully deployed to %s\n", item.URL())
	output += fmt.Sprintf("Change: %s\n", change)
	if token.ExpiresAt.Valid {
		output += fmt.Sprintf("Info: You can continue to deploy to this endpoint for another %s. The public URL will always be accessible.\n", util.FormatDuration(time.Until(token.ExpiresAt.V), true, 2, ""))
	}
	return c.Type("txt").SendString(output)
}

func validateToken(c fiber.Ctx) (*db.Token, bool, error) {
	value := c.Params("token")

	token, found := db.FindTokenByToken(value)
	if !found {
		return nil, false, textError(c, http.StatusNotFound, "Deploy URL not found. Check the URL is correct and has not been revoked.")
	}

	if token.ExpiresAt.Valid && token.ExpiresAt.V.Before(time.Now()) {
		return nil, false, textError(c, http.StatusGone, "Deploy URL has expired. Request a new deploy URL to continue deploying.")
	}

	return token, true, nil
}

func textError(c fiber.Ctx, status int, message string) error {
	c.Status(status)
	return c.Type("txt").SendString(message + "\n")
}

func notify(label, change, deployer, ipAddress, userAgent, publicURL string) {
	subject := fmt.Sprintf("%s deployed %s", deployer, label)

	body := fmt.Sprintf(`%s just deployed %s to %s

Change: %s
IP:     %s
Device: %s
`,
		deployer,
		label,
		publicURL,
		change,
		ipAddress,
		userAgent,
	)

	mail.Send(subject, body)
}
