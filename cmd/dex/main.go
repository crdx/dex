package main

import (
	"log"

	"crdx.org/col"
	"crdx.org/dex/cmd/dex/config"
	"crdx.org/dex/pkg/types"
	"crdx.org/duckopt/v2"
)

func getUsage() string {
	return `
		Usage:
			$0 [options] (paste | upload) [<paths>...] [--label NAME] [--uuid] [--force] [--type NAME] [--bare]
			$0 [options] redir <url> [--label NAME] [--force]
			$0 [options] cat <ref> [--force]
			$0 [options] rm <refs>...
			$0 [options] ls [--refs]
			$0 [options] (cp | mv) <from> <to>
			$0 [options] set <ref> (type | kind) <value>
			$0 [options] expose <ref> [--expiry DURATION]
			$0 [options] urls [rm <ids>...]
			$0 [options] deploys [--all]
			$0 [options] (conf [-e] | e <label> | tail | gc)

		Options:
			--refs              List refs
			-u, --uuid          Use UUID instead of filenames or numbers
			-t, --type TYPE     Content type e.g. application/json
			-f, --force         Use the force
			-e, --edit          Open in $EDITOR
			-l, --label NAME    Short name for this item
			-x, --expiry DUR    Token expiry duration e.g. 1d, 12h, 30m
			-b, --bare          Create empty item (requires --label)
			-a, --all           Show all results
			-v, --verbose       Be verbose
	`
}

type Opts struct {
	Redir       bool `docopt:"redir"`
	Paste       bool `docopt:"paste"`
	Upload      bool `docopt:"upload"`
	Cat         bool `docopt:"cat"`
	List        bool `docopt:"ls"`
	Gc          bool `docopt:"gc"`
	Copy        bool `docopt:"cp"`
	Move        bool `docopt:"mv"`
	Set         bool `docopt:"set"`
	Type        bool `docopt:"type"`
	Kind        bool `docopt:"kind"`
	Conf        bool `docopt:"conf"`
	Tail        bool `docopt:"tail"`
	Expose      bool `docopt:"expose"`
	QuickExpose bool `docopt:"e"`
	DeployURLs  bool `docopt:"urls"`
	Rm          bool `docopt:"rm"`
	Deploys     bool `docopt:"deploys"`

	URL   string   `docopt:"<url>"`
	Ref   string   `docopt:"<ref>"`
	Label string   `docopt:"<label>"`
	Value string   `docopt:"<value>"`
	From  string   `docopt:"<from>"`
	To    string   `docopt:"<to>"`
	IDs   []string `docopt:"<ids>"`

	Paths []string `docopt:"<paths>"`
	Refs  []string `docopt:"<refs>"`

	LabelF      string `docopt:"--label"`
	ContentType string `docopt:"--type"`
	Expiry      string `docopt:"--expiry"`

	RefsF   bool `docopt:"--refs"`
	Edit    bool `docopt:"--edit"`
	UUID    bool `docopt:"--uuid"`
	Force   bool `docopt:"--force"`
	Bare    bool `docopt:"--bare"`
	All     bool `docopt:"--all"`
	Verbose bool `docopt:"--verbose"`
}

func main() {
	log.SetFlags(0)
	opts := duckopt.MustBind[Opts](getUsage(), "$0")
	col.Init()

	if err := config.Load(); err != nil {
		log.Fatal(col.Red(err.Error()))
	}

	var err error

	switch {
	case opts.Set:
		err = set(opts.Ref, types.ParseKey(opts.Type, opts.Kind), opts.Value)
	case opts.Copy:
		err = copy(opts.From, opts.To)
	case opts.Move:
		err = move(opts.From, opts.To)
	case opts.Redir:
		err = redir(opts.LabelF, opts.URL, types.KindRedir, opts.Force)
	case opts.Paste:
		err = upload(opts.UUID, opts.LabelF, opts.Paths, types.KindPaste, opts.ContentType, opts.Force, opts.Bare)
	case opts.Upload:
		err = upload(opts.UUID, opts.LabelF, opts.Paths, types.KindFile, opts.ContentType, opts.Force, opts.Bare)
	case opts.Rm:
		if opts.DeployURLs {
			err = deployURLsDelete(opts.IDs)
		} else {
			err = delete(opts.Refs)
		}
	case opts.Cat:
		err = cat(opts.Ref, opts.Force)
	case opts.List:
		err = list(opts.Verbose, opts.RefsF)
	case opts.Conf:
		err = conf(opts.Edit)
	case opts.Gc:
		err = gc()
	case opts.Tail:
		err = tail(opts.Verbose)
	case opts.Expose:
		err = expose(opts.Ref, opts.Expiry)
	case opts.DeployURLs:
		err = deployURLs()
	case opts.Deploys:
		err = deploys(opts.Verbose, opts.All)
	case opts.QuickExpose:
		err = quickExpose(opts.Label)
	}

	if err != nil {
		log.Fatal(col.Red(err.Error()))
	}
}
