package main

import (
	"fmt"
	"game/engine/theme"
	"log"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type newReaderKeyMap struct {
	Enter    key.Binding
	Esc      key.Binding
	Quit     key.Binding
	Previous key.Binding
	Next     key.Binding
	Help     key.Binding
}

var newreaderKeys = newReaderKeyMap{
	Enter:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit")),
	Esc:      key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Help:     key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	Quit:     key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
	Previous: key.NewBinding(key.WithKeys("up"), key.WithHelp("↑", "previous")),
	Next:     key.NewBinding(key.WithKeys("down"), key.WithHelp("↓", "next")),
}

func (k newReaderKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Next, k.Previous, k.Enter, k.Esc, k.Quit,
	}
}

func (k newReaderKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next, k.Previous, k.Enter},
		{k.Help, k.Esc, k.Quit},
	}
}

const (
	name = iota
	book
	author
)

type NewReaderModel struct {
	inputs  []textinput.Model
	focused int
	help    help.Model
	keys    newReaderKeyMap
}

func nameValidator(s string) error {
	return nil
}

func (m *NewReaderModel) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

func (m *NewReaderModel) previousInput() {
	m.focused--
	if m.focused < 0 {
		m.focused = 0
	}
}

func InitialNewReaderModel() NewReaderModel {
	var inputs []textinput.Model = make([]textinput.Model, 3)

	inputs[name] = textinput.New()
	inputs[name].Focus()
	// inputs[name].Placeholder = ""
	inputs[name].CharLimit = 255
	inputs[name].Width = 30
	inputs[name].Prompt = ""
	inputs[name].Validate = nameValidator

	inputs[book] = textinput.New()
	inputs[book].CharLimit = 600
	inputs[book].Placeholder = "The Poppy War"
	inputs[book].Width = 50

	inputs[author] = textinput.New()
	inputs[author].CharLimit = 255
	inputs[author].Placeholder = "R.F. Kuang"
	inputs[author].Width = 30

	return NewReaderModel{
		// Our to-do list is a grocery list
		inputs:  inputs,
		focused: 0,
		help:    help.New(),
		keys:    newreaderKeys,
	}
}

func NewPlayerSave(name string, book string, author string) PlayerSave {
	existingBook, err := DB.FindBookByNameAuthor(author, book)
	var playerLibrary Library
	if err != nil {
		existingBook = Book{
			ID:                      uuid.New(),
			Name:                    book,
			Author:                  author,
			Progress:                0,
			KnowledgeIncrease:       100,
			KnowledgeRequirement:    1,
			IntelligenceIncrease:    1,
			IntelligenceRequirement: 1,
			Pages:                   100,
			Repeat:                  0,
		}
		DB.InsertOrUpdateBook(existingBook)
	}
	playerLibrary.AddBookToLibrary(existingBook)
	newReader := Reader{
		ID:               uuid.New(),
		Name:             name,
		FavouriteBook:    book,
		FavouriteAuthor:  author,
		IQ:               40 + rand.Intn(29), // 40-69
		Fun:              100,
		Knowledge:        0,
		Prestige:         0,
		CurrentReadLimit: 1,
		CurrentReads:     CurrentReads{BookIDs: []uuid.UUID{existingBook.ID}},
		ReadingSpeed:     1,
		Library:          playerLibrary,
		Inventory:        GameItemDatabase{make([]Item, 0)},
	}
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dir = filepath.Join(dir, "Documents", "IdleReader", newReader.Name+"_"+newReader.ID.String()+"_Reader_Save.bin")
	log.Println("NewPlayerSave | " + dir)
	return PlayerSave{
		Reader:   newReader,
		Filename: dir,
		Shop:     InitShop(),
	}
}

func (m NewReaderModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m NewReaderModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.Type {

		// These keys should exit the program.
		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyEsc:
			switched := InitialSaveMenuModel()
			return InitialRootModel().SwitchScreen(&switched)

		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				playersave := NewPlayerSave(m.inputs[name].Value(), m.inputs[book].Value(), m.inputs[author].Value())
				playersave.SavePlayerToFile()
				switched := InitialDashboardModel(&playersave, 1, 0, 0)
				return InitialRootModel().SwitchScreen(&switched)
			}

		case tea.KeyUp, tea.KeyShiftTab:
			m.previousInput()

		case tea.KeyDown, tea.KeyTab:
			m.nextInput()

		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)

}

func (m NewReaderModel) View() string {
	// The header

	return fmt.Sprintf(
		` %s

 %s
 %s

 %s
 %s  %s
 %s  %s

`,
		theme.Heading1.Render("New Reader Game"),
		theme.InputStyle.Width(30).Render("Username"),
		m.inputs[name].View(),
		theme.Heading2.Render("Favourite Book"),
		theme.InputStyle.Width(50).Render("Book Name"),
		theme.InputStyle.Width(30).Render("Author"),
		m.inputs[book].View(),
		m.inputs[author].View(),
	) + m.help.View(m.keys) + "\n"
}

func NewReader() {
	p := tea.NewProgram(InitialNewReaderModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
