package api

import (
	"regexp"
	"strings"

	"crdx.org/dex/db"
	"crdx.org/dex/pkg/types"
	"crdx.org/dex/pkg/util"
	"github.com/gofiber/fiber/v3"
	"github.com/samber/lo"
)

func List(c fiber.Ctx) error {
	var items []types.ListResponseItem

	for _, row := range db.FindItemsForList() {
		item, _ := db.FindItem(row.ID)
		contentType := item.ResolvedContentType()

		var meta string

		switch true {
		case item.Kind == types.KindRedir:
			meta = generateURLMeta(string(row.Content))
		case item.Kind == types.KindFile || !contentType.IsText:
			meta = generateFileMeta(row.Content)
		case item.Kind == types.KindPaste:
			meta = generateTextMeta(string(row.Content))
		}

		items = append(items, types.ListResponseItem{
			URL:         item.URL(),
			Ref:         item.Ref(),
			Kind:        item.Kind,
			Hits:        item.Hits,
			ContentType: contentType.MediaType,
			IsText:      contentType.IsText,
			Hash:        item.BlobHash,
			LastHitAt:   item.LastHitAt,
			Meta:        meta,
		})
	}

	return c.JSON(types.ListResponse{
		Items: items,
	})
}

func generateURLMeta(s string) string {
	return s
}

func generateFileMeta(b []byte) string {
	return util.FormatSize(len(b))
}

func generateTextMeta(s string) string {
	var items []string

	lines := lo.Compact(regexp.MustCompile("\n+").Split(s, -1))

	items = append(items, util.Itemise("line", "lines", len(lines)))
	items = append(items, util.Itemise("char", "chars", len(s)))
	items = append(items, util.FormatSize(len(s)))

	return strings.Join(items, ", ")
}
