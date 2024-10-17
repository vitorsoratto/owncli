package csvtodb

import (
	"os"

	"owncli/cmd/ui/csvtodb/filepicker"
	"owncli/cmd/ui/csvtodb/reader"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	logoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#400080")).
			Bold(true)

	tipMsgStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			Foreground(lipgloss.Color("#408444")).
			Italic(true)

	endingMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#008000")).
			Bold(true)
)

var CsvtodbCmd = &cobra.Command{
	Use:                   "csvtodb",
	DisableFlagsInUseLine: true,
	Short:                 "Uses a csv file and insert it into a selected database",
	Long:                  "Uses a csv file and insert it into a selected database",
	Run: func(cmd *cobra.Command, args []string) {
		var options filepicker.FilePickerOptions

		options.AllowedTypes = []string{".db", ".csv"}
		options.CurrentDirectory, _ = os.UserHomeDir()
		options.Output = &filepicker.Output{}

		tea.NewProgram(filepicker.InitialFilePicker(&options)).Run()

		csvPath := options.Output.SelectedCsvFile
		dbPath := options.Output.SelectedDBFile

		tea.NewProgram(reader.InitialReaderModel(csvPath, dbPath)).Run()
	},
}
