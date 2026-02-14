package main

import "crdx.org/dex/cmd/dex/config"

func conf(edit bool) error { //nolint:unparam
	if edit {
		config.Edit()
	} else {
		config.Print()
	}
	return nil
}
