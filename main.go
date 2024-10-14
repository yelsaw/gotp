package main

import (
	"encoding/base32"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type messageData struct {
	secret    string
	period    uint64
	code      string
	countdown int
	ticker    *time.Ticker
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("First arg requires: %s <full-totp-url>", os.Args[0])
	}

	url := os.Args[1]

	message, err := urlParser(url)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tea.NewProgram(message).Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// Calls otp.NewKeyFromURL() and parses keys into messageData struct
func urlParser(url string) (*messageData, error) {
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
		secret: secret,
		period: key.Period(),
		ticker: time.NewTicker(time.Second),
	}

	return message, nil
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

// Bubble Tea: View()
func (m messageData) View() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFDD00")).Bold(true)
	code := style.Render(m.code)
	count := style.Render(strconv.Itoa(m.countdown))
	text := fmt.Sprintf("\n\n\nToken: %s\n\nExpires in %s seconds\n\nPress q to quit\n\n", code, count)
	return text
}
