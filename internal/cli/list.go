package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"sahebg.dev/tick/internal/tickfile"
)

func renderStatus(s string) string {
	switch s {
	case "-":
		return statusTodo.Render(s)
	case "/":
		return statusActive.Render(s)
	case "x":
		return statusDone.Render(s)
	case "~":
		return statusDropped.Render(s)
	}
	return s
}

func runList(status, project, tag, priority string) error {
	path := tickfile.TickPath()

	tasks, err := tickfile.ParseFile(path)
	if err != nil {
		return err
	}

	for i, t := range tasks {
		if status != "" && t.Status != status {
			continue
		}
		if project != "" && t.Project != project {
			continue
		}
		if tag != "" && !contains(t.Tags, tag) {
			continue
		}
		if priority != "" && t.Priority != parsePriorityFlag(priority) {
			continue
		}

		fmt.Printf("%s  %s  %s", lineNumStyle.Render(fmt.Sprintf("%3d", i+1)), renderStatus(t.Status), t.Title)
		if t.Project != "" {
			fmt.Printf("  %s", projectStyle.Render("+"+t.Project))
		}
		if t.Duration > 0 {
			fmt.Printf("  %s", durationStyle.Render(tickfile.FormatDuration(t.Duration)))
		}
		if t.Priority > 0 {
			s := fmt.Sprintf("p%d", t.Priority)
			switch t.Priority {
			case 1:
				fmt.Printf("  %s", priorityP1.Render(s))
			case 2:
				fmt.Printf("  %s", priorityP2.Render(s))
			case 3:
				fmt.Printf("  %s", priorityP3.Render(s))
			}
		}
		if t.Note != "" {
			fmt.Printf("  %s", noteStyle.Render("; "+t.Note))
		}
		fmt.Println()
	}

	return nil
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		status, _ := cmd.Flags().GetString("status")
		project, _ := cmd.Flags().GetString("project")
		tag, _ := cmd.Flags().GetString("tag")
		priority, _ := cmd.Flags().GetString("priority")
		return runList(status, project, tag, priority)
	},
}

func init() {
	listCmd.Flags().String("status", "", "Filter by status")
	listCmd.Flags().String("project", "", "Filter by project")
	listCmd.Flags().String("tag", "", "Filter by tag")
	listCmd.Flags().String("priority", "", "Filter by priority (p1, p2, p3)")
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func parsePriorityFlag(s string) int {
	switch s {
	case "p1":
		return 1
	case "p2":
		return 2
	case "p3":
		return 3
	default:
		return 0
	}
}
