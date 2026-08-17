package pages

import (
	"svgastronomie-tui/ui/events"
	"fmt"

	tea "charm.land/bubbletea/v2"
)

type DisplayPageModel struct {
	Name string
}

func NewDisplayPage() DisplayPageModel {
	return DisplayPageModel{Name: "Noch kein Name gesetzt"}
}

func (m DisplayPageModel) Init() tea.Cmd {
	return nil
}

func (m DisplayPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case events.GlobalDataMsg:
		// Generischen Key auswerten
		if msg.Key == "username" {
			m.Name = msg.Value
		}

	case tea.KeyMsg:
		if msg.String() == "esc" {
			return m, events.NavigateTo(events.InputPageKey)
		}
	}
	return m, nil
}

func (m DisplayPageModel) View() tea.View {
	content := fmt.Sprintf(
		"--- DISPLAY SEITE ---\n\nEmpfangener Name: %s\n\n[Drücke 'ESC' um zurückzugehen]",
		m.Name,
	)
	return tea.NewView(content)
}
