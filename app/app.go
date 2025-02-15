package app

import (
	"encoding/base32"
	"fmt"
	"log"
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

// cleanString removes URL encoded chars from strings.
func cleanString(arg string) string {
	str, _ := url.QueryUnescape(arg)
	return str
}

// ArgParser captures a string or file path containing a URL.
func ArgParser(arg string) string {
	if _, err := os.Stat(arg); err == nil {
		data, err := os.ReadFile(arg)
		if err != nil {
			log.Fatalf("Unable to read known file path: %v", err)
		}
		str := strings.TrimSpace(string(data))
		return cleanString(str)
	}
	return cleanString(arg)
}

// UrlParser calls otp.NewKeyFromURL() and parses keys into messageData struct
func UrlParser(url string) (*messageData, error) {
	key, err := otp.NewKeyFromURL(url)
	if err != nil {
		return nil, err
	}
	secret := key.Secret()

	_, err = base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return nil, fmt.Errorf("secret is invalid: %v", err)
	}

	message := &messageData{
		provider: getProvider(url),
		account:  key.AccountName(),
		secret:   secret,
		period:   key.Period(),
		ticker:   time.NewTicker(time.Second),
	}

	return message, nil
}

// getProvider performs rudementary URL parsing and extracts a provider (if any)
func getProvider(url string) string {
	colon := strings.Split(url, ":")
	slash := strings.Split(colon[1], "/")
	return slash[3]
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
		log.Fatalf("Unable to retrieve code: %v", msg.err)
	}
	return m, nil
}

// View uses bubble View() to update terminal.
func (m messageData) View() string {

	yellow := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFDD00"))
	lime := lipgloss.NewStyle().Foreground(lipgloss.Color("#5CDE73"))
	provider := lipgloss.NewStyle().Bold(true).Render(m.provider)
	account := lime.Render(m.account)
	code := yellow.Bold(true).Render(m.code)
	count := yellow.Render(strconv.Itoa(m.countdown))
	const arrow = "\u2192"

	text := fmt.Sprintf(`%s %s %s

Token: %s

Regenerates in %s seconds

Press q to quit`, account, arrow, provider, code, count)

	return lipgloss.NewStyle().Padding(0, 1, 1).Render(text)
}

// GoTeaP creates a bubble tea program.
func Interactive(message *messageData) *tea.Program {
	return tea.NewProgram(message)
}
