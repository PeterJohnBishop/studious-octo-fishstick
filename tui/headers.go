package tui

import (
	"studious-octo-fishstick/api"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HeaderModel struct {
	keyInput   textinput.Model
	valueInput textinput.Model
	headers    []api.Header
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
		switch msg.String() {
		case "enter":
			if m.keyInput.Focused() {
				m.keyInput.Blur()
				m.valueInput.Focus()
			} else {
				m.valueInput.Blur()
				m.keyInput.Focus()
			}
		case "tab":
			return m, func() tea.Msg {
				return NextModelMsg{}
			}
		case "shift+tab":
			return m, func() tea.Msg {
				return PreviousModelMsg{}
			}
		}
		switch {
		case msg.Type == tea.KeyEnter:
			if m.keyInput.Value() != "" && m.valueInput.Value() != "" {
				h := api.Header{Key: m.keyInput.Value(), Value: m.valueInput.Value()}
				m.headers = append(m.headers, h)
				m.list.InsertItem(len(m.headers)-1, headerItem(h))

				// reset inputs
				m.keyInput.SetValue("")
				m.valueInput.SetValue("")
				m.keyInput.Focus()
				m.valueInput.Blur()
			}

		case msg.Type == tea.KeyDown:
			if len(m.headers) > 0 {
				m.showList = true
				m.list, _ = m.list.Update(msg)
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+d"))):
			if i := m.list.Index(); i >= 0 && i < len(m.headers) {
				m.headers = append(m.headers[:i], m.headers[i+1:]...)
				m.list.RemoveItem(i)
			}
		}

	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width - 4)
	}

	// Update text inputs
	var cmd tea.Cmd
	m.keyInput, cmd = m.keyInput.Update(msg)
	cmds = append(cmds, cmd)

	m.valueInput, cmd = m.valueInput.Update(msg)
	cmds = append(cmds, cmd)

	// Always update list
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m HeaderModel) View() string {
	s := "Headers:\n\n"
	s += m.keyInput.View() + " : " + m.valueInput.View()
	s += "\n\n[ Press Enter to add, Ctrl+D to delete header ].\n"

	if len(m.headers) > 0 {
		s += "\nHeaders:\n" + m.list.View()
	}

	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1).
		Width(55).
		Height(20).
		Render(s)

	return border
}

type headerItem api.Header

func (h headerItem) FilterValue() string { return h.Key }
func (h headerItem) Title() string       { return h.Key }
func (h headerItem) Description() string { return h.Value }
