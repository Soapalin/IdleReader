package main

import (
	"fmt"
	"game/engine/theme"
	"log"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *DashboardModel) BookshelfView() string {
	s := theme.Heading1.Render("Books") + "\n"
	start, end := m.bookPaginator.GetSliceBounds(len(m.ps.Reader.Library.Books))
	for i, b := range m.ps.Reader.Library.Books[start:end] {
		if b.Repeat > 0 {
			s += RenderEmojiOrFallback("✅", []string{"v"}) + " "
		} else {
			if b.Progress > 0 {
				s += RenderEmojiOrFallback("📖", []string{"-"}) + " "
			} else {
				s += RenderEmojiOrFallback("📓", []string{"x"}) + " "
			}
		}
		s += b.Name + ", " + b.Author
		if m.ps.Reader.CurrentReads.ContainsBook(b.ID) {
			s += " " + m.spinner.View()
		}
		if m.bs_cursor == (start + i) {
			s += theme.CursorArrow
		}
		s += "\n"
	}
	s += "\n"
	m.bookPaginator.SetTotalPages(len(m.ps.Reader.Library.Books))
	paginatorView := lipgloss.NewStyle().Width(m.width).AlignHorizontal(lipgloss.Position(0.5)).Render(m.bookPaginator.View())
	// pageNumber := fmt.Sprintf("<%d/%d>", m.bookPaginator.Page+1, m.bookPaginator.TotalPages)
	pageNumber := "<" + strconv.Itoa(m.bookPaginator.Page+1) + "/" + strconv.Itoa(m.bookPaginator.TotalPages) + ">"
	paginatorFull := lipgloss.JoinVertical(0, paginatorView, pageNumber)
	s += lipgloss.NewStyle().Width(m.width).AlignHorizontal(lipgloss.Position(0.5)).Render(paginatorFull)

	s += "\n" + theme.HelpIcon.Render("r") + theme.HelpText.Render(" read book • ")
	s += theme.HelpIcon.Render("enter") + theme.HelpText.Render(" book details • ")
	s += theme.HelpIcon.Render("esc") + theme.HelpText.Render(" back • ")
	s += theme.HelpIcon.Render("q") + theme.HelpText.Render(" quit")
	return s
}

func (m *DashboardModel) NextItemBookshelf() {
	if m.bs_cursor >= len(m.ps.Reader.Library.Books)-1 {
		return
	}
	m.bs_cursor++
	_, end := m.bookPaginator.GetSliceBounds(len(m.ps.Reader.Library.Books))

	if m.bs_cursor >= end {
		m.bookPaginator.Page++
	}

}

func (m *DashboardModel) PreviousItemBookshelf() {
	m.bs_cursor--
	if m.bs_cursor < 0 {
		m.bs_cursor = 0
		m.bookPaginator.Page = 0
		return
	}
	start, _ := m.bookPaginator.GetSliceBounds(len(m.ps.Reader.Library.Books))
	if m.bs_cursor < start {
		m.bookPaginator.Page--
	}

}

func (m *DashboardModel) PreviousBookPage() {
	m.bookPaginator.Page--
	_, m.bs_cursor = m.bookPaginator.GetSliceBounds(len(m.ps.Reader.Library.Books))
	m.bs_cursor--
	if m.bookPaginator.Page < 0 {
		m.bookPaginator.Page = 0
		m.bs_cursor = 0
	}
}

func (m *DashboardModel) NextBookPage() {
	if m.bookPaginator.Page >= m.bookPaginator.TotalPages-1 {
		m.bookPaginator.Page = m.bookPaginator.TotalPages - 1
		return
	}
	m.bookPaginator.Page++
	m.bs_cursor, _ = m.bookPaginator.GetSliceBounds(len(m.ps.Reader.Library.Books))

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
