package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pquerna/otp/totp"
)

type messageData struct {
	secret    string
	period    int
	code      string
	countdown int
	ticker    *time.Ticker
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("First arg requires: %s <full-totp-url>", os.Args[0])
	}

	otpUrl := os.Args[1]
	u, err := url.Parse(otpUrl)

	if err != nil {
		log.Fatalf("Unable to parse 'totp' url: %v", err)
	}

	secret := u.Query().Get("secret")
	if secret == "" {
		log.Fatalf("Missing 'secret' from url")
	}

	periodParam := u.Query().Get("period")
	period := 30
	if periodParam != "" {
		period, err = strconv.Atoi(periodParam)
		if err != nil {
			log.Fatalf("Invalid 'period' parameter: %v", err)
		}
	}

	if _, err = tea.NewProgram(messageData{
		secret:    secret,
		period:    period,
		countdown: period,
		ticker:    time.NewTicker(time.Second),
	}).Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}

func getCode(secret string) tea.Cmd {
	return func() tea.Msg {
		code, err := totp.GenerateCode(secret, time.Now())
		if err != nil {
			return errMsg{err}
		}
		return codeMsg{code}
	}
}

type codeMsg struct {
	code string
}

type errMsg struct {
	err error
}

func tickCmd(ticker *time.Ticker) tea.Cmd {
	return func() tea.Msg {
		return <-ticker.C
	}
}

// Bubble Tea: Init()
func (m messageData) Init() tea.Cmd {
	return tea.Batch(tickCmd(m.ticker), getCode(m.secret))
}

// Bubble Tea: Update()
func (m messageData) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	case codeMsg:
		m.code = msg.code
		m.countdown = m.period
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

// Bubble Tea: View()
func (m messageData) View() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFDD00")).Bold(true)
	code := style.Render(m.code)
	count := style.Render(strconv.Itoa(m.countdown))
	text := fmt.Sprintf("\n\n\nToken: %s\n\nExpires in %s seconds\n\nPress q to quit\n\n", code, count)
	return text
}
