package main

import (
	"errors"
	"os"
	"path"

	"crdx.org/dex/pkg/types"
	"crdx.org/dex/pkg/util"
)

func cat(ref string, force bool) error {
	res, err := get[types.CatResponse](path.Join("cat", ref), nil)
	if err != nil {
		return err
	}

	binary := res.Kind == types.KindFile || !res.IsText

	if binary && isInteractive() && !force {
		return errors.New("use -f to output a file")
	}

	if _, err := os.Stdout.Write(util.EnsureEndsIn(res.Content, '\n')); err != nil {
		return err
	}

	return nil
}
