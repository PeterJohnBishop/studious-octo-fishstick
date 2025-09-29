package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type MethodModel struct {
	methods  []string
	cursor   int
	selected map[int]struct{}
}

func InitialMethodModel() MethodModel {
	return MethodModel{
		methods:  []string{"GET", "POST", "PUT", "DELETE"},
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

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.methods)-1 {
				m.cursor++
			}

		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
				// send SaveMethodMsg
			}
		}
	}

	return m, nil
}

func (m MethodModel) View() string {

	s := ""

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

	return s
}
