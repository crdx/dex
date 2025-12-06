package util

import (
	"bytes"
	"crypto/sha1"
	"encoding/ascii85"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/samber/lo"
)

// MapTrimSpace is a lo.Map-compatible function that trims leading and trailing spaces from its
// arguments.
func MapTrimSpace(s string, _ int) string {
	return strings.TrimSpace(s)
}

func ToASCII85(b []byte) string {
	var buf bytes.Buffer
	encoder := ascii85.NewEncoder(&buf)
	lo.Must(encoder.Write(b))
	_ = encoder.Close()
	return buf.String()
}

func FromASCII85(b string) []byte {
	reader := ascii85.NewDecoder(strings.NewReader(b))
	var decoded bytes.Buffer
	lo.Must(io.Copy(&decoded, reader))
	return decoded.Bytes()
}

func ToSHA1(b []byte) string {
	hash := sha1.Sum(b)
	return hex.EncodeToString(hash[:])
}

func Truncate(s string, max int) string {
	if len(s) > max {
		return s[:max]
	}
	return s
}

func Hyperlink(link string, text string) string {
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", link, text)
}

func Pluralise(count int, singular, plural string) string {
	if count == 1 {
		return singular
	} else {
		return plural
	}
}

// FormatTimeSince formats a time as the relative time since t compared to time.Now().UTC(), so t
// should be in UTC.
//
// See documentation for FormatDuration for the remainder of the parameters.
func FormatTimeSince(t time.Time, long bool, precision int, suffix string) string {
	return FormatDuration(time.Now().UTC().Sub(t), long, precision, suffix)
}

// FormatTimeUntil formats a time as the relative time until t compared to time.Now().UTC(), so t
// should be in UTC.
//
// See documentation for FormatDuration for the remainder of the parameters.
func FormatTimeUntil(t time.Time, long bool, precision int, suffix string) string {
	return FormatDuration(t.Sub(time.Now().UTC()), long, precision, suffix)
}

// FormatDuration formats a duration as a relative time.
//
// Precision is the number of units to include. For example, for 65 seconds a precision of 1 would
// return "1 min" and a precision of 2 would return "1 min 5 secs".
func FormatDuration(duration time.Duration, long bool, precision int, suffix string) string {
	seconds := int(duration.Seconds())

	if seconds == 0 {
		if long {
			return "just now"
		} else {
			return "now"
		}
	}

	type Unit struct {
		longName  string
		shortName string
		seconds   int
	}

	units := []Unit{
		{"year", "y", 60 * 60 * 24 * 7 * 52},
		{"week", "w", 60 * 60 * 24 * 7},
		{"day", "d", 60 * 60 * 24},
		{"hour", "h", 60 * 60},
		{"min", "m", 60},
		{"sec", "s", 0},
	}

	var parts []string

	for _, unit := range units {
		if seconds < unit.seconds {
			continue
		}

		var partSeconds int
		if unit.seconds > 0 {
			partSeconds = seconds / unit.seconds
			seconds %= unit.seconds
		} else {
			partSeconds = seconds
		}

		if partSeconds > 0 {
			if long {
				parts = append(parts, fmt.Sprintf(
					"%d %s",
					partSeconds,
					Pluralise(partSeconds, unit.longName, unit.longName+"s"),
				))
			} else {
				parts = append(parts, fmt.Sprintf(
					"%d%s",
					partSeconds,
					unit.shortName,
				))
			}

			if precision > 0 {
				precision--
				if precision == 0 {
					break
				}
			}
		}
	}

	s := strings.Join(parts, " ")

	if suffix != "" {
		return s + " " + suffix
	}

	return s
}

// ToLocal converts a time to another timezone.
func ToLocal(t time.Time) time.Time {
	return t.In(lo.Must(time.LoadLocation("Europe/London")))
}

// EnsureEndsIn ensures xs ends in x (if it doesn't already).
func EnsureEndsIn(xs []byte, x byte) []byte {
	if len(xs) == 0 || xs[len(xs)-1] != x {
		return append(xs, x)
	}

	return xs
}

func Itemise(singular, plural string, n int) string {
	return fmt.Sprintf("%d %s", n, Pluralise(n, singular, plural))
}

// FormatSize formats a byte count as a human-readable size.
func FormatSize(n int) string {
	if n > 1024*1024 {
		return Itemise("MB", "MB", n/1024/1024)
	} else if n > 1024 {
		return Itemise("KB", "KB", n/1024)
	}
	return Itemise("byte", "bytes", n)
}

func ItemRequestPath(s string) string {
	if s == "/" {
		return s
	}
	return "/" + s
}

func ItemRef(s string) string {
	if s == "" {
		return "/"
	}
	return s
}

func PathExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func PassThrough(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
