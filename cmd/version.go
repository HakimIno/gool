package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "1.0.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version number of gool CLI tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gool version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
