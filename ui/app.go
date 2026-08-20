package ui

import (
	"fmt"
	"runtime/debug"

	"github.com/DerPeter77/svgastronomie-tui/ui/events"
	"github.com/DerPeter77/svgastronomie-tui/ui/pages"

	tea "charm.land/bubbletea/v2"
)

const modulePath string = "github.com/DerPeter77/svgastronomie-tui"

func Version() (version string, ok bool) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return version, ok
	}

	// 1. If consumed as a module dependency in another application:
	for _, dep := range info.Deps {
		if dep.Path == modulePath {
			return dep.Version, true
		}
	}

	// 2. If running directly from inside this module (e.g. tests or main binary):
	if info.Main.Path == modulePath && info.Main.Version != "" {
		return info.Main.Version, true
	}

	return version, false
}

type PageModel tea.Model

type AppModel struct {
	activePageKey events.PageKey
	pages         map[events.PageKey]PageModel
	sharedData    map[string]string // Zentraler Speicher für generische Daten
}

func NewAppModel() AppModel {
	return AppModel{
		activePageKey: events.ChooseRestaurantPageKey,
		pages: map[events.PageKey]PageModel{
			events.ChooseRestaurantPageKey: pages.NewChooseRestaurantPage(),
			events.ShowRestaurantPageKey:   pages.NewShowRestaurantPage(),
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

	case events.ChooseRestaurantMsg:
		model, cmd := m.pages[events.ShowRestaurantPageKey].Update(msg)
		m.pages[events.ShowRestaurantPageKey] = model
		return m, cmd
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
	version, ok := Version()
	if !ok {
		version = "not found ..."
	}
	if page, ok := m.pages[m.activePageKey]; ok {
		view := page.View()
		view.Content += fmt.Sprintf("\n\n\n\n\nVersion: %v", version)
		return view
	}
	return tea.NewView("Unbekannte Seite")
}
