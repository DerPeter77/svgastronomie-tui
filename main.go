package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/DerPeter77/svgastronomie-tui/ui"
)

func main() {
	p := tea.NewProgram(ui.NewAppModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Fehler beim Starten der App: %v\n", err)
		os.Exit(1)
	}
}
