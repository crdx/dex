package mail

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDSN(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected *Config
		err      string
	}{
		{
			name:  "full DSN with all fields",
			input: "smtp://user:pass@smtp.example.com:587?from=sender@example.com&to=recipient@example.com",
			expected: &Config{
				Host:     "smtp.example.com",
				Port:     "587",
				Username: "user",
				Password: "pass",
				From:     "sender@example.com",
				To:       "recipient@example.com",
			},
		},
		{
			name:  "DSN with from_name",
			input: "smtp://user:pass@smtp.example.com:587?from=sender@example.com&from_name=dex&to=recipient@example.com",
			expected: &Config{
				Host:     "smtp.example.com",
				Port:     "587",
				Username: "user",
				Password: "pass",
				From:     "sender@example.com",
				FromName: "dex",
				To:       "recipient@example.com",
			},
		},
		{
			name:  "DSN without credentials",
			input: "smtp://smtp.example.com:25?from=sender@example.com&to=recipient@example.com",
			expected: &Config{
				Host:     "smtp.example.com",
				Port:     "25",
				Username: "",
				Password: "",
				From:     "sender@example.com",
				To:       "recipient@example.com",
			},
		},
		{
			name:  "DSN with default port",
			input: "smtp://user:pass@smtp.example.com?from=sender@example.com&to=recipient@example.com",
			expected: &Config{
				Host:     "smtp.example.com",
				Port:     "587",
				Username: "user",
				Password: "pass",
				From:     "sender@example.com",
				To:       "recipient@example.com",
			},
		},
		{
			name:  "DSN with username only",
			input: "smtp://user@smtp.example.com:587?from=sender@example.com&to=recipient@example.com",
			expected: &Config{
				Host:     "smtp.example.com",
				Port:     "587",
				Username: "user",
				Password: "",
				From:     "sender@example.com",
				To:       "recipient@example.com",
			},
		},
		{
			name:  "invalid scheme",
			input: "http://smtp.example.com?from=sender@example.com&to=recipient@example.com",
			err:   "invalid scheme: expected smtp, got http",
		},
		{
			name:  "missing from parameter",
			input: "smtp://smtp.example.com?to=recipient@example.com",
			err:   "missing 'from' parameter in DSN",
		},
		{
			name:  "missing to parameter",
			input: "smtp://smtp.example.com?from=sender@example.com",
			err:   "missing 'to' parameter in DSN",
		},
		{
			name:  "invalid URL",
			input: "://invalid",
			err:   "invalid DSN",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			actual, err := parseDSN(testCase.input)

			if testCase.err != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, testCase.expected, actual)
			}
		})
	}
}
