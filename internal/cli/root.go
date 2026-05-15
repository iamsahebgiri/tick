package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:     "tick",
	Version: "0.0.1",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runList("", "", "", "")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(doneCmd)
	rootCmd.AddCommand(activeCmd)
	rootCmd.AddCommand(dropCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
