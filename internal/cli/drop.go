package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var dropCmd = &cobra.Command{
	Use:   "drop <line>",
	Short: "Soft-delete a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		lineNum, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid line number: %s", args[0])
		}
		return updateTaskStatus(lineNum, "~", "")
	},
}
