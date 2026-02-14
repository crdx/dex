package main

import (
	"log"

	"crdx.org/col"
	"crdx.org/dex/pkg/types"
)

func gc() error {
	res, err := post[types.GCResponse]("gc", nil)
	if err != nil {
		return err
	}

	for _, hash := range res.DeletedHashes {
		log.Printf(col.Yellow("deleted dangling blob %s"), hash)
	}

	return nil
}
