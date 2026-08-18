package pages

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/DerPeter77/svgastronomie-tui/ui/events"
	"github.com/DerPeter77/svgastronomie-tui/ui/styles"
)

// Types and Functions for Saving Restaurants in yaml File

// Bubbletea TUI

type ShowRestaurantModel struct {
	showRestaurant events.Restaurant
	err            error
}

func NewShowRestaurantPage() ShowRestaurantModel {
	return ShowRestaurantModel{
		showRestaurant: events.Restaurant{},
		err:            nil,
	}
}

func (m ShowRestaurantModel) Init() tea.Cmd {
	return GetRestaurantsCmd("SavedRestaurants.yaml")
}

func (m ShowRestaurantModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
		case "down":
		case "esc":
			return m, events.NavigateTo(events.ChooseRestaurantPageKey)
		case "q":
			return m, tea.Quit
		}
	case events.ChooseRestaurantMsg:
		m.showRestaurant = events.Restaurant(msg)
	}
	return m, nil
}

func (m ShowRestaurantModel) View() tea.View {
	headline := styles.Headline.Render("SV Restaurant TUI")

	text := styles.Text.Render("Ansicht vom Restaurant: " + m.showRestaurant.Name)

	if m.err != nil {
		text += styles.ErrorText.Render(m.err.Error())
	}

	finalString := lipgloss.JoinVertical(lipgloss.Left, headline, text)
	return tea.NewView(finalString)
}
