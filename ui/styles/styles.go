package styles

import "charm.land/lipgloss/v2"

var Headline lipgloss.Style = lipgloss.NewStyle().
	Foreground(lipgloss.Cyan).
	Bold(true).
	MarginBottom(1)

var Text lipgloss.Style = lipgloss.NewStyle().
	Foreground(lipgloss.White)

var SelectedText lipgloss.Style = lipgloss.NewStyle().
	Foreground(lipgloss.Blue).
	Bold(true)

var ErrorText lipgloss.Style = lipgloss.NewStyle().
	Foreground(lipgloss.BrightRed).
	Bold(true)

var Border lipgloss.Style = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.BrightBlack).
	Padding(0, 2)

var activeTabBorder = lipgloss.Border{
	Top:         "─",
	Bottom:      " ",
	Left:        "│",
	Right:       "│",
	TopLeft:     "╭",
	TopRight:    "╮",
	BottomLeft:  "┘",
	BottomRight: "└",
}

var tabBorder = lipgloss.Border{
	Top:         "─",
	Bottom:      "─",
	Left:        "│",
	Right:       "│",
	TopLeft:     "╭",
	TopRight:    "╮",
	BottomLeft:  "┴",
	BottomRight: "┴",
}

var Tab = lipgloss.NewStyle().
	Border(tabBorder, true).
	BorderForeground(lipgloss.White).
	Padding(0, 1)

var ActiveTab = Tab.
	Border(activeTabBorder, true).
	Foreground(lipgloss.BrightWhite).
	Bold(true).
	BorderForeground(lipgloss.White)
