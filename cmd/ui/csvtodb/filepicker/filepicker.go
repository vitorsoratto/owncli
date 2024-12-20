package filepicker

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	selectStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			PaddingRight(2).
			Background(lipgloss.Color("#212121")).
			Foreground(lipgloss.Color("#d5a4dd")).
			Bold(true)

	selectedStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			PaddingRight(2).
			Background(lipgloss.Color("#212121")).
			Foreground(lipgloss.Color("#008000")).
			Italic(true)

	fileStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			PaddingRight(2).
			Background(lipgloss.Color("#212121")).
			Foreground(lipgloss.Color("#008000")).
			Italic(true)
)

type (
	errMsg error
)

type Output struct {
	SelectedCsvFile string
	SelectedDBFile  string
}

type FilePickerOptions struct {
	AllowedTypes     []string
	CurrentDirectory string
	Output           *Output
}

type model struct {
	filePicker filepicker.Model
	err        error
	output     *Output
	quit       bool
}

func InitialFilePicker(filePickerOptions *FilePickerOptions) model {
	fp := filepicker.New()
	fp.AutoHeight = false
	fp.Height = 20
	fp.AllowedTypes = filePickerOptions.AllowedTypes
	fp.CurrentDirectory = filePickerOptions.CurrentDirectory

	return model{
		filePicker: fp,
		err:        nil,
		output:     filePickerOptions.Output,
	}
}

func (m model) Init() tea.Cmd {
	return m.filePicker.Init()
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quit = true
			return m, tea.Quit
		}

	case clearErrorMsg:
		m.err = nil
	}

	var cmd tea.Cmd
	m.filePicker, cmd = m.filePicker.Update(msg)

	if didSelect, path := m.filePicker.DidSelectFile(msg); didSelect {
		if (filepath.Ext(path) == ".csv" && m.output.SelectedCsvFile != "") || (filepath.Ext(path) == ".db" && m.output.SelectedDBFile != "") {
			m.err = errors.New(path + " is not a valid file.")

			return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
		}

		if m.output.SelectedCsvFile != "" && filepath.Ext(path) != ".db" {
			m.err = errors.New(path + " is not a .db file.")

			return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
		}

		if m.output.SelectedDBFile != "" && filepath.Ext(path) != ".csv" {
			m.err = errors.New(path + " is not a .csv file.")

			return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
		}

		if filepath.Ext(path) == ".csv" && m.output.SelectedCsvFile == "" {
			m.output.SelectedCsvFile = path
			m.filePicker.AllowedTypes = []string{".db"}
		} else if filepath.Ext(path) == ".db" && m.output.SelectedDBFile == "" {
			m.output.SelectedDBFile = path
		}

		if (m.output.SelectedCsvFile != "" && m.output.SelectedDBFile != "") {
			m.quit = true
			return m, tea.Quit
		}
	}
	return m, cmd
}

func (m model) View() string {
	var s strings.Builder
	s.WriteString("\n")
	if (m.quit) {
		return s.String();
	}
	if m.err != nil {
		s.WriteString(m.filePicker.Styles.DisabledFile.Render(m.err.Error()))
		s.WriteString("\n")
	} else {
		csvFile := m.output.SelectedCsvFile
		dbFile := m.output.SelectedDBFile

		renderSelectedFile := func(file, label string) {
			if file != "" {
				s.WriteString(selectedStyle.Render(fmt.Sprintf("Selected %s:", label)))
				s.WriteString(fileStyle.Render(file) + "\n")
			} else {
				s.WriteString(selectStyle.Render(fmt.Sprintf("Select the %s file", label)) + "\n")
			}
		}

		renderSelectedFile(csvFile, "csv")

		if csvFile != "" && dbFile == "" {
			s.WriteString(selectStyle.Render("Select the database file") + "\n")
		} else if dbFile != "" {
			renderSelectedFile(dbFile, "db")
		}
	}
	s.WriteString("\n" + m.filePicker.View())
	s.WriteString("\n press esc or ctrl+c to exit")
	return s.String()
}
