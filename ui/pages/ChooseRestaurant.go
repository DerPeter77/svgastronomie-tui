package pages

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/DerPeter77/svgastronomie-tui/ui/events"
	"github.com/DerPeter77/svgastronomie-tui/ui/styles"
	"github.com/goccy/go-yaml"
)

const savedDefaultPath string = "SavedRestaurants.yaml"

// Functions for Saving Restaurants in yaml File
func getSavedRestaurantsFromYAML(path string) (events.SavedRestaurants, error) {
	savedRestaurants := events.SavedRestaurants{}
	file, err := os.ReadFile(path)
	if err != nil {
		os.Create(savedDefaultPath)
		return savedRestaurants, err
	}
	yaml.Unmarshal(file, &savedRestaurants)
	return savedRestaurants, nil
}

func saveRestaurants(path string, savedRestaurants events.SavedRestaurants) error {
	bytes, err := yaml.Marshal(savedRestaurants)
	if err != nil {
		return err
	}
	os.WriteFile(path, bytes, 0644)
	return nil
}

func addRestaurant(path string, restaurant events.Restaurant) error {
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

type AddRestaurantCmdError error

func AddRestaurantCmd(path string, restaurant events.Restaurant) tea.Msg {
	err := addRestaurant(path, restaurant)
	return AddRestaurantCmdError(err)
}

func deleteRestaurant(path string, restaurant events.Restaurant) (deleted []events.Restaurant, err error) {
	saved, err := getSavedRestaurantsFromYAML(path)
	if err != nil {
		return nil, err
	}

	for index, value := range saved.Restaurants {
		if restaurant == value {
			deleted = append(deleted, value)
			saved.Restaurants = slices.Delete(saved.Restaurants, index, index+1)
			break
		}
	}

	err = saveRestaurants(path, saved)

	return deleted, err
}

type DeleteRestaurantMsg struct {
	deleted_restaurant []events.Restaurant
	err                error
}

func DeleteRestaurantCmd(path string, restaurant events.Restaurant) tea.Msg {
	deleted, err := deleteRestaurant(path, restaurant)
	return DeleteRestaurantMsg{deleted_restaurant: deleted, err: err}
}

type RestaurantsMsg struct {
	Restaurants events.SavedRestaurants
	err         error
}

func GetRestaurantsCmd(path string) tea.Cmd {
	restaurants, err := getSavedRestaurantsFromYAML(path)
	return func() tea.Msg { return RestaurantsMsg{Restaurants: restaurants, err: err} }
}

type ResetErrorTickMsg struct{}

func ErrorTick() tea.Msg {
	time.Sleep(time.Second * 2)
	return ResetErrorTickMsg{}
}

// Bubbletea TUI

type ChooseRestaurantModel struct {
	savedRestaurants events.SavedRestaurants
	cursorRestaurant int
	err              error

	addingMode bool
	nameInput  textinput.Model
	urlInput   textinput.Model
}

func NewChooseRestaurantPage() ChooseRestaurantModel {
	nameInput := textinput.New()
	nameInput.Placeholder = "Name"
	nameInput.SetVirtualCursor(true)
	nameInput.Blur()
	nameInput.CharLimit = 30
	nameInput.SetWidth(50)

	urlInput := textinput.New()
	urlInput.Placeholder = "Url"
	urlInput.SetVirtualCursor(true)
	urlInput.Blur()
	urlInput.CharLimit = 0
	urlInput.SetWidth(50)
	return ChooseRestaurantModel{
		savedRestaurants: events.SavedRestaurants{Restaurants: []events.Restaurant{}},
		cursorRestaurant: 0,
		err:              nil,
		addingMode:       false,
		nameInput:        nameInput,
		urlInput:         urlInput,
	}
}

func (m ChooseRestaurantModel) Init() tea.Cmd {
	return tea.Batch(ErrorTick, GetRestaurantsCmd(savedDefaultPath))
}

func (m ChooseRestaurantModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case ResetErrorTickMsg:
		m.err = nil
		return m, ErrorTick
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if len(m.savedRestaurants.Restaurants) == 0 {
				m.err = errors.New("Füge erstmal Restaurants hinzu!")
				return m, nil
			}
			length := len(m.savedRestaurants.Restaurants)
			m.cursorRestaurant = (m.cursorRestaurant - 1 + length) % length
		case "down":
			if len(m.savedRestaurants.Restaurants) == 0 {
				m.err = errors.New("Füge erstmal Restaurants hinzu!")
				return m, nil
			}
			length := len(m.savedRestaurants.Restaurants)
			m.cursorRestaurant = (m.cursorRestaurant + 1) % length
		case "enter":
			if m.addingMode {
				name := m.nameInput.Value()
				url := m.urlInput.Value()
				m.addingMode = false
				m.nameInput.Reset()
				m.urlInput.Reset()
				return m, tea.Batch(func() tea.Msg {
					return AddRestaurantCmd(savedDefaultPath, events.Restaurant{Name: name, Url: url})
				}, GetRestaurantsCmd(savedDefaultPath))
			} else {
				return m, tea.Batch(
					events.NavigateTo(events.ShowRestaurantPageKey),
					func() tea.Msg { return events.ChooseRestaurantMsg(m.savedRestaurants.Restaurants[m.cursorRestaurant]) },
				)
			}
		case "a":
			if !m.addingMode {
				m.addingMode = true
				m.nameInput.Focus()
				m.urlInput.Blur()
				return m, nil
			}
		case "tab":
			if m.addingMode {
				if m.nameInput.Focused() {
					m.nameInput.Blur()
					m.urlInput.Focus()
				} else {
					m.nameInput.Focus()
					m.urlInput.Blur()
				}
			}
		case "d":
			return m, func() tea.Msg {
				return DeleteRestaurantCmd(savedDefaultPath, m.savedRestaurants.Restaurants[m.cursorRestaurant])
			}
		case "r":
			return m, GetRestaurantsCmd(savedDefaultPath)
		case "q":
			return m, tea.Quit
		}
	case RestaurantsMsg:
		m.savedRestaurants = msg.Restaurants
		m.err = msg.err
	case AddRestaurantCmdError:
		m.err = msg
	case DeleteRestaurantMsg:
		if m.err != nil {
			m.err = fmt.Errorf("Error: %v", msg.err.Error())
		} else {
			m.err = fmt.Errorf("Deleted Restaurants: %v", msg.deleted_restaurant[0].Name)
		}
		return m, GetRestaurantsCmd(savedDefaultPath)
	}

	m.nameInput, cmd = m.nameInput.Update(msg)
	cmds = append(cmds, cmd)

	m.urlInput, cmd = m.urlInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, nil
}

func (m ChooseRestaurantModel) View() tea.View {
	headline := styles.Headline.Render("SV Restaurant TUI")

	text := styles.Text.Render("Wähle das Restaurant aus!")

	// Selection Mode
	if m.addingMode == false {
		var restaurants []string
		for index, value := range m.savedRestaurants.Restaurants {
			textstyle := styles.Text
			if index == m.cursorRestaurant {
				textstyle = styles.SelectedText
			}

			restaurants = append(restaurants, textstyle.Render(fmt.Sprintf("Restaurant: %v - Url: %v", value.Name, value.Url)))
		}

		text += styles.Text.Render("\n\n" + strings.Join(restaurants, "\n"))

		if m.err != nil {
			text += styles.ErrorText.Render("\n\n" + m.err.Error())
		}

		text += styles.Text.Render("\n\n\n\n\n[R] Neu laden   -   [A] Neues Restaurant   -   [D] Ausgewähltes Restaurant löschen")
	} else {
		// Adding Mode
		text += styles.Text.Render("\n" + m.nameInput.View())
		text += styles.Text.Render("\n" + m.urlInput.View())
	}

	finalString := lipgloss.JoinVertical(lipgloss.Left, headline, text)
	view := tea.NewView(finalString)
	view.AltScreen = true
	return view
}
