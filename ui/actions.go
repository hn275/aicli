package ui

import (
	"log"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func RenderInput() *tea.Program {
	p := tea.NewProgram(InitialModel())
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	return p
}

func MockRequest(wg *sync.WaitGroup, m *model) {
	time.Sleep(time.Second)
	m.output = "result"
	wg.Done()
}
