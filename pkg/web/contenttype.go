package web

import "strings"

type ContentType struct {
	MediaType string
	IsText    bool
}

func (self ContentType) String() string {
	return self.MediaType
}

// ContentTypeHeader returns the content type formatted for an HTTP header, including charset
// for text-based types.
func (self ContentType) ContentTypeHeader() string {
	if self.MediaType == "" {
		return ""
	}
	if self.IsText {
		return self.MediaType + "; charset=utf-8"
	}
	return self.MediaType
}

func ContentTypeFrom(mediaType string) ContentType {
	return ContentType{
		MediaType: mediaType,
		IsText:    IsText(mediaType),
	}
}

func IsText(mediaType string) bool {
	if strings.HasPrefix(mediaType, "text/") {
		return true
	}

	if strings.HasSuffix(mediaType, "+xml") || strings.HasSuffix(mediaType, "+json") {
		return true
	}

	textTypes := map[string]bool{
		"application/json":                  true,
		"application/xml":                   true,
		"application/javascript":            true,
		"application/xhtml+xml":             true,
		"application/x-www-form-urlencoded": true,
		"application/sql":                   true,
		"application/graphql":               true,
		"application/ld+json":               true,
		"application/rss+xml":               true,
		"application/atom+xml":              true,
		"application/mathml+xml":            true,
		"application/csv":                   true,
		"application/yaml":                  true,
		"application/x-sh":                  true,
		"application/x-csh":                 true,
		"application/x-perl":                true,
		"application/x-python":              true,
		"application/x-ruby":                true,
	}

	return textTypes[mediaType]
}
