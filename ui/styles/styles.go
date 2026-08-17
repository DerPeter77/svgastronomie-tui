package styles

import "charm.land/lipgloss/v2"

var Headline lipgloss.Style = lipgloss.NewStyle().
	Foreground(lipgloss.Cyan).
	Bold(true).
	MarginBottom(1)

var Text lipgloss.Style = lipgloss.NewStyle().
	Foreground(lipgloss.White)

var SelectedText lipgloss.Style = lipgloss.NewStyle().
	Foreground(lipgloss.BrightBlue).
	Bold(true)

var Border lipgloss.Style = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.BrightBlack).
	Padding(1, 2)
