package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BodyModel struct {
	textarea textarea.Model
	focused  bool
}

func InitialBodyModel() BodyModel {
	ta := textarea.New()
	ta.Placeholder = "{\n  \"key\": \"value\"\n}"
	ta.SetWidth(50)
	ta.SetHeight(10)
	ta.ShowLineNumbers = true
	ta.Focus()

	return BodyModel{
		textarea: ta,
		focused:  true,
	}
}

func (m BodyModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m BodyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

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
		}
	}

	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m BodyModel) View() string {
	title := lipgloss.NewStyle().Bold(false).Render("Request Body (JSON):\n")

	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1).
		Width(55).
		Height(20).
		Render(lipgloss.JoinVertical(lipgloss.Left, title, m.textarea.View()))

	return border
}
