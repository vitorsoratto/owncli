package textinput

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var titleStyle = lipgloss.NewStyle().
	PaddingLeft(2).
	PaddingRight(2).
	Background(lipgloss.Color("#212121")).
	Foreground(lipgloss.Color("#d5a4dd")).
	Bold(true)

type (
	errMsg error
)

type Output struct {
	Output  string
	Answers []string
}

func (o *Output) AddAnswer(val string) {
	o.Answers = append(o.Answers, val)
}

func (o *Output) update(val string) {
	o.Output = val
}

type model struct {
	textInput textinput.Model
	err       error
	output    *Output
	header    string
}

func InitialTextInputModel(output *Output, header string) model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput: ti,
		err:       nil,
		output:    output,
		header:    titleStyle.Render(header),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			if len(m.output.Answers) == 2 {
				m.output.AddAnswer(m.textInput.Value())
				m.output.update("Done!")
				return m, tea.Quit
			}

			if len(m.output.Answers) < 3 {
				m.output.AddAnswer(m.textInput.Value())
				m.textInput.Reset()
			}

			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	msg := fmt.Sprintf(
		"%s\n\n%s\n\n",
		m.header,
		m.textInput.View(),
	)

	for i, v := range m.output.Answers {
		o := fmt.Sprintf("Answer %v: %s  -  Length:%v", i+1, v, len(m.output.Answers))
		msg += lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f904fd")).
			Bold(true).
			Render(o)

		msg += "\n"
	}

	o := fmt.Sprintf("\n%s\n", m.output.Output)
	msg += lipgloss.NewStyle().
		Foreground(lipgloss.Color("#008000")).
		Bold(true).
		Render(o)

	return msg
}
