package cmd

import (
	"os"

	"owncli/cmd/csvtodb"
	"owncli/cmd/fbadm"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "owncli",
	Short: "A cli created to solve my own problems",
	Long:  "A cli created to solve my own problems",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(csvtodb.CsvtodbCmd)
	rootCmd.AddCommand(fbadm.FBAdmCmd)
}
