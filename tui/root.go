package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	activeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true) // pink + bold
	inactiveStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))            // gray
)

type Root struct {
	method   MethodModel
	headers  HeaderModel
	endpoint EndpointModel
	body     BodyModel
	response ResponseModel
	active   int
}

type PreviousModelMsg struct{}
type NextModelMsg struct{}
type SendRequestMsg struct{}
type ResetMsg struct{}

func InitialRootModel() Root {
	return Root{
		method:   InitialMethodModel(),
		headers:  InitialHeaderModel(),
		endpoint: InitialEndpointModel(),
		body:     InitialBodyModel(),
		response: InitialResponseModel(),
		active:   0,
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
		case "ctrl+c":
			return m, tea.Quit
		}

	case PreviousModelMsg:
		m.active = (m.active + 4) % 5
		return m, nil
	case NextModelMsg:
		m.active = (m.active + 1) % 5
		return m, nil

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
	case 2: // endpoint
		updated, cmd = m.endpoint.Update(msg)
		m.endpoint = updated.(EndpointModel)
	case 3: // body
		updated, cmd = m.body.Update(msg)
		m.body = updated.(BodyModel)
	case 4: // response
		updated, cmd = m.response.Update(msg)
		m.response = updated.(ResponseModel)
	}

	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Root) View() string {

	methodView := m.method.View()
	headersView := m.headers.View()
	endpointView := m.endpoint.View()
	bodyView := m.body.View()
	responseView := m.response.View()

	if m.active == 0 {
		methodView = activeStyle.Render(methodView)
	} else {
		methodView = inactiveStyle.Render(methodView)
	}

	if m.active == 1 {
		headersView = activeStyle.Render(headersView)
	} else {
		headersView = inactiveStyle.Render(headersView)
	}

	if m.active == 2 {
		endpointView = activeStyle.Render(endpointView)
	} else {
		endpointView = inactiveStyle.Render(endpointView)
	}

	if m.active == 3 {
		bodyView = activeStyle.Render(bodyView)
	} else {
		bodyView = inactiveStyle.Render(bodyView)
	}

	if m.active == 4 {
		responseView = activeStyle.Render(responseView)
	} else {
		responseView = inactiveStyle.Render(responseView)
	}

	rowTop := "\n\n"
	sepH := lipgloss.NewStyle().Width(4).Render("")
	rowTop += lipgloss.JoinHorizontal(
		lipgloss.Top,
		methodView,
		sepH,
		headersView,
		sepH,
		endpointView,
		sepH,
		bodyView,
	)

	rowBottom := "\n\n"
	rowBottom += lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		sepH,
		responseView,
	)

	sepV := lipgloss.NewStyle().Height(1).Render("\n")
	col := lipgloss.JoinVertical(
		lipgloss.Left,
		rowTop,
		rowBottom,
		sepV,
		"[ Use shift+tab (back) and tab (forward) to navigate between sections. CTRL+C to quit ]",
		sepV,
	)

	return col
}
