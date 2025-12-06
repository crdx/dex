package api

import (
	"crdx.org/dex/db"
	"crdx.org/dex/pkg/types"
	"github.com/gofiber/fiber/v3"
	"github.com/lithammer/shortuuid/v3"
	"github.com/samber/lo"
)

func Copy(c fiber.Ctx) error {
	var req types.CopyRequest
	lo.Must0(c.Bind().Body(&req))

	item, found := db.FindItemByRef(req.From, req.From)
	if !found {
		return errRefNotFound(c, req.From)
	}

	if req.To != "" {
		if _, found := db.FindItemByLabel(req.To); found {
			return errDuplicateRef(c, req.To)
		}
		if _, found := db.FindItemByUUID(req.To); found {
			return errDuplicateRef(c, req.To)
		}
	}

	newItem := db.Item{
		Label:       req.To,
		UUID:        shortuuid.New(),
		Kind:        item.Kind,
		BlobHash:    item.BlobHash,
		ContentType: item.ContentType,
	}

	db.CreateItem(&newItem)

	return c.JSON(types.CopyResponse{
		URL: newItem.URL(),
	})
}
