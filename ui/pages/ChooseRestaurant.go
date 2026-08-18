package pages

import (
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/DerPeter77/svgastronomie-tui/ui/styles"
	"github.com/goccy/go-yaml"
)

type Restaurant struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

type SavedRestaurants struct {
	Restaurants []Restaurant `yaml:"restaurants"`
}

func getSavedRestaurantsFromYAML(path string) (SavedRestaurants, error) {
	savedRestaurants := SavedRestaurants{}
	file, err := os.ReadFile(path)
	if err != nil {
		return savedRestaurants, err
	}
	yaml.Unmarshal(file, &savedRestaurants)
	return savedRestaurants, nil
}

func saveRestaurants(path string, savedRestaurants SavedRestaurants) error {
	bytes, err := yaml.Marshal(savedRestaurants)
	if err != nil {
		return err
	}
	os.WriteFile(path, bytes, 0644)
	return nil
}

func addRestaurant(path string, restaurant Restaurant) error {
	saved, err := getSavedRestaurantsFromYAML(path)
	if err != nil {
		return err
	}
	saved.Restaurants = append(saved.Restaurants, restaurant)
	err = saveRestaurants(path, saved)
	if err != nil {
		return err
	}
	return nil
}

func testAddRestaurant(path string, restaurant Restaurant) tea.Cmd {
	err := addRestaurant(path, restaurant)
	return func() tea.Msg { return err }
}

type ChooseRestaurantModel struct {
	name             string
	savedRestaurants SavedRestaurants
}

func NewChooseRestaurantPage() ChooseRestaurantModel {
	return ChooseRestaurantModel{name: ""}
}

func (m ChooseRestaurantModel) Init() tea.Cmd {
	return testAddRestaurant("SavedRestaurants.yaml", Restaurant{Name: "Test Init", Url: "Iergendeine URL"})
}

func (m ChooseRestaurantModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ChooseRestaurantModel) View() tea.View {
	headline := styles.Headline.Render("SV Restaurant TUI")

	text := styles.Text.Render("Gib die URL von deinem Restaurant ein!")

	finalString := lipgloss.JoinVertical(lipgloss.Left, headline, text)
	return tea.NewView(finalString)
}
