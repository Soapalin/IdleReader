package main

import (
	"fmt"
	fileutils "game/engine/utils"
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

	//os check
	SYSTEM = runtime.GOOS

	if SYSTEM == "windows" {
		if len(os.Getenv("WT_SESSION")) > 0 {
			TERMINAL = "terminal"
		} else if len(os.Getenv("TERM_PROGRAM")) > 0 {
			TERMINAL = os.Getenv("TERM_PROGRAM")
		} else {
			TERMINAL = "consolehost"
		}
	}

	if len(os.Getenv("DEBUG")) > 0 {
		dir := createDocumentFolder()

		debugFile := filepath.Join(dir, "debug.log")
		if fileutils.IsFileExists(debugFile) {
			fileutils.KeepFromEnd(debugFile, 1000000)
		}
		f, err := tea.LogToFile(debugFile, "[DEBUG]")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
		DB.CreateAllBooksTable()
		DB.CreateAllItemsTable()
	}

	p := tea.NewProgram(InitialRootModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}

func PlatformPath(path string) string {
	if SYSTEM == "windows" {
		return filepath.Join(path, "Documents", "IdleReader")
	} else if SYSTEM == "linux" {
		return filepath.Join(path, ".IdleReader")
	} else if SYSTEM == "darwin" {
		// macOS
		return filepath.Join(path, "Library", "Application Support", "IdleReader")
	} else {
		return filepath.Join(path, "Documents", "IdleReader")
	}
}

func createDocumentFolder() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dir = PlatformPath(dir)

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
