package cli

import "github.com/charmbracelet/lipgloss"

var (
	statusTodo    = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	statusActive  = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	statusDone    = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	statusDropped = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))

	projectStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	priorityP1    = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	priorityP2    = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	priorityP3    = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	durationStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	noteStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true)
	lineNumStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)
