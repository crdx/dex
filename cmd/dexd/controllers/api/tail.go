package api

import (
	"fmt"
	"net/url"
	"slices"

	"crdx.org/dex/cmd/dexd/env"
	"crdx.org/dex/db"
	"crdx.org/dex/pkg/types"
	"crdx.org/dex/pkg/util"
	"github.com/gofiber/fiber/v3"
)

func Tail(c fiber.Ctx) error {
	since := int64(fiber.Query[int](c, "since"))
	verbose := fiber.Query[bool](c, "verbose")

	allItems := db.MapBy[int64]("ID", db.FindItemsUnscoped())
	itemIDs := db.Pluck[int64](db.FindItems(), "ID")

	var logs []types.TailResponseItem
	for _, log := range db.FindLogsSince(since) {
		item := allItems[log.ItemID]
		ref := item.Ref()
		requestPath := util.ItemRequestPath(ref)

		// Root gets a lot of bot traffic, so omit it by default.
		if requestPath == "/" && !verbose {
			continue
		}

		// Parse the request so we can ignore everything except the path when figuring out whether
		// this request was for an item that has since been renamed.
		loggedPath, err := url.Parse(log.Request)
		if err != nil {
			panic(fmt.Sprintf("url.Parse(%q) failed: %s", log.Request, err))
		}

		logs = append(logs, types.TailResponseItem{
			ID:        log.ID,
			BaseURL:   env.BaseURL(),
			CreatedAt: log.CreatedAt,
			Method:    log.Method,
			Request:   log.Request,
			IPAddress: log.IPAddress,
			UserAgent: log.UserAgent,
			Renamed:   requestPath != loggedPath.Path,
			Ref:       ref,
			Exists:    slices.Contains(itemIDs, log.ItemID),
		})
	}

	return c.JSON(types.TailResponse{
		Logs: logs,
	})
}
