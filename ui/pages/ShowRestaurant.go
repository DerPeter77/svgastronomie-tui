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
	showRestaurant Restaurant
	err            error
}

func NewShowRestaurantPage() ShowRestaurantModel {
	return ShowRestaurantModel{
		showRestaurant: Restaurant{},
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
			events.NavigateTo(events.ChooseRestaurantPageKey)
			return m, nil
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ShowRestaurantModel) View() tea.View {
	headline := styles.Headline.Render("SV Restaurant TUI")

	text := styles.Text.Render("Wähle das Restaurant aus!\n")

	if m.err != nil {
		text += styles.ErrorText.Render(m.err.Error())
	}

	finalString := lipgloss.JoinVertical(lipgloss.Left, headline, text)
	return tea.NewView(finalString)
}
