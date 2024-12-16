package tui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	primaryColor   = lipgloss.Color("#7D56F4")
	secondaryColor = lipgloss.Color("#F2F2F2")
	errorColor     = lipgloss.Color("#FF0000")

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor)

	DoneStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#4CAF50"))

	TableStyle = table.Styles{

		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor),

		Cell: lipgloss.NewStyle().
			Foreground(secondaryColor),

		Selected: lipgloss.NewStyle().
			Bold(true),
	}

	LogErrorStyle = lipgloss.NewStyle().Foreground(errorColor)
	LogInfoStyle  = lipgloss.NewStyle().Foreground(secondaryColor)

	// Additional styling for frames or borders if needed
	BorderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2).
			Margin(1, 2)
)
