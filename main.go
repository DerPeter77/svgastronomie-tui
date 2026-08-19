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

	// Example usage of the svgastronomie package
	// restaurantURL := "https://sv-gastronomie.de/menu/Gastronomie%20Weidm%C3%BCller,%20Detmold/Mittagsmen%C3%BC%20CTC"
	// fmt.Println("Scraping restaurant from URL:", restaurantURL)
	// restaurant, err := sv.ScrapeRestaurant(restaurantURL, nil)
	// if err != nil {
	// 	println("Error scraping restaurant:", err.Error())
	// 	return
	// }
	// fmt.Println("DONE")

	// // write to json file
	// var jsonData bytes.Buffer
	// encoder := json.NewEncoder(&jsonData)
	// encoder.SetEscapeHTML(false)
	// encoder.SetIndent("", "  ")
	// if err := encoder.Encode(restaurant); err != nil {
	// 	println("Error marshalling restaurant to JSON:", err.Error())
	// 	return
	// }
	// fileName := "restaurant.json"
	// err = os.WriteFile(fileName, jsonData.Bytes(), 0644)
	// if err != nil {
	// 	println("Error writing JSON to file:", err.Error())
	// 	return
	// }
	// fmt.Println("JSON data written to", fileName)
}
