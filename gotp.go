package gotp

import (
	"encoding/base32"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

var (
	// accentColor #hex value for accent (default yellow).
	accentColor = "#FFDD00"
	// emailColor #hex value for email (default green).
	emailColor = "#5CDE73"
	// textMessage defines terminal output.
	textMessage = `%s %s %s

Token: %s

Regenerates in %s seconds

Press q to quit`
)

// messageData struct stores bubble tea and parsed url data.
type messageData struct {
	provider  string
	secret    string
	account   string
	period    uint64
	code      string
	countdown int
	ticker    *time.Ticker
}

// codeMsg is used in bubble tea Update() to display an OTP code
// with additional details.
type codeMsg struct {
	code string
}

// errMsg is used in bubble tea Update() to display an error.
type errMsg struct {
	err error
}

// messageTheme struct contains display placeholders used in View().
type messageTheme struct {
	accent   lipgloss.Style
	email    lipgloss.Style
	provider string
	account  string
	code     string
	count    string
	arrow    string
}

func initTheme(m messageData) messageTheme {
	t := &messageTheme{
		accent:   lipgloss.NewStyle().Foreground(lipgloss.Color(accentColor)),
		email:    lipgloss.NewStyle().Foreground(lipgloss.Color(emailColor)),
		provider: lipgloss.NewStyle().Bold(true).Render(m.provider),
		arrow:    "\u2192",
	}

	t.code = t.accent.Bold(true).Render(m.code)
	t.count = t.accent.Render(strconv.Itoa(m.countdown))
	t.account = t.email.Render(m.account)
	return *t
}

// cleanString removes URL encoded chars from strings and validates input.
func cleanString(arg string) string {
	str, _ := url.QueryUnescape(arg)
	return strings.TrimSpace(str)
}

// ArgParser captures a string or file path containing a URL.
func ArgParser(arg string) (string, error) {
	if arg == "" {
		return "", fmt.Errorf("empty argument provided")
	}

	// Check if argument is a file path
	if stat, err := os.Stat(arg); err == nil && !stat.IsDir() {
		data, err := os.ReadFile(arg)
		if err != nil {
			return "", fmt.Errorf("unable to read file: %v", err)
		}
		return cleanString(string(data)), nil
	} else if err != nil && !os.IsNotExist(err) {
		// If there's an error other than file not existing, return it
		return "", fmt.Errorf("error accessing file: %v", err)
	}

	return cleanString(arg), nil
}

// UrlParser calls otp.NewKeyFromURL() and parses keys into messageData struct
func UrlParser(rawURL string) (*messageData, error) {
	if rawURL == "" {
		return nil, fmt.Errorf("empty URL provided")
	}

	if !strings.HasPrefix(rawURL, "otpauth://") {
		return nil, fmt.Errorf("invalid OTP URL format, must start with otpauth://")
	}

	key, err := otp.NewKeyFromURL(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OTP URL: %v", err)
	}
	secret := key.Secret()

	// Validate secret using Base32 (correct encoding for TOTP)
	_, err = base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(secret))
	if err != nil {
		return nil, fmt.Errorf("secret is invalid (must be Base32 encoded): %v", err)
	}

	message := &messageData{
		provider: getProvider(rawURL),
		account:  key.AccountName(),
		secret:   secret,
		period:   key.Period(),
		ticker:   time.NewTicker(time.Second),
	}

	return message, nil
}

// getProvider extracts the provider/issuer from the OTP URL using proper URL parsing.
func getProvider(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	// Try to get issuer from query parameter first
	if issuer := u.Query().Get("issuer"); issuer != "" {
		return issuer
	}

	// Fallback to path extraction for backward compatibility
	pathParts := strings.Split(u.Path, ":")
	if len(pathParts) >= 2 {
		return strings.TrimPrefix(pathParts[0], "/")
	}

	return ""
}

// getCode generates a time-based code.
func getCode(secret string) tea.Cmd {
	return func() tea.Msg {
		code, err := totp.GenerateCode(secret, time.Now())
		if err != nil {
			return errMsg{err}
		}
		return codeMsg{code}
	}
}

// tickCmd is used in bubble tea Init() and Update().
func tickCmd(ticker *time.Ticker) tea.Cmd {
	return func() tea.Msg {
		return <-ticker.C
	}
}

// Init initializes bubble tea Batch() with tickCmd() and getCode().
func (m messageData) Init() tea.Cmd {
	return tea.Batch(tickCmd(m.ticker), getCode(m.secret))
}

// Update uses bubble Update() to display messages.
func (m messageData) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			m.ticker.Stop()
			return m, tea.Quit
		}
	case codeMsg:
		m.code = msg.code
		m.countdown = int(m.period)
		return m, tickCmd(m.ticker)
	case time.Time:
		m.countdown--
		if m.countdown <= 0 {
			return m, getCode(m.secret)
		}
		return m, tickCmd(m.ticker)
	case errMsg:
		m.ticker.Stop()
		return m, tea.Quit
	}
	return m, nil
}

// View uses bubble View() to update terminal.
func (m messageData) View() string {
	t := initTheme(m)
	text := fmt.Sprintf(textMessage, t.account, t.arrow, t.provider, t.code, t.count)
	return lipgloss.NewStyle().Padding(0, 1, 1).Render(text)
}

// Interactive creates a bubble tea program.
func Interactive(message *messageData) *tea.Program {
	return tea.NewProgram(message)
}
