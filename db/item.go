package db

import (
	"cmp"
	"mime"
	"net/http"
	"net/url"
	"path"
	"path/filepath"

	"crdx.org/dex/cmd/dexd/env"
	"crdx.org/dex/pkg/types"
	"crdx.org/dex/pkg/util"
	"crdx.org/dex/pkg/web"
)

// FindOrCreateBlob finds an existing blob by content hash, or creates a new one.
func FindOrCreateBlob(content []byte) *Blob {
	hash := util.ToSHA1(content)
	if blob, found := FindBlobByHash(hash); found {
		return blob
	}
	return CreateBlob(&Blob{
		Hash:    hash,
		Content: content,
	})
}

func (self *Item) Content() []byte {
	blob, _ := FindBlobByHash(self.BlobHash)
	return blob.Content
}

func (self *Item) URL() string {
	if path, err := url.JoinPath(env.BaseURL(), self.Ref()); err == nil {
		return path
	}
	return ""
}

func (self *Item) Ref() string {
	return cmp.Or(self.Label, self.UUID)
}

func (self *Item) Base() string {
	return path.Base(self.Ref())
}

func (self *Item) Ext() string {
	return filepath.Ext(self.Ref())
}

func (self *Item) ResolvedContentType() web.ContentType {
	if self.Kind == types.KindRedir {
		return web.ContentType{}
	}

	value := self.ContentType

	if value == "" {
		value = mime.TypeByExtension(self.Ext())
	}

	if value == "" {
		value = http.DetectContentType(self.Content())
	}

	if value == "application/octet-stream" && self.Kind == types.KindPaste {
		value = "text/plain"
	}

	parsedValue, _, err := mime.ParseMediaType(value)
	if err != nil {
		parsedValue = value
	}

	return web.ContentTypeFrom(parsedValue)
}
