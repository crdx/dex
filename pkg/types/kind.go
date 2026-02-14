package types

type Kind = string

var Kinds = []Kind{
	KindFile,
	KindPaste,
	KindRedir,
}

var (
	KindFile  Kind = "file"
	KindPaste Kind = "paste"
	KindRedir Kind = "redir"
)
