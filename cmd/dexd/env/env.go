package env

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
)

const (
	ModeDevelopment = "development"
	ModeProduction  = "production"

	LogTypeAll    = "all"
	LogTypeDisk   = "disk"
	LogTypeStderr = "stderr"
	LogTypeNone   = "none"
)

var env map[string]string

var (
	mode = func() string { return env["MODE"] }

	Debug      = func() bool { return truthy(env["DEX_DEBUG"]) }
	Production = func() bool { return env["MODE"] == ModeProduction }

	Host = func() string { return env["HOST"] }
	Port = func() string { return env["PORT"] }

	DatabaseName     = func() string { return env["DB_NAME"] }
	DatabaseUsername = func() string { return env["DB_USERNAME"] }
	DatabasePassword = func() string { return env["DB_PASSWORD"] }
	DatabaseProtocol = func() string { return env["DB_PROTOCOL"] }
	DatabaseAddress  = func() string { return env["DB_ADDRESS"] }

	BaseURL = func() string { return env["BASE_URL"] }
	APIKey  = func() string { return env["API_KEY"] }

	TrustedProxies = func() string { return env["TRUSTED_PROXIES"] }
)

func Init() {
	if env == nil {
		env = map[string]string{}
	}

	for _, v := range os.Environ() {
		name, value, ok := strings.Cut(v, "=")
		if !ok {
			continue
		}

		env[name] = value
	}
}

func InitFrom(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	env = parse(string(b))
	Init()

	return nil
}

func Validate() error {
	var errs []error
	e := func(err error) {
		errs = append(errs, err)
	}

	e(require(mode, "MODE"))

	if Production() && Port() == "" {
		// In development no port means use a random port, but this will never be correct for production.
		return fmt.Errorf("running in production but no port set")
	}

	e(require(Host, "HOST"))
	e(require(DatabaseName, "DB_NAME"))

	e(require(APIKey, "API_KEY"))
	e(require(BaseURL, "BASE_URL"))

	if err := errors.Join(errs...); err != nil {
		return errors.New("missing environment variables:\n" + err.Error())
	}

	return nil
}

func parse(s string) map[string]string {
	m := map[string]string{}

	for line := range strings.SplitSeq(s, "\n") {
		line := strings.TrimSpace(line)

		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		name, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		if len(value) > 0 {
			n := len(value) - 1
			if value[0] == '"' && value[n] == '"' {
				value = value[1:n]
			}
		}

		m[name] = value
	}

	return m
}

func require(f func() string, name string) error {
	if f() == "" {
		return fmt.Errorf("    %s required", name)
	}
	return nil
}

func requireIn(f func() string, name string, values []string, canBeEmpty bool) error {
	if !canBeEmpty {
		if err := require(f, name); err != nil {
			return err
		}
	}

	value := f()

	if canBeEmpty && value == "" {
		return nil
	}

	if !slices.Contains(values, value) {
		s := ""
		if canBeEmpty {
			s = `, or the empty string ("")`
		}

		return fmt.Errorf(
			`%s contains an invalid value (must be one of: "%s"%s)`,
			name,
			strings.Join(values, `", "`),
			s,
		)
	}

	return nil
}

func or(name string, default_ string) string { //nolint:unused
	value := env[name]

	if value != "" {
		return value
	} else {
		return default_
	}
}

func truthy(s string) bool {
	return slices.Contains([]string{"true", "1", "yes"}, s)
}
