package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Header struct {
	Key   string
	Value string
}

type HeaderModel struct {
	keyInput   textinput.Model
	valueInput textinput.Model
	headers    []Header
	list       list.Model
	focusIndex int
	selected   int
	showList   bool
}

type errMsg error

func InitialHeaderModel() HeaderModel {
	keyInput := textinput.New()
	keyInput.Placeholder = "Header Key"
	keyInput.Focus()
	keyInput.CharLimit = 32
	keyInput.Width = 20

	valueInput := textinput.New()
	valueInput.Placeholder = "Header Value"
	valueInput.CharLimit = 64
	valueInput.Width = 30

	delegate := list.NewDefaultDelegate()
	l := list.New([]list.Item{}, delegate, 50, 10)
	l.Title = "Headers"

	return HeaderModel{
		keyInput:   keyInput,
		valueInput: valueInput,
		list:       l,
		selected:   -1,
	}
}

func (m HeaderModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m HeaderModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {

		// Submit header
		case msg.Type == tea.KeyEnter:
			if m.keyInput.Value() != "" && m.valueInput.Value() != "" {
				h := Header{Key: m.keyInput.Value(), Value: m.valueInput.Value()}
				m.headers = append(m.headers, h)
				m.list.InsertItem(len(m.headers)-1, headerItem(h))

				// reset inputs
				m.keyInput.SetValue("")
				m.valueInput.SetValue("")
				m.keyInput.Focus()
			}

		// Navigate list
		case msg.Type == tea.KeyDown:
			if len(m.headers) > 0 {
				m.showList = true
				m.list, _ = m.list.Update(msg)
			}

		// Delete selected header
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+d"))):
			if i := m.list.Index(); i >= 0 && i < len(m.headers) {
				m.headers = append(m.headers[:i], m.headers[i+1:]...)
				m.list.RemoveItem(i)
			}

		// Save
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+s"))):
			fmt.Println("Saving headers:", m.headers)
			return m, tea.Quit

		// Quit
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width - 4)
	}

	// Update text inputs
	var cmd tea.Cmd
	if m.keyInput.Focused() {
		m.keyInput, cmd = m.keyInput.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.valueInput, cmd = m.valueInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Always update list
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m HeaderModel) View() string {
	s := "Enter Header Key/Value\n\n"
	s += m.keyInput.View() + " : " + m.valueInput.View()
	s += "\n\nPress Enter to add, Ctrl+D to delete, Ctrl+S to save, Ctrl+C to quit.\n"

	if len(m.headers) > 0 {
		s += "\nHeaders:\n" + m.list.View()
	}

	return s
}

type headerItem Header

func (h headerItem) FilterValue() string { return h.Key }
func (h headerItem) Title() string       { return h.Key }
func (h headerItem) Description() string { return h.Value }
