package main

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"crdx.org/dex/cmd/dex/config"
	"crdx.org/dex/pkg/types"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
)

func isInteractive() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

type prettyTable struct {
	t table.Writer
}

// newPrettyTable creates a new pretty table with the given headers. A pretty table has borders and
// headers. Best suited for displaying tables.
func newPrettyTable(headers []any) *prettyTable {
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.Style().Format.Header = text.FormatDefault
	t.AppendHeader(headers)
	return &prettyTable{t: t}
}

// AddRow adds a row to the table.
func (self *prettyTable) AddRow(row []any) {
	self.t.AppendRow(row)
}

// AddRows adds multiple rows to the table.
func (self *prettyTable) AddRows(rows [][]any) {
	for _, row := range rows {
		self.t.AppendRow(row)
	}
}

// Render renders the table to stdout.
func (self *prettyTable) Render() {
	self.t.SetOutputMirror(os.Stdout)
	self.t.Render()
}

func post[T any](path string, payload any) (T, error) {
	return request[T](http.MethodPost, path, payload, nil)
}

func get[T any](path string, qs map[string]string) (T, error) {
	return request[T](http.MethodGet, path, nil, qs)
}

func request[T any](method string, path string, payload any, qs map[string]string) (T, error) {
	var success T
	var failure types.FailureResponse

	req := config.Request().
		SetSuccessResult(&success).
		SetErrorResult(&failure)

	if payload != nil {
		req.SetBody(payload)
	}

	if qs != nil {
		req.SetQueryParams(qs)
	}

	res, err := req.Send(method, config.Endpoint(path))
	if err != nil {
		return success, err
	}

	if res.IsSuccessState() {
		return success, nil
	}

	return success, errors.New(failure.Message)
}

// nextLabel returns the next label available for use.
func nextLabel() (string, error) {
	res, err := get[types.ListResponse]("list", nil)
	if err != nil {
		return "", err
	}

	value := 0
	for _, item := range res.Items {
		n, err := strconv.Atoi(item.Ref)
		if err != nil {
			continue
		}

		if n > value {
			value = n
		}
	}

	return strconv.Itoa(value + 1), nil
}

func parseDuration(s string) (time.Duration, error) {
	if daysStr, found := strings.CutSuffix(s, "d"); found {
		days, err := strconv.Atoi(daysStr)
		if err != nil {
			return 0, err
		}
		return time.Duration(days) * 24 * time.Hour, nil
	}
	return time.ParseDuration(s)
}

func mustParseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}
