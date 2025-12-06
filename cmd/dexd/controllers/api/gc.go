package api

import (
	"crdx.org/dex/db"
	"crdx.org/dex/pkg/types"
	"github.com/gofiber/fiber/v3"
)

func GC(c fiber.Ctx) error {
	var deletedHashes []string
	for _, blob := range db.FindBlobs() {
		if _, found := db.FindItemByBlobHash(blob.Hash); !found {
			deletedHashes = append(deletedHashes, blob.Hash)
			blob.HardDelete()
		}
	}
	return c.JSON(types.GCResponse{
		DeletedHashes: deletedHashes,
	})
}
