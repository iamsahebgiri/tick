package cli

import (
	"fmt"

	"github.com/iamsahebgiri/tick/internal/tickfile"
)

func updateTaskStatus(taskIdx int, newStatus string, note string) error {
	path := tickfile.TickPath()

	tasks, err := tickfile.ParseFile(path)
	if err != nil {
		return err
	}

	if taskIdx < 1 || taskIdx > len(tasks) {
		return fmt.Errorf("task %d not found", taskIdx)
	}

	target := &tasks[taskIdx-1]
	fileLine := target.Line

	target.Status = newStatus
	if note != "" {
		target.Note = note
	}

	lines, err := tickfile.ReadLines(path)
	if err != nil {
		return err
	}

	newLine := tickfile.FormatTask(*target)
	lineIdx := fileLine - 1
	if lineIdx < 0 || lineIdx >= len(lines) {
		return fmt.Errorf("line %d out of range", fileLine)
	}
	lines[lineIdx] = newLine

	return tickfile.WriteLines(path, lines)
}
