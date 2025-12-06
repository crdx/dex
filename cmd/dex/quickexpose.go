package main

import "crdx.org/dex/pkg/types"

func quickExpose(ref string) error {
	if err := uploadContent(ref, []byte{}, types.KindPaste, "", false, true); err != nil {
		return err
	}
	return expose(ref, "1d")
}
