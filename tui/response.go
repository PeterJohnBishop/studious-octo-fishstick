package tui

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	blurredStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	focusedButton = focusedStyle.Render("[ SEND ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("SEND"))
)

type ResponseModel struct {
	focused    bool
	loading    bool
	errMessage string
	response   string
	viewport   viewport.Model
}

func (m ResponseModel) Init() tea.Cmd {
	return nil
}

func InitialResponseModel() *ResponseModel {
	vp := viewport.New(237, 20)
	vp.YPosition = 0
	vp.HighPerformanceRendering = false
	vp.SetContent("")
	// add border + scrollbar style
	vp.Style = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)
	return &ResponseModel{
		focused:    false,
		response:   "",
		loading:    false,
		errMessage: "",
		viewport:   vp,
	}
}

func addLineNumbers(content string) string {
	lines := strings.Split(content, "\n")
	for i := range lines {
		lines[i] = fmt.Sprintf("%3d | %s", i+1, lines[i])
	}
	return strings.Join(lines, "\n")
}

func (m *ResponseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			return m, func() tea.Msg { return NextModelMsg{} }
		case "shift+tab":
			return m, func() tea.Msg { return PreviousModelMsg{} }
		case "enter":
			m.focused = true
			m.loading = true
			m.errMessage = ""
			m.response = ""
			m.viewport.SetContent("")
			return m, func() tea.Msg { return SendRequestMsg{} }
		}

		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd

	case ResponseMsg:
		// fmt.Println("DEBUG ResponseMsg received:", len(msg.response), "chars")
		m.loading = false
		m.response = msg.response
		m.errMessage = msg.err
		if msg.response != "" {
			withLines := addLineNumbers(msg.response)
			// fmt.Println("DEBUG setting viewport content with", len(strings.Split(withLines, "\n")), "lines")
			m.viewport.SetContent(withLines)
			m.viewport.GotoTop()
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
		b.WriteString("requesting...\n")
	} else if m.errMessage != "" {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render(m.errMessage))
		b.WriteString("\n")
	} else if m.response != "" {
		b.WriteString(m.viewport.View())
		b.WriteString("\n")
	}

	return b.String()
}
