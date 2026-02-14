package main

import "crdx.org/dex/pkg/types"

func redir(label string, url string, kind types.Kind, force bool) error {
	if label == "" {
		var err error
		label, err = nextLabel()
		if err != nil {
			return err
		}
	}

	return uploadContent(label, []byte(url), kind, "", force, false)
}
