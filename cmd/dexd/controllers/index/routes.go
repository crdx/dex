package index

import (
	"crdx.org/dex/db"
	"crdx.org/dex/pkg/types"
	"crdx.org/dex/pkg/util"
	"github.com/gofiber/fiber/v3"
)

func InitRoutes(app *fiber.App) {
	app.Get("*", Get)
}

func Get(c fiber.Ctx) error {
	ref := util.ItemRef(c.Params("*"))

	item, found := db.FindItemByRef(ref, ref)
	if !found {
		return c.SendStatus(404)
	}

	if c.Get("if-none-match") == item.BlobHash {
		return c.SendStatus(fiber.StatusNotModified)
	}

	db.IncrementHits(item.ID)
	item.UpdateLastHitAt(db.Now())

	util.AssertPublicIP(c.IP())

	db.CreateLog(&db.Log{
		ItemID:    item.ID,
		Method:    util.Truncate(c.Method(), 10),
		Request:   util.Truncate(c.OriginalURL(), 200),
		IPAddress: util.Truncate(c.IP(), 15),
		UserAgent: util.Truncate(c.Get("user-agent"), 200),
	})

	content := item.Content()

	if item.Kind == types.KindRedir {
		if c.Query("raw") == "1" {
			return c.Send(content)
		} else {
			return c.Redirect().To(string(content))
		}
	}

	c.Set("etag", item.BlobHash)

	if ref != "/" {
		c.Set("access-control-allow-origin", "*")
	}

	if item.Kind == types.KindFile {
		c.Attachment(item.Base())
	}

	err := c.Send(content)

	// Content-Type is overridden somehow if it does not come after c.Send.
	if header := item.ResolvedContentType().ContentTypeHeader(); header != "" {
		c.Set("content-type", header)
	}

	return err
}
