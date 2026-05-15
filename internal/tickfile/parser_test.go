package tickfile

import (
	"os"
	"testing"
	"time"
)

func TestParseLine_todo(t *testing.T) {
	task, _ := parseLine("- :: Buy oat milk", 1)
	if task.Status != "-" {
		t.Errorf("status = %q, want '-'", task.Status)
	}
	if task.Title != "Buy oat milk" {
		t.Errorf("title = %q, want 'Buy oat milk'", task.Title)
	}
	if task.Note != "" {
		t.Errorf("note = %q, want ''", task.Note)
	}
	if task.Project != "" {
		t.Errorf("project = %q, want ''", task.Project)
	}
}

func TestParseLine_active(t *testing.T) {
	task, _ := parseLine("/ =2h +groww #android :: Fix trace batching crash on cold start", 1)
	if task.Status != "/" {
		t.Errorf("status = %q, want '/'", task.Status)
	}
	if task.Title != "Fix trace batching crash on cold start" {
		t.Errorf("title = %q, want 'Fix trace batching crash on cold start'", task.Title)
	}
	if task.Duration != 2*time.Hour {
		t.Errorf("duration = %v, want 2h0m0s", task.Duration)
	}
	if task.Project != "groww" {
		t.Errorf("project = %q, want 'groww'", task.Project)
	}
	if len(task.Tags) != 1 || task.Tags[0] != "android" {
		t.Errorf("tags = %v, want ['android']", task.Tags)
	}
}

func TestParseLine_done(t *testing.T) {
	task, _ := parseLine("x =3h +groww #fno :: Hook P&L dashboard to ledger API", 1)
	if task.Status != "x" {
		t.Errorf("status = %q, want 'x'", task.Status)
	}
}

func TestParseLine_dropped(t *testing.T) {
	task, _ := parseLine("~ +blog :: Add comments system ; killed scope - not worth the spam fight", 1)
	if task.Status != "~" {
		t.Errorf("status = %q, want '~'", task.Status)
	}
	if task.Note != "killed scope - not worth the spam fight" {
		t.Errorf("note = %q, want 'killed scope - not worth the spam fight'", task.Note)
	}
}

func TestParseLine_priority(t *testing.T) {
	tests := []struct {
		line     string
		expected int
	}{
		{"- p1 +proj :: task", 1},
		{"- p2 +proj :: task", 2},
		{"- p3 +proj :: task", 3},
		{"- :: task", 0},
	}
	for _, tc := range tests {
		task, _ := parseLine(tc.line, 1)
		if task.Priority != tc.expected {
			t.Errorf("parseLine(%q).Priority = %q, want %q", tc.line, task.Priority, tc.expected)
		}
	}
}

func TestParseLine_duration(t *testing.T) {
	tests := []struct {
		line     string
		expected time.Duration
	}{
		{"- =15m :: task", 15 * time.Minute},
		{"- =2h :: task", 2 * time.Hour},
		{"- =1d :: task", 24 * time.Hour},
		{"- =1w :: task", 7 * 24 * time.Hour},
		{"- :: task", 0},
	}
	for _, tc := range tests {
		task, _ := parseLine(tc.line, 1)
		if task.Duration != tc.expected {
			t.Errorf("parseLine(%q).Duration = %v, want %v", tc.line, task.Duration, tc.expected)
		}
	}
}

func TestParseLine_project(t *testing.T) {
	task, _ := parseLine("- +groww :: work on stuff", 1)
	if task.Project != "groww" {
		t.Errorf("project = %q, want 'groww'", task.Project)
	}
}

func TestParseLine_noProjectIsValid(t *testing.T) {
	task, _ := parseLine("- :: untagged task", 1)
	if task.Project != "" {
		t.Errorf("project = %q, want ''", task.Project)
	}
}

func TestParseLine_tags(t *testing.T) {
	task, _ := parseLine("- +proj #backend #api :: task", 1)
	if len(task.Tags) != 2 || task.Tags[0] != "backend" || task.Tags[1] != "api" {
		t.Errorf("tags = %v, want ['backend', 'api']", task.Tags)
	}
}

func TestParseLine_tagsInlineInTitle(t *testing.T) {
	task, _ := parseLine("- :: update #backend config", 1)
	if len(task.Tags) != 1 || task.Tags[0] != "backend" {
		t.Errorf("tags = %v, want ['backend']", task.Tags)
	}
}

func TestParseLine_tagsInNote(t *testing.T) {
	task, _ := parseLine("- :: fix it ; see #bug123", 1)
	if len(task.Tags) != 1 || task.Tags[0] != "bug123" {
		t.Errorf("tags = %v, want ['bug123']", task.Tags)
	}
}

func TestParseLine_mentions(t *testing.T) {
	task, _ := parseLine("- :: review with @saheb", 1)
	if len(task.Mentions) != 1 || task.Mentions[0] != "saheb" {
		t.Errorf("mentions = %v, want ['saheb']", task.Mentions)
	}
}

func TestParseLine_urls(t *testing.T) {
	task, _ := parseLine("- :: review docs ; https://internal.dev/docs/deploy", 1)
	if len(task.URLs) != 1 || task.URLs[0] != "https://internal.dev/docs/deploy" {
		t.Errorf("urls = %v, want ['https://internal.dev/docs/deploy']", task.URLs)
	}
}

func TestParseLine_note(t *testing.T) {
	task, _ := parseLine("- :: task description ; contextual note here", 1)
	if task.Title != "task description" {
		t.Errorf("title = %q, want 'task description'", task.Title)
	}
	if task.Note != "contextual note here" {
		t.Errorf("note = %q, want 'contextual note here'", task.Note)
	}
}

func TestParseLine_missingDelimiter(t *testing.T) {
	task, err := parseLine("- just a line without delimiter", 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if task.Status != "-" {
		t.Errorf("status = %q, want '-'", task.Status)
	}
	if task.Title != "" {
		t.Errorf("title = %q, want ''", task.Title)
	}
}

func TestParseLine_noStatus(t *testing.T) {
	task, _ := parseLine(":: bare task", 1)
	if task.Status != "" {
		t.Errorf("status = %q, want ''", task.Status)
	}
}

func TestParseLine_firstProjectWins(t *testing.T) {
	task, _ := parseLine("- +first +second :: task", 1)
	if task.Project != "first" {
		t.Errorf("project = %q, want 'first'", task.Project)
	}
}

func TestParseLine_firstDurationWins(t *testing.T) {
	task, _ := parseLine("- =1d =2h :: task", 1)
	if task.Duration != 24*time.Hour {
		t.Errorf("duration = %v, want 24h0m0s", task.Duration)
	}
}

func TestParseLine_priorityOnlyBeforeDelimiter(t *testing.T) {
	task, _ := parseLine("- :: task p1 is in title and ignored", 1)
	if task.Priority != 0 {
		t.Errorf("priority = %q, want '' (p1 after ::)", task.Priority)
	}
}

func TestParseLine_durationOnlyBeforeDelimiter(t *testing.T) {
	task, _ := parseLine("- :: API = broken", 1)
	if task.Duration != 0 {
		t.Errorf("duration = %v, want 0 (= after ::)", task.Duration)
	}
}

func TestParseLine_lineNumber(t *testing.T) {
	task, _ := parseLine("- :: task", 42)
	if task.Line != 42 {
		t.Errorf("Line = %d, want 42", task.Line)
	}
}

func TestParseLine_multipleStatusCharactersInTitle(t *testing.T) {
	task, _ := parseLine("- :: foo - bar / baz x done", 1)
	if task.Status != "-" {
		t.Errorf("status = %q, want '-'", task.Status)
	}
}

func TestParseLine_p1InHeaderOnly(t *testing.T) {
	task, _ := parseLine("- p1 +proj :: task", 1)
	if task.Priority != 1 {
		t.Errorf("priority = %q, want 'p1'", task.Priority)
	}
}

func TestParseFile_validFile(t *testing.T) {
	f, _ := os.CreateTemp("", "tick_test_*.tick")
	defer os.Remove(f.Name())
	f.WriteString("- :: task one\n")
	f.WriteString("/ :: task two\n")
	f.WriteString("x :: task three\n")
	f.Close()

	tasks, err := ParseFile(f.Name())
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if len(tasks) != 3 {
		t.Fatalf("got %d tasks, want 3", len(tasks))
	}
	if tasks[0].Title != "task one" {
		t.Errorf("tasks[0].Title = %q, want 'task one'", tasks[0].Title)
	}
	if tasks[1].Status != "/" {
		t.Errorf("tasks[1].Status = %q, want '/'", tasks[1].Status)
	}
}

func TestParseFile_skipsComments(t *testing.T) {
	f, _ := os.CreateTemp("", "tick_test_*.tick")
	defer os.Remove(f.Name())
	f.WriteString(";tick:0.0.1\n")
	f.WriteString("; this is a comment\n")
	f.WriteString("- :: task one\n")
	f.WriteString("\n")
	f.WriteString("- :: task two\n")
	f.Close()

	tasks, err := ParseFile(f.Name())
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if len(tasks) != 2 {
		t.Fatalf("got %d tasks, want 2", len(tasks))
	}
}

func TestParseFile_emptyFile(t *testing.T) {
	f, _ := os.CreateTemp("", "tick_test_*.tick")
	defer os.Remove(f.Name())
	f.Close()

	tasks, err := ParseFile(f.Name())
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("got %d tasks, want 0", len(tasks))
	}
}

func TestParseFile_onlyComments(t *testing.T) {
	f, _ := os.CreateTemp("", "tick_test_*.tick")
	defer os.Remove(f.Name())
	f.WriteString(";tick:0.0.1\n")
	f.WriteString("; nothing to see here\n")
	f.Close()

	tasks, err := ParseFile(f.Name())
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("got %d tasks, want 0", len(tasks))
	}
}

func TestParseFile_notFoundIsEmpty(t *testing.T) {
	tasks, err := ParseFile("/tmp/nonexistent_tick_file_42.tick")
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("got %d tasks, want 0", len(tasks))
	}
}

func TestParseFile_allFields(t *testing.T) {
	line := "/ =2h p1 +groww #android #bug :: Fix crash on cold start ; review with @saheb"
	f, _ := os.CreateTemp("", "tick_test_*.tick")
	defer os.Remove(f.Name())
	f.WriteString(line + "\n")
	f.Close()

	tasks, err := ParseFile(f.Name())
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("got %d tasks, want 1", len(tasks))
	}

	task := tasks[0]
	if task.Status != "/" {
		t.Errorf("Status = %q, want '/'", task.Status)
	}
	if task.Duration != 2*time.Hour {
		t.Errorf("Duration = %v, want 2h0m0s", task.Duration)
	}
	if task.Priority != 1 {
		t.Errorf("Priority = %q, want 'p1'", task.Priority)
	}
	if task.Project != "groww" {
		t.Errorf("Project = %q, want 'groww'", task.Project)
	}
	if task.Title != "Fix crash on cold start" {
		t.Errorf("Title = %q, want 'Fix crash on cold start'", task.Title)
	}
	if task.Note != "review with @saheb" {
		t.Errorf("Note = %q, want 'review with @saheb'", task.Note)
	}
	if len(task.Tags) != 2 || task.Tags[0] != "android" || task.Tags[1] != "bug" {
		t.Errorf("Tags = %v, want ['android', 'bug']", task.Tags)
	}
	if len(task.Mentions) != 1 || task.Mentions[0] != "saheb" {
		t.Errorf("Mentions = %v, want ['saheb']", task.Mentions)
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
	}{
		{"=15m", 15 * time.Minute},
		{"=2h", 2 * time.Hour},
		{"=1d", 24 * time.Hour},
		{"=1w", 7 * 24 * time.Hour},
		{"15m", 15 * time.Minute},
		{"2h", 2 * time.Hour},
		{"=0m", 0},
		{"=abc", 0},
		{"", 0},
		{"=", 0},
	}
	for _, tc := range tests {
		d, _ := parseDuration(tc.input)
		if d != tc.expected {
			t.Errorf("parseDuration(%q) = %v, want %v", tc.input, d, tc.expected)
		}
	}
}
