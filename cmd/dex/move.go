package main

import (
	"fmt"

	"crdx.org/dex/pkg/types"
)

func move(ref string, label string) error {
	payload := types.MoveRequest{
		FromRef: ref,
		ToLabel: label,
	}

	res, err := post[types.MoveResponse]("mv", payload)
	if err != nil {
		return err
	}
	fmt.Println(res.URL)
	return nil
}
