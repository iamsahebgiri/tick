package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done <line>",
	Short: "Mark a task as done",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		lineNum, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid line number: %s", args[0])
		}

		note, _ := cmd.Flags().GetString("note")

		return updateTaskStatus(lineNum, "x", note)
	},
}

func init() {
	doneCmd.Flags().String("note", "", "Completion note")
}
