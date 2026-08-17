package pages

import (
	"svgastronomie-tui/ui/events"
	"fmt"

	tea "charm.land/bubbletea/v2"
)

type InputPageModel struct {
	Name string
}

func NewInputPage() InputPageModel {
	return InputPageModel{Name: ""}
}

func (m InputPageModel) Init() tea.Cmd {
	return nil
}

func (m InputPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if len(m.Name) > 0 {
				submittedName := m.Name
				m.Name = ""

				// Generelle GlobalDataMsg senden UND navigieren
				return m, tea.Batch(
					func() tea.Msg {
						return events.GlobalDataMsg{
							Key:   "username",
							Value: submittedName,
						}
					},
					events.NavigateTo(events.DisplayPageKey),
				)
			}
		case "backspace":
			if len(m.Name) > 0 {
				m.Name = m.Name[:len(m.Name)-1]
			}
		default:
			if len(msg.String()) == 1 {
				m.Name += msg.String()
			}
		}
	}
	return m, nil
}

func (m InputPageModel) View() tea.View {
	content := fmt.Sprintf(
		"--- INPUT SEITE ---\n\nName eingeben: %s_\n\n[Drücke 'Enter' zum Bestätigen]",
		m.Name,
	)
	return tea.NewView(content)
}
