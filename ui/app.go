package ui

import (
	"svgastronomie-tui/ui/events"
	"svgastronomie-tui/ui/pages"

	tea "charm.land/bubbletea/v2"
)

type PageModel tea.Model

type AppModel struct {
	activePageKey events.PageKey
	pages         map[events.PageKey]PageModel
	sharedData    map[string]string // Zentraler Speicher für generische Daten
}

func NewAppModel() AppModel {
	return AppModel{
		activePageKey: events.InputPageKey,
		pages: map[events.PageKey]PageModel{
			events.InputPageKey:   pages.NewInputPage(),
			events.DisplayPageKey: pages.NewDisplayPage(),
		},
		sharedData: make(map[string]string),
	}
}

func (m AppModel) Init() tea.Cmd {
	if page, ok := m.pages[m.activePageKey]; ok {
		return page.Init()
	}
	return nil
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case events.NavigationMsg:
		if targetPage, ok := m.pages[msg.To]; ok {
			m.activePageKey = msg.To
			return m, targetPage.Init()
		}

	// Jede GlobalDataMsg im zentralen Speicher sichern UND an alle Seiten verteilen
	case events.GlobalDataMsg:
		m.sharedData[msg.Key] = msg.Value

		for key, page := range m.pages {
			updatedPage, cmd := page.Update(msg)
			if p, ok := updatedPage.(PageModel); ok {
				m.pages[key] = p
			}
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	}

	// Update für die aktuell aktive Seite
	current := m.pages[m.activePageKey]
	updated, cmd := current.Update(msg)

	if p, ok := updated.(PageModel); ok {
		m.pages[m.activePageKey] = p
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m AppModel) View() tea.View {
	if page, ok := m.pages[m.activePageKey]; ok {
		return page.View()
	}
	return tea.NewView("Unbekannte Seite")
}
