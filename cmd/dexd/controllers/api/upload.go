package api

import (
	"crdx.org/dex/db"
	"crdx.org/dex/pkg/types"
	"crdx.org/dex/pkg/util"
	"github.com/gofiber/fiber/v3"
	"github.com/lithammer/shortuuid/v3"
	"github.com/samber/lo"
)

func Upload(c fiber.Ctx) error {
	var req types.UploadRequest
	lo.Must0(c.Bind().Body(&req))

	bytes := util.FromASCII85(req.Content)
	contentHash := util.ToSHA1(bytes)

	if req.ContentHash != contentHash {
		return errInvalidContentHash(c)
	}

	var item *db.Item

	if req.Label != "" {
		if isReservedLabel(req.Label) {
			return errReservedRef(c, req.Label)
		}
		if foundItem, found := db.FindItemByLabel(req.Label); found {
			if !req.Force {
				return errDuplicateRef(c, req.Label)
			}
			item = foundItem
		}
		if foundItem, found := db.FindItemByUUID(req.Label); found {
			if !req.Force {
				return errDuplicateRef(c, req.Label)
			}
			item = foundItem
		}
	}

	blob := db.FindOrCreateBlob(bytes)

	if item != nil {
		item.UpdateKind(req.Kind)
		item.UpdateBlobHash(blob.Hash)
		if req.ContentType != "" {
			item.UpdateContentType(req.ContentType)
		}
	} else {
		item = db.CreateItem(&db.Item{
			UUID:     shortuuid.New(),
			Label:    req.Label,
			Kind:     req.Kind,
			BlobHash: blob.Hash,
		})
	}

	return c.JSON(types.UploadResponse{
		URL: item.URL(),
	})
}
