package main

import (
	"fmt"

	"crdx.org/dex/pkg/types"
)

func copy(ref string, label string) error {
	payload := types.CopyRequest{
		From: ref,
		To:   label,
	}

	res, err := post[types.CopyResponse]("cp", payload)
	if err != nil {
		return err
	}
	fmt.Println(res.URL)
	return nil
}
