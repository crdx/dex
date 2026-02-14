package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"crdx.org/dex/pkg/types"
	"crdx.org/dex/pkg/util"
)

func upload(uuid bool, label string, paths []string, kind types.Kind, contentType string, force bool, bare bool) error {
	if bare {
		if label == "" && !uuid {
			return errors.New("--bare requires --label or --uuid")
		}
		return uploadContent(label, []byte{}, kind, contentType, force, false)
	}

	if len(paths) == 0 {
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		if label == "" && !uuid {
			if label, err = nextLabel(); err != nil {
				return err
			}
		}
		return uploadContent(label, content, kind, contentType, force, false)
	}

	if len(paths) > 1 && label != "" {
		return errors.New("cannot specify label with multiple paths")
	}

	customLabel := label != ""

	var errs []error
	for _, path := range paths {
		if !customLabel {
			if uuid {
				label = ""
			} else {
				label = filepath.Base(path)
			}
		}

		content, err := os.ReadFile(path) //nolint:gosec // G304: path is from user-provided CLI args
		if err != nil {
			errs = append(errs, err)
			continue
		}

		errs = append(errs, uploadContent(label, content, kind, contentType, force, false))
	}

	return errors.Join(errs...)
}

func uploadContent(label string, content []byte, kind types.Kind, contentType string, force bool, quiet bool) error {
	payload := types.UploadRequest{
		Label:       label,
		Content:     util.ToASCII85(content),
		ContentHash: util.ToSHA1(content),
		ContentType: contentType,
		Kind:        kind,
		Force:       force,
	}

	res, err := post[types.UploadResponse]("upload", payload)
	if err != nil {
		return err
	}

	if !quiet {
		fmt.Println(res.URL)
	}
	return nil
}
