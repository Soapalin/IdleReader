package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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

	// encoding check
	in, err := exec.Command("chcp", "65001").Output()
	if err != nil {
		panic(err)
	}
	log.Println("chcp | " + string(in))

	// encoding check
	in, err = exec.Command("chcp").Output()
	if err != nil {
		panic(err)
	}
	log.Println("chcp | " + string(in))

	p := tea.NewProgram(InitialRootModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
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
