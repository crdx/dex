package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"crdx.org/col"
	"crdx.org/dex/pkg/types"
	"crdx.org/dex/pkg/util"
)

func deployURLs() error {
	res, err := get[types.DeployURLsResponse]("urls", nil)
	if err != nil {
		return err
	}

	if len(res.Tokens) == 0 {
		return errors.New("no deploy urls")
	}

	table := newPrettyTable([]any{"ID", "Item", "Deploy URL", "Expires", "Created"})
	for _, token := range res.Tokens {
		var expires string
		if token.ExpiresAt != "" {
			expires = util.FormatTimeUntil(mustParseTime(token.ExpiresAt), false, 1, "")
		} else {
			expires = "Never"
		}

		table.AddRow([]any{
			token.ID,
			token.Ref,
			token.DeployURL,
			expires,
			util.FormatTimeSince(mustParseTime(token.CreatedAt), false, 1, "ago"),
		})
	}

	table.Render()
	return nil
}

func deployURLsDelete(ids []string) error {
	for _, id := range ids {
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid id: %s", id)
		}

		_, err = post[types.DeployURLsDeleteResponse]("urls/delete/"+strconv.FormatInt(idInt, 10), nil)
		if err != nil {
			return err
		}

		fmt.Printf("Deleted %d\n", idInt)
	}

	return nil
}

func deploys(verbose bool, all bool) error {
	qs := map[string]string{}
	if all {
		qs["all"] = "1"
	}
	res, err := get[types.DeploymentsResponse]("deployments", qs)
	if err != nil {
		return err
	}

	if len(res.Deployments) == 0 {
		return errors.New("no deployments")
	}

	headers := []any{"Item", "Deployer", "Change", "IP", "When"}
	if verbose {
		headers = append(headers, "User Agent")
		headers = append(headers, "ID")
	}

	table := newPrettyTable(headers)
	for _, d := range res.Deployments {
		row := []any{
			d.Ref,
			d.Deployer,
			util.Truncate(d.Change, 100),
			d.IPAddress,
			util.FormatTimeSince(mustParseTime(d.CreatedAt), false, 1, "ago"),
		}
		if verbose {
			row = append(row, util.Truncate(d.UserAgent, 50))
			row = append(row, d.ID)
		}
		if d.Deleted {
			for i, cell := range row {
				row[i] = col.Dim(cell)
			}
		}
		table.AddRow(row)
	}

	table.Render()
	return nil
}

func expose(ref string, expiry string) error {
	var expirySeconds int64

	if expiry != "" {
		duration, err := parseDuration(expiry)
		if err != nil {
			return fmt.Errorf("invalid expiry: %w", err)
		}
		expirySeconds = int64(duration.Seconds())
	}

	payload := types.ExposeRequest{
		Ref:           ref,
		ExpirySeconds: expirySeconds,
	}

	res, err := post[types.ExposeResponse]("expose", payload)
	if err != nil {
		return err
	}

	fmt.Printf("Public URL: %s\n", res.PublicURL)
	if res.Warning != "" {
		fmt.Printf("Deploy URL: %s %s\n", res.DeployURL, col.Yellow("(pre-existing)"))
	} else {
		fmt.Printf("Deploy URL: %s\n", res.DeployURL)
	}

	if res.ExpiresAt != "" {
		t, _ := time.Parse(time.RFC3339, res.ExpiresAt)
		fmt.Printf("Expires at: %s\n", t.Format("2 Jan 2006 15:04"))
	}

	fmt.Printf("\nRun curl %s to find out how to deploy.\n", res.DeployURL)

	return nil
}
