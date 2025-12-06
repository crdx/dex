package main

import (
	"fmt"
	"path"
	"time"

	"crdx.org/col"
	"crdx.org/dex/pkg/types"
	"crdx.org/dex/pkg/util"
)

func tail(verbose bool) error {
	var since int64
	qs := map[string]string{
		"since":   fmt.Sprint(since),
		"verbose": fmt.Sprint(verbose),
	}
	for {
		res, err := get[types.TailResponse]("tail", qs)
		if err != nil {
			return err
		}

		for _, log := range res.Logs {
			if log.ID > since {
				since = log.ID
			}
			printTailResponse(log)
		}

		qs["since"] = fmt.Sprint(since)

		if !isInteractive() {
			break
		}

		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

func printTailResponse(res types.TailResponseItem) {
	var requestStr string
	if res.Exists && !res.Renamed {
		requestStr = col.Yellow(util.Hyperlink(path.Join(res.BaseURL, res.Request), res.Request))
	} else {
		requestStr = col.Red(res.Request)
	}

	var renamedStr string
	if res.Renamed {
		requestPath := util.ItemRequestPath(res.Ref)
		if res.Exists {
			renamedStr = fmt.Sprintf("(→ %s) ", col.Yellow(util.Hyperlink(path.Join(res.BaseURL, res.Ref), requestPath)))
		} else {
			renamedStr = fmt.Sprintf("(→ %s) ", col.Red(requestPath))
		}
	}

	fmt.Printf(
		"[%s] %s %s %s%s\n",
		util.ToLocal(res.CreatedAt).Format(time.DateTime),
		col.Cyan(util.Hyperlink("https://whatismyipaddress.com/ip/"+res.IPAddress, res.IPAddress)),
		requestStr,
		renamedStr,
		col.Magenta(res.UserAgent),
	)
}
