package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Root struct {
	method  MethodModel
	headers HeaderModel
	active  int
}

type PreviousModelMsg struct{}
type NextModelMsg struct{}
type SaveMethodMsg struct{}
type SaveHeadersMsg struct{}
type SaveEndpointMsg struct{}
type SaveBodyMsg struct{}
type SaveResponseMsg struct{}
type ResetMsg struct{}

func InitialRootModel() Root {
	return Root{
		method:  InitialMethodModel(),
		headers: InitialHeaderModel(),
		active:  0,
	}
}

func (m Root) Init() tea.Cmd {
	return nil
}

func (m Root) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case ">":
			if m.active < 1 {
				m.active++
			}
			return m, nil
		case "<":
			if m.active > 0 {
				m.active--
			}
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		}

	case PreviousModelMsg:
		if m.active > 0 {
			m.active--
		}
	case NextModelMsg:
		if m.active < 1 {
			m.active++
		}
	}

	var updated tea.Model
	var cmd tea.Cmd

	switch m.active {
	case 0: // endpoint (method)
		updated, cmd = m.method.Update(msg)
		m.method = updated.(MethodModel)
	case 1: // headers
		updated, cmd = m.headers.Update(msg)
		m.headers = updated.(HeaderModel)
	}

	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Root) View() string {
	top := "\nWelcome to the Postgres TUI Client!\n\n"
	// Row of children
	top += lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.method.View(),
		m.headers.View(),
	)

	// Bottom child below
	col := lipgloss.JoinVertical(lipgloss.Left,
		top,
		"CTRL+C to quit",
	)

	return col
}
