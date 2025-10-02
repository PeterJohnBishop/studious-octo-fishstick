package tui

import (
	"studious-octo-fishstick/api"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type EndpointModel struct {
	endpoint textinput.Model
	selected int
	params   []api.Params
}

func InitialEndpointModel() EndpointModel {
	endpoint := textinput.New()
	endpoint.Placeholder = "https://api.example.com/resource"
	endpoint.Focus()
	endpoint.CharLimit = 256
	endpoint.Width = 256

	return EndpointModel{
		endpoint: endpoint,
		selected: -1,
	}
}

func (m EndpointModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m EndpointModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			params, err := api.ExtractQueryParams(m.endpoint.Value())
			if err == nil {
				m.params = params
			}
			return m, nil
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

	var cmd tea.Cmd
	m.endpoint, cmd = m.endpoint.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m EndpointModel) View() string {
	s := "Endpoint:\n\n"
	s += m.endpoint.View() + "\n\n"

	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1).
		Width(55).
		Height(20).
		Render(s)
	return border
}
