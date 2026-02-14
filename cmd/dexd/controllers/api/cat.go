package api

import (
	"crdx.org/dex/db"
	"crdx.org/dex/pkg/types"
	"crdx.org/dex/pkg/util"
	"github.com/gofiber/fiber/v3"
)

func Cat(c fiber.Ctx) error {
	ref := util.ItemRef(c.Params("*"))

	if item, found := db.FindItemByRef(ref, ref); found {
		contentType := item.ResolvedContentType()
		return c.JSON(types.CatResponse{
			Kind:        item.Kind,
			Content:     item.Content(),
			ContentType: contentType.MediaType,
			IsText:      contentType.IsText,
		})
	} else {
		return errRefNotFound(c, ref)
	}
}
