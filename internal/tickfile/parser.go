package tickfile

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/iamsahebgiri/tick/internal/model"
)

var (
	statusRegex = regexp.MustCompile(`^([-/x~])`)
	priorityRe  = regexp.MustCompile(`\b(p[1-3])\b`)
	durationRe  = regexp.MustCompile(`=(\d+[mhdw])`)
	projectRe   = regexp.MustCompile(`\+([A-Za-z0-9-]+)`)
	tagRe       = regexp.MustCompile(`#(\w+)`)
	mentionRe   = regexp.MustCompile(`@(\w+)`)
	urlRe       = regexp.MustCompile(`https?://[^\s]+`)
)

func ParseFile(path string) ([]model.Task, error) {
	data, err := readFile(path)
	if err != nil {
		return nil, err
	}

	var tasks []model.Task
	lines := strings.Split(data, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}

		task, err := parseLine(line, i+1)
		if err != nil {
			continue
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func parseLine(line string, lineNum int) (model.Task, error) {
	task := model.Task{Line: lineNum}

	if match := statusRegex.FindStringSubmatch(line); match != nil {
		task.Status = match[1]
		line = line[1:]
	} else {
		return task, nil
	}

	if idx := strings.Index(line, "::"); idx != -1 {
		meta := strings.TrimSpace(line[:idx])
		narrative := strings.TrimSpace(line[idx+2:])

		title := narrative
		note := ""
		if idx := strings.Index(narrative, ";"); idx != -1 {
			title = strings.TrimSpace(narrative[:idx])
			note = strings.TrimSpace(narrative[idx+1:])
		}

		task.Title = title
		task.Note = note

		for token := range strings.FieldsSeq(meta) {
			if projectRe.MatchString(token) && task.Project == "" {
				task.Project = strings.TrimPrefix(token, "+")
			} else if priorityRe.MatchString(token) && task.Priority == 0 {
				task.Priority = parsePriority(token)
			} else if durationRe.MatchString(token) && task.Duration == 0 {
				task.Duration, _ = parseDuration(token)
			}
		}
	}

	for _, match := range tagRe.FindAllStringSubmatch(line, -1) {
		task.Tags = append(task.Tags, match[1])
	}

	for _, match := range mentionRe.FindAllStringSubmatch(line, -1) {
		task.Mentions = append(task.Mentions, match[1])
	}

	for _, match := range urlRe.FindAllString(line, -1) {
		task.URLs = append(task.URLs, match)
	}

	return task, nil
}

func parseDuration(dur string) (time.Duration, error) {
	dur = strings.TrimPrefix(dur, "=")
	if len(dur) < 2 {
		return 0, nil
	}

	numStr := dur[:len(dur)-1]
	unit := dur[len(dur)-1]

	n, err := parseInt(numStr)
	if err != nil {
		return 0, nil
	}

	switch unit {
	case 'm':
		return time.Duration(n) * time.Minute, nil
	case 'h':
		return time.Duration(n) * time.Hour, nil
	case 'd':
		return time.Duration(n) * 24 * time.Hour, nil
	case 'w':
		return time.Duration(n) * 7 * 24 * time.Hour, nil
	default:
		return 0, nil
	}
}

func parsePriority(priority string) int {
	switch priority {
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

func formatDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	switch {
	case h >= 168 && h%168 == 0:
		return fmt.Sprintf("=%dw", h/168)
	case h >= 24 && h%24 == 0:
		return fmt.Sprintf("=%dd", h/24)
	case h > 0 && m == 0:
		return fmt.Sprintf("=%dh", h)
	case h > 0:
		return fmt.Sprintf("=%dh%dm", h, m)
	default:
		return fmt.Sprintf("=%dm", m)
	}
}

func parseInt(s string) (int, error) {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, nil
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}

func ParseDuration(dur string) (time.Duration, error) {
	return parseDuration(dur)
}

func FormatDuration(d time.Duration) string {
	return formatDuration(d)
}

func ParseAddInput(input string) model.Task {
	task := model.Task{Status: "-"}

	metaInput := input
	titleInput := ""

	if idx := strings.Index(input, "::"); idx != -1 {
		metaInput = strings.TrimSpace(input[:idx])
		titleInput = strings.TrimSpace(input[idx+2:])
	}

	for _, match := range tagRe.FindAllStringSubmatch(input, -1) {
		task.Tags = append(task.Tags, match[1])
	}
	for _, match := range mentionRe.FindAllStringSubmatch(input, -1) {
		task.Mentions = append(task.Mentions, match[1])
	}
	for _, match := range urlRe.FindAllString(input, -1) {
		task.URLs = append(task.URLs, match)
	}

	for _, token := range strings.Fields(metaInput) {
		if projectRe.MatchString(token) && task.Project == "" {
			task.Project = strings.TrimPrefix(token, "+")
		} else if priorityRe.MatchString(token) && task.Priority == 0 {
			task.Priority = parsePriority(token)
		} else if durationRe.MatchString(token) && task.Duration == 0 {
			task.Duration, _ = parseDuration(token)
		}
	}

	if titleInput != "" {
		task.Title = titleInput
	} else {
		rest := input
		rest = tagRe.ReplaceAllString(rest, "")
		rest = mentionRe.ReplaceAllString(rest, "")
		rest = projectRe.ReplaceAllString(rest, "")
		rest = priorityRe.ReplaceAllString(rest, "")
		rest = durationRe.ReplaceAllString(rest, "")
		rest = strings.Join(strings.Fields(rest), " ")
		task.Title = rest
	}

	return task
}

func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return string(data), nil
}
