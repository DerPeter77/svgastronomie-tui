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
	loading           bool
	err               error
}

type ScrapeRestaurantMsg struct {
	restaurant *sv.Restaurant
	err        error
}

func ScrapeRestaurantCmd(url string) tea.Cmd {
	return func() tea.Msg {
		scrapedRestaurant, err := sv.ScrapeRestaurant(url, nil)
		if err != nil {
			return ScrapeRestaurantMsg{restaurant: nil, err: err}
		}

		return ScrapeRestaurantMsg{restaurant: scrapedRestaurant, err: nil}
	}
}

func NewShowRestaurantPage() ShowRestaurantModel {
	return ShowRestaurantModel{
		showRestaurant:    events.Restaurant{},
		scrapedRestaurant: sv.Restaurant{},
		activeDayTab:      0,
		loading:           false,
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
		case "left", "a":
			length := len(m.scrapedRestaurant.Week.Days)
			if length == 0 {
				return m, nil
			}
			m.activeDayTab = (m.activeDayTab - 1 + length) % length
		case "right", "d":
			length := len(m.scrapedRestaurant.Week.Days)
			if length == 0 {
				return m, nil
			}
			m.activeDayTab = (m.activeDayTab + 1) % length
		case "esc":
			return m, events.NavigateTo(events.ChooseRestaurantPageKey)
		case "q":
			return m, tea.Quit
		}
	case events.ChooseRestaurantMsg:
		m.showRestaurant = events.Restaurant(msg)
		m.loading = true
		m.err = nil
		m.scrapedRestaurant = sv.Restaurant{}
		m.activeDayTab = 0

		return m, ScrapeRestaurantCmd(m.showRestaurant.Url)

	case ScrapeRestaurantMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		if msg.restaurant != nil {
			m.scrapedRestaurant = *msg.restaurant
		}

		daysLen := len(m.scrapedRestaurant.Week.Days)
		if daysLen == 0 {
			m.activeDayTab = 0
		} else if m.activeDayTab >= daysLen {
			m.activeDayTab = daysLen - 1
		}
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

	if m.loading {
		text += styles.Text.Render("Lade Speiseplan...")

		finalString := lipgloss.JoinVertical(lipgloss.Left, headline, text)
		view := tea.NewView(finalString)
		view.AltScreen = true
		return view
	}

	if len(m.scrapedRestaurant.Week.Days) == 0 {
		text += styles.Text.Render("Kein Speiseplan verfuegbar.")

		if m.err != nil {
			text += "\n" + styles.ErrorText.Render(m.err.Error())
		}

		finalString := lipgloss.JoinVertical(lipgloss.Left, headline, text)
		view := tea.NewView(finalString)
		view.AltScreen = true
		return view
	}

	// Dishes per tab
	var dishes []string
	for _, dish := range m.scrapedRestaurant.Week.Days[m.activeDayTab].Dishes {
		tags := ""
		if len(dish.Tags) > 0 {
			tags = " - "
			tags += strings.Join(dish.Tags, ", ")
		}
		dishes = append(dishes, styles.Border.Render(styles.Text.Render(fmt.Sprintf("%v - %.2f€ \n%v%v", dish.Name, dish.Price, dish.Description, tags))))
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
