package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	theme "game/engine/theme"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/google/uuid"

	// "github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
)

type keyMap struct {
	Enter key.Binding
	Quit  key.Binding
	Esc   key.Binding
	Help  key.Binding
	Up    key.Binding
	Down  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},   // first column
		{k.Help, k.Quit}, // second column
	}
}

var keys = keyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select choice"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
	),
}

type SaveMenuModel struct {
	choices       []string         // items on the to-do list
	cursor        int              // which to-do list item our cursor is pointing at
	selected      map[int]struct{} // which to-do items are selected
	continue_save string
	load_save     bool
	load_input    textinput.Model
	model_err     string
	help          help.Model
	keys          keyMap
}

type PlayerSave struct {
	Reader   Reader
	Filename string
	Shop     Shop
}

func (ps *PlayerSave) AlreadyOwned(id uuid.UUID) bool {
	for _, b := range ps.Reader.Library.Books {
		if b.ID == id {
			return true
		}
	}
	for _, i := range ps.Reader.Inventory.Items {
		if i.ID == id {
			return true
		}
	}
	return false
}

func GetLatestSave() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Println("GetLatestSave | os.UserHomeDir")
		panic(err)
	}
	dir = PlatformPath(dir)

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Println("GetLatestSave | os.ReadDir")
		log.Println(dir)
		// panic(err)
	}

	log.Println(files)
	var reader_save string
	var last_saved time.Time
	for _, sv := range files {
		if strings.HasSuffix(sv.Name(), "Reader_Save.bin") {
			if reader_save == "" {
				reader_save = sv.Name()
				last_saved_info, err := os.Stat(filepath.Join(dir, sv.Name()))
				if err != nil {
					panic(err)
				}
				last_saved = last_saved_info.ModTime()
				log.Println(reader_save)
			} else {
				new_save_info, err := os.Stat(filepath.Join(dir, sv.Name()))
				if err != nil {
					panic(err)
				}
				new_save_time := new_save_info.ModTime()
				if new_save_time.After(last_saved) {
					reader_save = sv.Name()
				}
			}
		}
	}

	if reader_save == "" {
		return ""
	}

	path := filepath.Join(dir, reader_save)
	return path
}

func LoadPlayerFromFile(filename string) (PlayerSave, error) {
	f, err := os.Open(filename)
	if err != nil {
		log.Println("LoadPlayerFromFile | f.Open")
		log.Println(filename)
		panic(err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}

	// data := make([]byte, fi.Size())
	// n, err := f.Read(data)

	// if err != nil && err != io.EOF {
	// 	log.Println("LoadPlayerFromFile | f.Read")
	// 	log.Println(filename)
	// 	panic(err)
	// }
	// if n == 0 {
	// 	log.Println("n = 0")
	// }
	// var playersave PlayerSave

	// err = json.Unmarshal(data, &playersave)
	// if err != nil {
	// 	log.Println("Unmarshal")
	// 	panic(err)
	// }
	data := make([]byte, fi.Size())
	n, err := f.Read(data)
	if err != nil && err != io.EOF {
		log.Println("LoadPlayerFromFile | f.Read")
		log.Println(filename)
		panic(err)
	}
	if n == 0 {
		log.Println("n = 0")
	}

	buff := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buff)

	playersave := PlayerSave{
		Filename: filename,
	}

	dec_err := dec.Decode(&playersave)
	if dec_err != nil {
		log.Println("LoadPlayerFromFile | dec.Decode")
		log.Println(playersave)
		log.Println(dec_err)
		return playersave, errors.New("save file not valid! Try another file")
	}
	playersave.Shop.LoadShopTable(&playersave.Reader.Library)
	// log.Println("LoadPlayerFromFile | Decode Successful")
	log.Println(playersave)

	return playersave, nil
}

func (r *PlayerSave) SavePlayerToFile() {
	f, err := os.Create(r.Filename)
	if err != nil {
		log.Println("SavePlayerToFile | os.Create")
		log.Println(r.Filename)
		panic(err)
	}
	defer f.Close()

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)

	error_enc := enc.Encode(*r)
	if error_enc != nil {
		panic(error_enc)
	}

	if _, err := f.Write(buff.Bytes()); err != nil {
		panic(err)
	}

	// g, err := os.Create(r.Filename)
	// if err != nil {
	// 	log.Println("SavePlayerToFile | os.Create")
	// 	log.Println(r.Filename)
	// 	panic(err)
	// }
	// defer g.Close()

	// b, err := json.Marshal(*r)
	// log.Println("SavePlayerString")
	// log.Println(string(b))
	// // var playersave PlayerSave
	// // err = json.Unmarshal(b, &playersave)
	// // log.Println(playersave)

	// if err != nil {
	// 	panic(err)
	// }

	// if _, err := g.Write(b); err != nil {
	// 	panic(err)
	// }
}

func InitialSaveMenuModel() SaveMenuModel {
	save_file := GetLatestSave()
	var choices []string
	if save_file == "" {
		choices = []string{"New Game", "Load Game"}
	} else {
		// readerName := strings.Replace(save_file, "saves/", "", -1)
		readerName := filepath.Base(save_file)
		readerName = strings.Split(readerName, "_")[0]
		choices = []string{"Continue - " + readerName, "New Game", "Load Game"}
	}

	load_input := textinput.New()
	load_input.CharLimit = 500
	load_input.Width = 100
	load_input.Placeholder = "C:/saves/example.bin"
	// load_input.Prompt = "Enter the file path of your save"
	return SaveMenuModel{
		choices:       choices,
		continue_save: save_file,
		load_save:     false,
		load_input:    load_input,
		model_err:     "",
		help:          help.New(),
		keys:          keys,
		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m SaveMenuModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m SaveMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			if msg.String() == "q" {
				if !m.load_save {
					return m, tea.Quit
				}
			} else {
				return m, tea.Quit
			}

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case tea.KeyEsc.String():
			m.load_save = false
			m.load_input.SetValue("")
		case "?":
			m.help.ShowAll = !m.help.ShowAll
		case "enter":
			switch strings.TrimSpace(strings.Split(m.choices[m.cursor], "-")[0]) {
			case "Continue":
				log.Println("Continue")
				if m.continue_save != "" {
					playersave, err := LoadPlayerFromFile(m.continue_save)
					if err != nil {
						m.model_err = err.Error()
					} else {
						switched := InitialDashboardModel(&playersave, 1, 0, 0)
						return InitialRootModel().SwitchScreen(&switched)
					}

				}
			case "New Game":
				log.Println("New Game")
				switched := InitialNewReaderModel()
				return InitialRootModel().SwitchScreen(&switched)
			case "Load Game":
				log.Println("Load Game")
				if m.load_save {
					if m.load_input.Value() != "" {
						// check if file exists
						if _, err := os.Stat(m.load_input.Value()); err == nil {
							playersave, err := LoadPlayerFromFile(strings.TrimSpace(m.load_input.Value()))
							if err != nil {
								m.model_err = err.Error()
							} else {
								switched := InitialDashboardModel(&playersave, 1, 0, 0)
								return InitialRootModel().SwitchScreen(&switched)
							}
						} else if errors.Is(err, os.ErrNotExist) {
							// file does not exists
							m.model_err = "Save File does not exist!"
						} else {
							panic(err)
						}
					}
				} else {
					m.load_save = true
					m.load_input.Focus()
				}

			default:
				log.Println(m.choices[m.cursor])
				log.Println(strings.Split(m.choices[m.cursor], "-")[0])
			}
		}
	}
	if m.load_save {
		var cmd tea.Cmd
		m.load_input, cmd = m.load_input.Update(msg)
		return m, cmd
	}
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m SaveMenuModel) View() string {
	// The header
	s := theme.Heading1.Render("Idle Book Reader")
	s += "\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = "x" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("[%s] %s\n", cursor, choice)
	}
	if m.load_save {
		m.load_input.Focus()
		s += "\n"
		s += theme.Heading2.Render("Enter the file path of your save")
		s += "\n"
		s += m.load_input.View()
	}
	if m.model_err != "" {
		s += "\n"
		s += m.model_err
	}

	// The footer
	s += "\n"
	// s += m.help.View(m.keys)
	s += "\n" + theme.HelpIcon.Render("?") + theme.HelpText.Render(" help • ")
	s += theme.HelpIcon.Render("q") + theme.HelpText.Render(" quit")

	// Send the UI for rendering
	return s
}
