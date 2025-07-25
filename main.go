package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jbeausoleil/task-tracker-bubbletea-extended/internal/task"
)

func main() {
	p := tea.NewProgram(task.InitialModel())
	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
