package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hn275/aicli/ui"
)

func main() {
	p := tea.NewProgram(ui.NewModel())
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}
