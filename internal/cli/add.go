package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sahebg.dev/tick/internal/tickfile"
)

var addCmd = &cobra.Command{
	Use:   "add <task>",
	Short: "Add a new task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := tickfile.TickPath()

		input := args[0]
		task := tickfile.ParseAddInput(input)

		if project, _ := cmd.Flags().GetString("project"); project != "" {
			task.Project = project
		}
		if priority, _ := cmd.Flags().GetString("priority"); priority != "" {
			task.Priority = parsePriorityFlag(priority)
		}
		if durStr, _ := cmd.Flags().GetString("duration"); durStr != "" {
			d, err := tickfile.ParseDuration(durStr)
			if err != nil {
				return fmt.Errorf("invalid duration: %s", durStr)
			}
			task.Duration = d
		}

		line := tickfile.FormatTask(task)

		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("opening file: %w", err)
		}
		defer f.Close()

		fi, _ := f.Stat()
		if fi.Size() > 0 {
			f.WriteString("\n")
		}
		f.WriteString(line + "\n")

		fmt.Println(task.Title)
		return nil
	},
}

func init() {
	addCmd.Flags().String("project", "", "Project name")
	addCmd.Flags().String("priority", "", "Priority (p1, p2, p3)")
	addCmd.Flags().String("duration", "", "Duration (e.g. 2h, 30m, 1d, 1w)")
}
