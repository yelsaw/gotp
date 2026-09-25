package gotp

import (
	"os"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestCleanString(t *testing.T) {
	test_strs := []struct {
		unexpected string
		expected   string
	}{
		{"", ""},
		{"  space  ", "space"},
		{"not%20url%20encoded", "not url encoded"},
		{"hello+world", "hello world"},
	}

	for _, test := range test_strs {
		for_res := cleanString(test.unexpected)
		if for_res != test.expected {
			t.Errorf("unexpected(%q) = %q, expected %q", test.unexpected, for_res, test.expected)
		}
	}
}

func TestArgParser(t *testing.T) {

	tmpF, err := os.CreateTemp("", "test_gotp_*.txt")
	if err != nil {
		t.Fatalf("Temp file failed to create: %v", err)
	}
	defer os.Remove(tmpF.Name())

	testURL := "otpauth://totp/Microsoft:you@youremail.com?algorithm=SHA1&digits=6&issuer=Microsoft&period=30&secret=VXYU6YKSNBZELU23"

	if _, err := tmpF.WriteString(testURL); err != nil {
		t.Fatalf("Temp file cannot be written to: %v", err)
	}
	tmpF.Close()

	test_args := []struct {
		name     string
		input    string
		hasError bool
		expected string
	}{
		{
			name:     "valid file path",
			input:    tmpF.Name(),
			hasError: false,
			expected: testURL,
		},
		{
			name:     "valid direct string",
			input:    testURL,
			hasError: false,
			expected: testURL,
		},
		{
			name:     "empty input",
			input:    "",
			hasError: true,
			expected: "",
		},
		{
			name:     "non-existent file treated as string",
			input:    "/non/existent/file.txt",
			hasError: false,
			expected: "/non/existent/file.txt",
		},
	}

	for _, test := range test_args {
		t.Run(test.name, func(t *testing.T) {
			result, err := ArgParser(test.input)

			if test.hasError && err == nil {
				t.Errorf("Expected error for input %q, but got none", test.input)
			}

			if !test.hasError && err != nil {
				t.Errorf("Unexpected error for input %q: %v", test.input, err)
			}

			if result != test.expected {
				t.Errorf("ArgParser(%q) = %q, expected %q", test.input, result, test.expected)
			}
		})
	}
}

func TestUrlParser(t *testing.T) {
	test_urls := []struct {
		name             string
		input            string
		hasError         bool
		expectedProvider string
		expectedAccount  string
	}{
		{
			name:             "OTP URL with issuer",
			input:            "otpauth://totp/Microsoft:me@microsoft.com?secret=JBSWY3DPEHPK3PXP&issuer=Microsoft&period=30",
			hasError:         false,
			expectedProvider: "Microsoft",
			expectedAccount:  "me@microsoft.com",
		},
		{
			name:             "valid OTP URL without issuer",
			input:            "otpauth://totp/Service:me@service.com?secret=JBSWY3DPEHPK3PXP&period=30",
			hasError:         false,
			expectedProvider: "Service",
			expectedAccount:  "me@service.com",
		},
		{
			name:     "empty URL",
			input:    "",
			hasError: true,
		},
		{
			name:     "invalid protocol",
			input:    "http://example.com",
			hasError: true,
		},
		{
			name:     "invalid secret",
			input:    "otpauth://totp/Service:me@service.com?secret=NOTAREALPUPPY!@#&issuer=Service",
			hasError: true,
		},
	}

	for _, test := range test_urls {
		t.Run(test.name, func(t *testing.T) {
			result, err := UrlParser(test.input)

			if test.hasError && err == nil {
				t.Errorf("Expected error for input %q, but got none", test.input)
			}

			if !test.hasError {
				if err != nil {
					t.Errorf("Input failed %q: %v", test.input, err)
				} else {
					if result.provider != test.expectedProvider {
						t.Errorf("Provider missing %q, got %q", test.expectedProvider, result.provider)
					}
					if result.account != test.expectedAccount {
						t.Errorf("Account missing %q, got %q", test.expectedAccount, result.account)
					}
					if result.ticker == nil {
						t.Error("Ticker missing")
					}
					result.ticker.Stop() // Clean up
				}
			}
		})
	}
}

func TestGetProvider(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "otpauth://totp/Microsoft:user@example.com?secret=JBSWY3DPEHPK3PXP&issuer=Microsoft",
			expected: "Microsoft",
		},
		{
			input:    "otpauth://totp/Service:user@example.com?secret=JBSWY3DPEHPK3PXP",
			expected: "Service",
		},
		{
			input:    "not-down-with-otp",
			expected: "",
		},
		{
			input:    "",
			expected: "",
		},
	}

	for _, test := range tests {
		result := getProvider(test.input)
		if result != test.expected {
			t.Errorf("getProvider(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestMessageData(t *testing.T) {
	// Test that messageData implements tea.Model interface
	var _ interface {
		Init() tea.Cmd
		Update(tea.Msg) (tea.Model, tea.Cmd)
		View() string
	} = &messageData{}

	// Create a test message data
	msg := &messageData{
		provider:  "TestService",
		account:   "test@example.com",
		secret:    "JBSWY3DPEHPK3PXP",
		period:    30,
		code:      "123456",
		countdown: 15,
		ticker:    time.NewTicker(time.Second),
	}
	defer msg.ticker.Stop()

	// Test Init
	cmd := msg.Init()
	if cmd == nil {
		t.Error("Init should return a command")
	}

	// Test View
	view := msg.View()
	if view == "" {
		t.Error("View should not return empty string")
	}
	if !strings.Contains(view, "test@example.com") {
		t.Error("View should contain account name")
	}
	if !strings.Contains(view, "TestService") {
		t.Error("View should contain provider name")
	}
	if !strings.Contains(view, "123456") {
		t.Error("View should contain OTP code")
	}

	// Test Update with quit key
	model, cmd := msg.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if model == nil {
		t.Error("Update should return a model")
	}
	if cmd == nil {
		t.Error("Update should return tea.Quit command")
	}

	// Test Update with time tick
	model, cmd = msg.Update(time.Now())
	if model == nil {
		t.Error("Update should return a model")
	}
	if cmd == nil {
		t.Error("Update should return a command")
	}
}
