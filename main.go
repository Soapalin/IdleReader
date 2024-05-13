package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/text/encoding/charmap"
)

var SYSTEM string
var TERMINAL string

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

	//os check
	if runtime.GOOS == "windows" {
		SYSTEM = "windows"
	}

	if SYSTEM == "windows" {
		if len(os.Getenv("WT_SESSION")) > 0 {
			TERMINAL = "terminal"
		} else if len(os.Getenv("TERM_PROGRAM")) > 0 {
			TERMINAL = os.Getenv("TERM_PROGRAM")
		} else {
			TERMINAL = "consolehost"
		}
		// encoding check
		// in, err := exec.Command("chcp", "65001").Output()
		// if err != nil {
		// 	panic(err)
		// }
		// log.Println("chcp | " + string(in))
		// in, err = exec.Command("chcp").Output()
		// if err != nil {
		// 	panic(err)
		// }
		// log.Println("chcp | " + string(in))
	}

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

func RenderEmojiOrFallback(emoji string, fallback []string) string {
	emoji = strings.Trim(emoji, " ")

	if SYSTEM == "windows" && TERMINAL == "consolehost" {
		e := charmap.CodePage850.NewEncoder()
		out, err := e.String(emoji)
		for err != nil {
			// log.Println("RenderEmojiOrFallback | " + err.Error())
			for _, f := range fallback {
				out, err = e.String(f)
				if err == nil {
					return f
				}
			}
		}
		return out
	} else {
		return emoji
	}

}
