package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "[DEBUG]")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	// CreateAllBookLibBin()
	// CreateAllGameItemBin()

	currentUser, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	log.Println(currentUser)

	p := tea.NewProgram(InitialRootModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
