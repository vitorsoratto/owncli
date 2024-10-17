package reader

import (
	"strings"

	"owncli/cmd/csvtodb/schema"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type readCsvMsg struct{}

type readingCSVMsg struct {
	plans     map[int]schema.Plan
	workouts  map[int]schema.Workout
	exercises []schema.Exercise
}

type dataInsertedMsg struct{}

type Model struct {
	err     error
	csvPath string
	output  string
	spinner spinner.Model
	quit    bool
}

func InitialReaderModel(csvPath string, dbPath string) Model {
	schema.InitDB(dbPath)
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#d5a4dd"))

	return Model{
		err:     nil,
		csvPath: csvPath,
		spinner: s,
	}
}

func (m Model) Read() tea.Cmd {
	return func() tea.Msg {
		plans, workouts, exercises, err := schema.ReadCSV(m.csvPath)
		if err != nil {
			m.err = err
			return nil
		}

		return readingCSVMsg{
			plans:     plans,
			workouts:  workouts,
			exercises: exercises,
		}
	}
}

func Insert(plans map[int]schema.Plan, workouts map[int]schema.Workout, exercises []schema.Exercise) tea.Cmd {
	return func() tea.Msg {
		schema.InsertPlans(plans)
		schema.InsertWorkouts(workouts)
		schema.InsertExercises(exercises)

		return dataInsertedMsg{}
	}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		return readCsvMsg{}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case readCsvMsg:
		m.output = "Reading csv file"
		return m, tea.Batch(m.spinner.Tick, m.Read())

	case readingCSVMsg:
		m.output = "Inserting data"
		return m, tea.Batch(m.spinner.Tick, Insert(msg.plans, msg.workouts, msg.exercises))

	case dataInsertedMsg:
		m.output = "Data inserted on the database!"
		m.quit = true
		return m, tea.Quit

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	var str strings.Builder

	str.WriteString("\n")
	str.WriteString(m.output)

	if m.quit {
		return str.String()
	}

	str.WriteString(" ")
	str.WriteString(m.spinner.View())

	return str.String()
}
