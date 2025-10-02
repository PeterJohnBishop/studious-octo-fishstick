package tui

import (
	"bytes"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	focusedButton = focusedStyle.Render("[ SEND ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("SEND"))
)

type ResponseModel struct {
	focused    bool
	response   string
	loading    bool
	errMessage string
}

func (m ResponseModel) Init() tea.Cmd {
	return nil
}

func InitialResponseModel() ResponseModel {
	return ResponseModel{
		focused:    false,
		response:   "",
		loading:    false,
		errMessage: "",
	}
}

func (m ResponseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			return m, func() tea.Msg {
				return NextModelMsg{}
			}
		case "shift+tab":
			return m, func() tea.Msg {
				return PreviousModelMsg{}
			}
		case "enter":
			if m.focused {
				m.loading = true
				m.errMessage = ""
				return m, func() tea.Msg { return SendRequestMsg{} }
			}
		}
	}
	return m, nil
}

func (m ResponseModel) View() string {
	var b bytes.Buffer

	if m.focused {
		b.WriteString(focusedButton)
	} else {
		b.WriteString(blurredButton)
	}
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString("Loading...\n")
	} else if m.errMessage != "" {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render(m.errMessage))
		b.WriteString("\n")
	} else if m.response != "" {
		jsonBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1).
			Width(60).
			Render(m.response)
		b.WriteString(jsonBox)
		b.WriteString("\n")
	}

	return b.String()
}
