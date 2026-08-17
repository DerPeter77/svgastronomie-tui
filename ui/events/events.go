package events

import tea "charm.land/bubbletea/v2"

type PageKey string

const (
	InputPageKey   PageKey = "input_page"
	DisplayPageKey PageKey = "display_page"
)

type NavigationMsg struct {
	To PageKey
}

type ErrorMsg struct {
	Err error
}

// Generelle Nachricht zur Datenübertragung zwischen Seiten
type GlobalDataMsg struct {
	Key   string
	Value string
}

func NavigateTo(key PageKey) tea.Cmd {
	return func() tea.Msg {
		return NavigationMsg{To: key}
	}
}