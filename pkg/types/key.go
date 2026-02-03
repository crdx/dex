package types

type Key = string

var (
	KeyContentType Key = "ContentType"
	KeyKind        Key = "Kind"
)

func ParseKey(contentType bool, kind bool) Key {
	if contentType {
		return KeyContentType
	} else if kind {
		return KeyKind
	} else {
		panic("invalid key")
	}
}
