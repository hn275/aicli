package ui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func RenderInput() string {
	p := tea.NewProgram(InitialModel())
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	return InputData
}
