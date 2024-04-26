package main

import (
	"fmt"
	"game/engine/theme"
	"log"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *DashboardModel) BookshelfView() string {
	s := theme.Heading1.Render("Books") + "\n"
	for i, b := range m.ps.Reader.Library.Books {
		if b.Repeat > 0 {
			s += "✅"
		} else {
			if b.Progress > 0 {
				s += "📖"
			} else {
				s += "📓"
			}
		}
		s += b.Name + ", " + b.Author
		if m.ps.Reader.CurrentReads.ContainsBook(b.ID) {
			s += " " + m.spinner.View()
		}
		if m.bs_cursor == i {
			s += theme.CursorArrow
		}
		s += "\n"
	}
	s += theme.Heading1.Render("Inventory") + "\n"
	for index, i := range m.ps.Reader.Inventory.Items {
		s += "• " + i.Name
		if m.bs_cursor == (len(m.ps.Reader.Library.Books) + index) {
			s += theme.CursorArrow
		}
		s += "\n"
	}

	s += "\n" + theme.HelpIcon.Render("r") + theme.HelpText.Render(" read book • ")
	s += theme.HelpIcon.Render("enter") + theme.HelpText.Render(" book details • ")
	s += theme.HelpIcon.Render("esc") + theme.HelpText.Render(" back • ")
	s += theme.HelpIcon.Render("q") + theme.HelpText.Render(" quit")
	return s
}

func (m *DashboardModel) NextItemBookshelf() {
	m.bs_cursor++
	if m.bs_cursor >= len(m.ps.Reader.Inventory.Items)+len(m.ps.Reader.Library.Books) {
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
	if m.bs_cursor >= len(m.ps.Reader.Library.Books) {
		return
	}
	cr := m.ps.Reader.CurrentReads
	book_owned := m.ps.Reader.Library.Books
	var book_available Library

	for _, b := range book_owned {
		if !cr.ContainsBook(b.ID) {
			book_available.AddBookToLibrary(b)
		}
	}
	m.errorMessage = "\n" + cr.String("DIGITS")
	canRead := m.ps.Reader.CurrentReadLimit - len(cr.BookIDs)
	digit := len(cr.BookIDs) + 1
	for canRead > 0 {
		m.errorMessage += fmt.Sprintf("%d. [select a book to read]\n", digit)
		digit++
		canRead--
	}

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
		if index >= len(cr.BookIDs) {
			cr.AddBookToLibrary(m.ps.Reader.Library.Books[m.bs_cursor].ID)
		} else {
			if cr.ContainsBook(m.ps.Reader.Library.Books[m.bs_cursor].ID) {
				m.bookChange = false
				m.errorMessage = "You are already reading this book!"
				return m, nil
			}
			cr.ReplaceBook(index, m.ps.Reader.Library.Books[m.bs_cursor].ID)
		}
		m.bookChange = false
		m.errorMessage = "Current Book changed successfully!"

		var cmd []tea.Cmd
		cmd = append(cmd, m.progress[index].DecrPercent(1))
		b_value, err := m.ps.Reader.Library.GetBookPointerByID(cr.BookIDs[index])
		if err != nil {
			panic(err)
		}
		cmd = append(cmd, m.progress[index].IncrPercent(b_value.Progress))

		return m, tea.Batch(cmd...)
	}
	return m, nil
}
