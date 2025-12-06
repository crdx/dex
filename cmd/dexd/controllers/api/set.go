package api

import (
	"slices"

	"crdx.org/dex/db"
	"crdx.org/dex/pkg/types"
	"github.com/gofiber/fiber/v3"
	"github.com/samber/lo"
)

func Set(c fiber.Ctx) error {
	var req types.SetRequest
	lo.Must0(c.Bind().Body(&req))

	item, found := db.FindItemByRef(req.Ref, req.Ref)
	if !found {
		return errRefNotFound(c, req.Ref)
	}

	var modified bool

	switch req.Key {
	case types.KeyContentType:
		if item.ContentType != req.Value {
			modified = true
		}
		item.UpdateContentType(req.Value)
	case types.KeyKind:
		if !slices.Contains(types.Kinds, req.Value) {
			return errInvalidKind(c, req.Value)
		}

		if item.Kind != req.Value {
			modified = true
		}

		item.UpdateKind(req.Value)
	}

	return c.JSON(types.SetResponse{Modified: modified})
}
