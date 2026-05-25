package api

import (
	"fmt"

	"crdx.org/dex/db"
	"crdx.org/dex/pkg/types"
	"crdx.org/dex/pkg/util"
	"github.com/gofiber/fiber/v3"
)

func Delete(c fiber.Ctx) error {
	ref := util.ItemRef(c.Params("*"))

	if item, found := db.FindItemByRef(ref, ref); found {
		item.Delete()
		return c.JSON(types.DeleteResponse{
			Message: fmt.Sprintf("item %s deleted", ref),
		})
	} else {
		return errRefNotFound(c, ref)
	}
}
