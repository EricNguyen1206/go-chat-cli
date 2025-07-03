package client

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	client   *WSClient
	username string
	messages []string
	input    string
}

func StartClientUI(wsURL string, username string) {
	client, err := NewWSClient(wsURL)
	if err != nil {
		fmt.Println("❌ Connection failed:", err)
		return
	}

	m := model{
		client:   client,
		username: username,
		messages: []string{"✅ Connected to server."},
		input:    "",
	}

	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Println("❌ Error running program:", err)
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(waitForMessage(m.client), tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyEnter:
			if m.input != "" {
				m.client.send <- m.input
				m.input = ""
			}
			return m, nil

		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
			return m, nil

		default:
			m.input += msg.String()
			return m, nil
		}

	case string:
		m.messages = append(m.messages, msg)
		return m, waitForMessage(m.client)

	default:
		return m, nil
	}
}

func (m model) View() string {
	output := "💬 Chat Room:\n\n"

	// Display the last 20 messages
	start := 0
	if len(m.messages) > 20 {
		start = len(m.messages) - 20
	}
	for _, msg := range m.messages[start:] {
		output += msg + "\n"
	}

	output += "\n👉 " + m.input

	return output
}

func waitForMessage(c *WSClient) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-c.recv
		if !ok {
			return tea.Quit() // Channel is closed → exit program
		}
		return msg
	}
}


