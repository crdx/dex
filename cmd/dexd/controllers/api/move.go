package api

import (
	"crdx.org/dex/db"
	"crdx.org/dex/pkg/types"
	"github.com/gofiber/fiber/v3"
	"github.com/samber/lo"
)

func Move(c fiber.Ctx) error {
	var req types.MoveRequest
	lo.Must0(c.Bind().Body(&req))

	item, found := db.FindItemByRef(req.FromRef, req.FromRef)
	if !found {
		return errRefNotFound(c, req.FromRef)
	}

	if req.ToLabel != "" {
		if _, found := db.FindItemByLabel(req.ToLabel); found {
			return errDuplicateRef(c, req.ToLabel)
		}
		if _, found := db.FindItemByUUID(req.ToLabel); found {
			return errDuplicateRef(c, req.ToLabel)
		}
	}

	item.UpdateLabel(req.ToLabel)

	return c.JSON(types.MoveResponse{
		URL: item.URL(),
	})
}
