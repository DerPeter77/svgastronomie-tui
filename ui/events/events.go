package events

import tea "charm.land/bubbletea/v2"

type PageKey string

const (
	ChooseRestaurantPageKey PageKey = "chooseRestaurantPage"
	ShowRestaurantPageKey   PageKey = "showRestaurantPage"
)

type NavigationMsg struct {
	To PageKey
}

type ErrorMsg struct {
	Err error
}

func NavigateTo(key PageKey) tea.Cmd {
	return func() tea.Msg {
		return NavigationMsg{To: key}
	}
}
