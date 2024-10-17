package reader

import (
	"owncli/cmd/csvtodb/schema"

	tea "github.com/charmbracelet/bubbletea"
)

type readCsvMsg struct{}

type readingCSVMsg struct {
	plans     map[int]schema.Plan
	workouts  map[int]schema.Workout
	exercises []schema.Exercise
}

type dataInsertedMsg struct{}

type model struct {
	err     error
	csvPath string
	output  string
}

func InitialReaderModel(csvPath string, dbPath string) model {
	schema.InitDB(dbPath)
	return model{
		err:     nil,
		csvPath: csvPath,
	}
}

func (m model) Read() tea.Cmd {
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

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		return readCsvMsg{}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case readCsvMsg:
		m.output = "Reading csv file..."
		return m, m.Read()

	case readingCSVMsg:
		m.output = "Inserting data..."
		return m, Insert(msg.plans, msg.workouts, msg.exercises)

	case dataInsertedMsg:
		m.output = "Data inserted on the database!"
		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() string {
	return m.output
}
