package main

import (
	"errors"
	"log"
	"path"

	"crdx.org/col"
	"crdx.org/dex/pkg/types"
)

func delete(refs []string) error {
	var errs []error
	for _, ref := range refs {
		res, err := post[types.DeleteResponse](path.Join("delete", ref), nil)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		log.Println(col.Green(res.Message))
	}
	return errors.Join(errs...)
}
