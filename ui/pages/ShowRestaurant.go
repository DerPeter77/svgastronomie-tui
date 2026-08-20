package pages

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/DerPeter77/svgastronomie-tui/ui/events"
	"github.com/DerPeter77/svgastronomie-tui/ui/styles"

	sv "github.com/EchterTimo/go-svgastronomie"
)

// Types and Functions for Saving Restaurants in yaml File

// Bubbletea TUI

type ShowRestaurantModel struct {
	width, height     int
	showRestaurant    events.Restaurant
	scrapedRestaurant sv.Restaurant
	activeDayTab      int
	err               error
}

func NewShowRestaurantPage() ShowRestaurantModel {
	return ShowRestaurantModel{
		showRestaurant:    events.Restaurant{},
		scrapedRestaurant: sv.Restaurant{},
		activeDayTab:      0,
		err:               nil,
	}
}

func (m ShowRestaurantModel) Init() tea.Cmd {
	return tea.Batch(GetRestaurantsCmd("SavedRestaurants.yaml"), tea.RequestWindowSize)
}

func (m ShowRestaurantModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "left":
			length := len(m.scrapedRestaurant.Week.Days)
			m.activeDayTab = (m.activeDayTab - 1 + length) % length
		case "right":
			length := len(m.scrapedRestaurant.Week.Days)
			m.activeDayTab = (m.activeDayTab + 1) % length
		case "esc":
			return m, events.NavigateTo(events.ChooseRestaurantPageKey)
		case "q":
			return m, tea.Quit
		}
	case events.ChooseRestaurantMsg:
		m.showRestaurant = events.Restaurant(msg)

		scrapedRestaurant, err := sv.ScrapeRestaurant(m.showRestaurant.Url, nil)
		if err != nil {
			m.err = err
			return m, nil
		}

		m.scrapedRestaurant = *scrapedRestaurant
	}
	return m, nil
}

func (m ShowRestaurantModel) View() tea.View {
	headline := styles.Headline.Render("SV Restaurant TUI")

	text := styles.Text.Render(fmt.Sprintf("Speiseplan vom Restaurant: %v", m.showRestaurant.Name))
	text += "\n\n"

	// Days Tabs
	var days []string
	for index, day := range m.scrapedRestaurant.Week.Days {
		var daystring string
		if index == m.activeDayTab {
			daystring = styles.ActiveTab.Render(fmt.Sprintf("%v", day.Time.Format("_2.01 - Mon")))
		} else {
			daystring = styles.Tab.Render(fmt.Sprintf("%v", day.Time.Format("_2.01 - Mon")))
		}
		days = append(days, daystring)
	}

	tabsRow := lipgloss.JoinHorizontal(lipgloss.Top, days...)

	text += tabsRow
	text += "\n\n"

	// Dishes per tab
	var dishes []string
	for _, dish := range m.scrapedRestaurant.Week.Days[m.activeDayTab].Dishes {
		tags := ""
		if len(dish.Tags) > 0 {
			tags = " - "
			tags += strings.Join(dish.Tags, ", ")
		}
		dishes = append(dishes, styles.Border.Render(styles.Text.Render(fmt.Sprintf("%v - %v€ \n%v%v", dish.Name, dish.Price, dish.Description, tags))))
	}

	text += strings.Join(dishes, "\n")

	// Fehleranzeige
	if m.err != nil {
		text += "\n" + styles.ErrorText.Render(m.err.Error())
	}

	finalString := lipgloss.JoinVertical(lipgloss.Left, headline, text)
	view := tea.NewView(finalString)
	view.AltScreen = true
	return view
}
