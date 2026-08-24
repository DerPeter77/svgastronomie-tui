package pages

import (
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/DerPeter77/svgastronomie-tui/ui/events"
	"github.com/DerPeter77/svgastronomie-tui/ui/styles"

	sv "github.com/EchterTimo/go-svgastronomie"
)

var sandspinner = []string{"⠁", "⠂", "⠄", "⡀", "⡈", "⡐", "⡠", "⣀", "⣁", "⣂", "⣄", "⣌", "⣔",
	"⣤", "⣥", "⣦", "⣮", "⣶", "⣷", "⣿", "⡿", "⠿", "⢟", "⠟", "⡛", "⠛", "⠫", "⢋", "⠋", "⠍", "⡉", "⠉", "⠑", "⠡", "⢁"}

// var circlespinner = []string{"◜", "◠", "◝", "◞", "◡", "◟"}

// var starspinner = []string{"✶", "✸", "✹", "✺", "✹", "✷"}

type scrapedRestaurantMsg struct {
	restaurant sv.Restaurant
	err        error
}

func scrapeRestaurant(url string) tea.Msg {
	scrapedRestaurant, err := sv.ScrapeRestaurant(url, nil)

	return scrapedRestaurantMsg{
		restaurant: *scrapedRestaurant,
		err:        err,
	}
}

type ShowRestaurantModel struct {
	width, height     int
	showRestaurant    events.Restaurant
	scrapedRestaurant sv.Restaurant
	activeDayTab      int
	err               error

	spinner spinner.Model
}

func NewShowRestaurantPage() ShowRestaurantModel {
	spinner := spinner.New()
	spinner.Spinner.Frames = sandspinner
	spinner.Spinner.FPS = time.Second / 12
	return ShowRestaurantModel{
		showRestaurant:    events.Restaurant{},
		scrapedRestaurant: sv.Restaurant{},
		activeDayTab:      0,
		err:               nil,
		spinner:           spinner,
	}
}

func (m ShowRestaurantModel) Init() tea.Cmd {
	return tea.Batch(tea.RequestWindowSize, m.spinner.Tick)
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
			m.activeDayTab = (m.activeDayTab - 1 + length) % length
		case "right", "d":
			length := len(m.scrapedRestaurant.Week.Days)
			m.activeDayTab = (m.activeDayTab + 1) % length
		case "esc":
			return m, events.NavigateTo(events.ChooseRestaurantPageKey)
		case "q":
			return m, tea.Quit
		}
	case events.ChooseRestaurantMsg:
		m.showRestaurant = events.Restaurant(msg)

		return m, func() tea.Msg { return scrapeRestaurant(m.showRestaurant.Url) }
	case scrapedRestaurantMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.scrapedRestaurant = sv.Restaurant(msg.restaurant)
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)

	return m, cmd
}

func (m ShowRestaurantModel) View() tea.View {
	headline := styles.Headline.Render("SV Restaurant TUI")

	text := styles.Text.Render(fmt.Sprintf("Speiseplan vom Restaurant: %v", m.showRestaurant.Name))
	text += "\n\n"

	if len(m.scrapedRestaurant.Week.Days) > 0 {
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
			dishes = append(dishes, styles.Border.Render(styles.Text.Render(fmt.Sprintf("%v - %.2f€ \n%v%v", dish.Name, dish.Price, dish.Description, tags))))
		}

		text += strings.Join(dishes, "\n")
	} else {
		text += "\n " + m.spinner.View()
	}

	// Fehleranzeige
	if m.err != nil {
		text += "\n" + styles.ErrorText.Render(m.err.Error())
	}

	finalString := lipgloss.JoinVertical(lipgloss.Left, headline, text)
	view := tea.NewView(finalString)
	view.AltScreen = true
	return view
}
