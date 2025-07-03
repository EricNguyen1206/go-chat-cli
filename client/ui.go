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
		fmt.Println("❌ Kết nối thất bại:", err)
		return
	}

	m := model{
		client:   client,
		username: username,
		messages: []string{"✅ Đã kết nối đến server."},
		input:    "",
	}

	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Println("❌ Lỗi chạy chương trình:", err)
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

	// Hiển thị tối đa 20 tin nhắn gần nhất
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
		return <-c.recv
	}
}
