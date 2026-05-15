package tickfile

import (
	"fmt"
	"os"
	"strings"

	"github.com/iamsahebgiri/tick/internal/model"
)

func FormatTask(task model.Task) string {
	var metaParts []string

	if task.Duration > 0 {
		metaParts = append(metaParts, formatDuration(task.Duration))
	}
	if task.Project != "" {
		metaParts = append(metaParts, "+"+task.Project)
	}
	if task.Priority > 0 {
		metaParts = append(metaParts, fmt.Sprintf("p%d", task.Priority))
	}
	for _, tag := range task.Tags {
		metaParts = append(metaParts, "#"+tag)
	}
	for _, m := range task.Mentions {
		metaParts = append(metaParts, "@"+m)
	}

	meta := strings.Join(metaParts, " ")

	var sb strings.Builder
	sb.WriteString(task.Status)
	if meta != "" {
		sb.WriteString(" " + meta)
	}
	sb.WriteString(" :: ")
	sb.WriteString(task.Title)
	if task.Note != "" {
		sb.WriteString(" ; " + task.Note)
	}
	return sb.String()
}

func ReadLines(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	content := strings.TrimRight(string(data), "\n")
	if content == "" {
		return nil, nil
	}
	return strings.Split(content, "\n"), nil
}

func WriteLines(path string, lines []string) error {
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}

func TickPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "tasks.tick"
	}
	return home + "/tasks.tick"
}
