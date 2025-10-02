package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MethodModel struct {
	methods  []string
	cursor   int
	selected map[int]struct{}
}

func InitialMethodModel() MethodModel {
	return MethodModel{
		methods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		selected: make(map[int]struct{}),
	}
}

func (m MethodModel) Init() tea.Cmd {
	return nil
}

func (m MethodModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.methods)-1 {
				m.cursor++
			}

		case "enter", " ":
			if _, ok := m.selected[m.cursor]; ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected = make(map[int]struct{})
				m.selected[m.cursor] = struct{}{}
			}
		case "tab":
			return m, func() tea.Msg {
				return NextModelMsg{}
			}
		case "Shift+tab":
			return m, func() tea.Msg {
				return PreviousModelMsg{}
			}
		}
	}

	return m, nil
}

func (m MethodModel) View() string {

	s := "Method:\n\n"

	for i, choice := range m.methods {

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1).
		Width(55).
		Height(20).
		Render(s)

	return border
}
