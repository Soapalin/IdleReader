package main

import (
	"game/engine/theme"
	"log"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *DashboardModel) BookshelfView() string {
	s := theme.Heading1.Render("Books") + "\n"
	for i, b := range m.ps.Library.Books {
		s += b.Name + ", " + b.Author
		if m.bs_cursor == i {
			s += theme.CursorArrow
		}
		s += "\n"
	}
	s += theme.Heading1.Render("Inventory") + "\n"
	for index, i := range m.ps.Inventory.Items {
		s += i.Name
		if m.bs_cursor == (len(m.ps.Library.Books) + index) {
			s += theme.CursorArrow
		}
		s += "\n"
	}
	return s
}

func (m *DashboardModel) NextItemBookshelf() {
	m.bs_cursor++
	if m.bs_cursor >= len(m.ps.Inventory.Items)+len(m.ps.Library.Books) {
		m.bs_cursor = 0
	}
}

func (m *DashboardModel) PreviousItemBookshelf() {
	m.bs_cursor--
	if m.bs_cursor < 0 {
		m.bs_cursor = 0
	}
}

func (m *DashboardModel) TrySwitchBook() {
	if m.bs_cursor >= len(m.ps.Library.Books) {
		return
	}
	cr := m.ps.Reader.CurrentReads
	book_owned := m.ps.Library.Books
	var book_available Library

	for _, b := range book_owned {
		if !cr.ContainsBook(b) {
			book_available.AddBookToLibrary(b)
		}
	}
	m.errorMessage = "\n" + cr.String("DIGITS")

	m.errorMessage += "What book would you like to replace? (1-3)"
	m.bookChange = true

}

func (m *DashboardModel) ChooseBook(msg string) (tea.Model, tea.Cmd) {
	log.Println("ChooseBook")
	log.Println(m.bookChange)
	if m.bookChange {
		index, err := strconv.Atoi(msg)
		index = index - 1

		if index >= m.ps.Reader.CurrentReadLimit || err != nil {
			m.errorMessage = "Please choose between available options displayed"
			m.bookChange = false
			return m, nil
		}
		cr := &m.ps.Reader.CurrentReads
		if index >= len(cr.Books) {
			cr.AddBookToLibrary(m.ps.Library.Books[m.bs_cursor])
		} else {
			if cr.ContainsBook(m.ps.Library.Books[m.bs_cursor]) {
				m.bookChange = false
				m.errorMessage = "You are already reading this book!"
				return m, nil
			}
			cr.ReplaceBook(index, m.ps.Library.Books[m.bs_cursor])
		}
		m.bookChange = false
		m.errorMessage = "Current Book changed successfully!"

		var cmd []tea.Cmd
		cmd = append(cmd, m.progress[index].DecrPercent(1))
		cmd = append(cmd, m.progress[index].DecrPercent(cr.Books[index].Progress))

		return m, tea.Batch(cmd...)
	}
	return m, nil
}
