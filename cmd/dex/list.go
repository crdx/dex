package main

import (
	"errors"
	"fmt"

	"crdx.org/dex/pkg/types"
	"crdx.org/dex/pkg/util"
	"github.com/samber/lo"
)

func list(verbose bool, refs bool) error {
	res, err := get[types.ListResponse]("list", nil)
	if err != nil {
		return err
	}

	if refs {
		for _, item := range res.Items {
			if item.Ref == "/" {
				continue
			}
			fmt.Println(item.Ref)
		}
		return nil
	}

	if len(res.Items) == 0 {
		return errors.New("no items")
	}

	headers := []any{"Type", "Item", "Hits", "Last Hit"}
	if verbose {
		headers = append(headers, "Content Type", "Meta")
	}

	table := newPrettyTable(headers)
	for _, item := range res.Items {
		row := []any{
			item.Kind,
			item.URL,
			item.Hits,
		}

		if item.LastHitAt.Valid {
			row = append(row, util.FormatTimeSince(item.LastHitAt.V, false, 1, "ago"))
		} else {
			row = append(row, "Never")
		}

		if verbose {
			row = append(row, lo.If(item.ContentType == "", "...").Else(item.ContentType))
			row = append(row, item.Meta)
		}

		table.AddRow(row)
	}

	table.Render()
	return nil
}
