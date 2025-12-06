package main

import (
	"crdx.org/dex/pkg/types"
)

func set(ref string, key types.Key, value string) error {
	payload := types.SetRequest{
		Ref:   ref,
		Key:   key,
		Value: value,
	}

	_, err := post[types.SetResponse]("set", payload)
	if err != nil {
		return err
	}
	return nil
}
