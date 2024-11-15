package fbadm

import (
	"fmt"

	"owncli/cmd/ui/fbadm/textinput"

	// "firebase.google.com/go/v4@latest"

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

var FBAdmCmd = &cobra.Command{
	Use:                   "fbadm",
	DisableFlagsInUseLine: true,
	Short:                 "Usages for the Firebase Admin",
	Long:                  "",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := tea.NewProgram(textinput.InitialModel()).Run()
		if err != nil {
			fmt.Println(err)
		}
	},
}
