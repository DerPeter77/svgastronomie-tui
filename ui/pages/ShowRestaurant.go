package pages

import (
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/DerPeter77/svgastronomie-tui/ui/events"
	"github.com/DerPeter77/svgastronomie-tui/ui/styles"

	sv "github.com/EchterTimo/go-svgastronomie"
)

func Version() (version string, ok bool) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return version, ok
	}

	// 1. If consumed as a module dependency in another application:
	for _, dep := range info.Deps {
		if dep.Path == "github.com/DerPeter77/svgastronomie-tui" {
			return dep.Version, true
		}
	}

	// 2. If running directly from inside this module (e.g. tests or main binary):
	if info.Main.Path == "github.com/DerPeter77/svgastronomie-tui" && info.Main.Version != "" {
		return info.Main.Version, true
	}

	return version, false
}

// Spinners
// var sandspinner = []string{"⠁", "⠂", "⠄", "⡀", "⡈", "⡐", "⡠", "⣀", "⣁", "⣂", "⣄", "⣌", "⣔",
// 	"⣤", "⣥", "⣦", "⣮", "⣶", "⣷", "⣿", "⡿", "⠿", "⢟", "⠟", "⡛", "⠛", "⠫", "⢋", "⠋", "⠍", "⡉", "⠉", "⠑", "⠡", "⢁"}

var circlespinner = []string{"◜", "◠", "◝", "◞", "◡", "◟"}

// var starspinner = []string{"✶", "✸", "✹", "✺", "✹", "✷"}

type scrapedRestaurantMsg struct {
	restaurant sv.Restaurant
	err        error
}

func scrapeRestaurant(url string) tea.Msg {
	scrapedRestaurant, err := sv.ScrapeRestaurant(url, nil)
	if err != nil {
		return scrapedRestaurantMsg{err: err}
	}

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

	isLoading bool
	spinner   spinner.Model
}

func NewShowRestaurantPage() ShowRestaurantModel {
	spinner := spinner.New()
	spinner.Spinner.Frames = circlespinner
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
			if !m.isLoading {
				length := len(m.scrapedRestaurant.Week.Days)
				m.activeDayTab = (m.activeDayTab - 1 + length) % length
			}
		case "right", "d":
			if !m.isLoading {
				length := len(m.scrapedRestaurant.Week.Days)
				m.activeDayTab = (m.activeDayTab + 1) % length
			}
		case "esc":
			newModel := NewShowRestaurantPage()
			return newModel, events.NavigateTo(events.ChooseRestaurantPageKey)
		case "q":
			return m, tea.Quit
		}
	case events.ChooseRestaurantMsg:
		m.showRestaurant = events.Restaurant(msg)
		m.isLoading = true
		return m, func() tea.Msg { return scrapeRestaurant(m.showRestaurant.Url) }
	case scrapedRestaurantMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.scrapedRestaurant = sv.Restaurant(msg.restaurant)
		m.isLoading = false
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)

	return m, cmd
}

func (m ShowRestaurantModel) View() tea.View {
	version, ok := Version()
	if !ok {
		version = "not found ..."
	}
	headline := styles.Headline.Render(fmt.Sprintf("SV Restaurant TUI   %v", version))

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
		if len(m.scrapedRestaurant.Week.Days[m.activeDayTab].Dishes) > 0 {
			for _, dish := range m.scrapedRestaurant.Week.Days[m.activeDayTab].Dishes {
				tags := ""
				if len(dish.Tags) > 0 {
					tags = " - "
					tags += strings.Join(dish.Tags, ", ")
				}
				dishes = append(dishes,
					styles.Border.Render(
						lipgloss.Wrap(
							styles.Text.Render(fmt.Sprintf("%v - %.2f€ \n%v%v", dish.Name, dish.Price, dish.Description, tags)),
							m.width-10, "",
						),
					),
				)
			}
		} else {
			dishes = append(dishes, styles.Border.Render("Für den Tag wurden keine Gerichte gefunden :("))
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
