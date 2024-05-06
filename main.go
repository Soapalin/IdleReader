package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)


func main() {
	os.Setenv("DEBUG", "true")
	if len(os.Getenv("DEBUG")) > 0 {
		dir := createDocumentFolder()

		debugFile := filepath.Join(dir, "debug.log")
		f, err := tea.LogToFile(debugFile, "[DEBUG]")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
		DB.CreateAllBooksTable()
		DB.CreateAllItemsTable()
	}

	// CreateAllBookLibBin()
	// CreateAllGameItemBin()

	// UpdateAllBooksLibrary()
	// UpdateAllGameItemDatabase()

	p := tea.NewProgram(InitialRootModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}

func createDocumentFolder() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dir = filepath.Join(dir, "Documents", "IdleReader")
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return dir
}
