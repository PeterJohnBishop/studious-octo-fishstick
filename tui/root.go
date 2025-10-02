package tui

import (
	"studious-octo-fishstick/api"

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
	response *ResponseModel
	active   int
}

type PreviousModelMsg struct{}
type NextModelMsg struct{}
type SendRequestMsg struct{}
type ResetMsg struct{}
type ResponseMsg struct {
	response string
	err      string
}

func doRequest(method string, headers []api.Header, endpoint string, params []api.Params, body string) tea.Cmd {
	return func() tea.Msg {
		resp, err := api.SendRequest(method, headers, endpoint, params, &body)
		if err != nil {
			return ResponseMsg{response: "", err: err.Error()}
		}

		pretty := api.PrettyPrintJSON(resp)
		// fmt.Println("DEBUG pretty:", pretty) // 👈 confirm we actually have JSON

		return ResponseMsg{response: pretty, err: ""}
	}
}

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
	case SendRequestMsg:
		var method string = m.method.methods[m.method.cursor]
		var headers []api.Header = m.headers.headers
		var endpoint string = m.endpoint.endpoint.Value()
		var params []api.Params = m.endpoint.params
		var body string = m.body.textarea.Value()

		m.active = 4
		m.response.errMessage = ""
		m.response.response = ""
		m.response.loading = true
		return m, doRequest(method, headers, endpoint, params, body)

	case ResponseMsg:
		updated, cmd := m.response.Update(msg)
		m.response = updated.(*ResponseModel)
		return m, cmd
	case ResetMsg:
		m.method = InitialMethodModel()
		m.headers = InitialHeaderModel()
		m.endpoint = InitialEndpointModel()
		m.body = InitialBodyModel()
		m.response = InitialResponseModel()
		m.active = 0
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
		m.response = updated.(*ResponseModel)
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
