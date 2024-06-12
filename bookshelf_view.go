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
	s := theme.Heading1.Render("Books") + "\n\n"
	s += theme.Heading2.Render("Search") + "\n"
	s += m.bookshelf_input.View()
	s += "\n"

	start, end := m.bookPaginator.GetSliceBounds(len(m.bookshelfLibrary.Books))
	for i, b := range m.bookshelfLibrary.Books[start:end] {
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
	m.bookPaginator.SetTotalPages(len(m.bookshelfLibrary.Books))
	paginatorView := lipgloss.NewStyle().Width(m.width).AlignHorizontal(lipgloss.Position(0.5)).Render(m.bookPaginator.View())
	// pageNumber := fmt.Sprintf("<%d/%d>", m.bookPaginator.Page+1, m.bookPaginator.TotalPages)
	pageNumber := "<" + strconv.Itoa(m.bookPaginator.Page+1) + "/" + strconv.Itoa(m.bookPaginator.TotalPages) + ">"
	paginatorFull := lipgloss.JoinVertical(0, paginatorView, pageNumber)
	s += lipgloss.NewStyle().Width(m.width).AlignHorizontal(lipgloss.Position(0.5)).Render(paginatorFull)

	s += "\n" + theme.HelpIcon.Render("enter") + theme.HelpText.Render(" read book/submit search • ")
	s += theme.HelpIcon.Render("ctrl+F") + theme.HelpText.Render(" begin search • ")
	s += theme.HelpIcon.Render("ctrl+X") + theme.HelpText.Render(" clear search • ")
	s += theme.HelpIcon.Render("i") + theme.HelpText.Render(" book info • ")
	s += theme.HelpIcon.Render("tab/shift+tab") + theme.HelpText.Render(" switch tabs • ")
	s += theme.HelpIcon.Render("esc") + theme.HelpText.Render(" back • ")
	s += theme.HelpIcon.Render("q") + theme.HelpText.Render(" quit")
	return s
}

func (m *DashboardModel) NextItemBookshelf() {
	if m.bookshelf_input.Focused() {
		return
	}
	if m.bs_cursor >= len(m.bookshelfLibrary.Books)-1 {
		return
	}
	m.bs_cursor++
	_, end := m.bookPaginator.GetSliceBounds(len(m.bookshelfLibrary.Books))

	if m.bs_cursor >= end {
		m.bookPaginator.Page++
	}

}

func (m *DashboardModel) PreviousItemBookshelf() {
	if m.bookshelf_input.Focused() {
		return
	}
	m.bs_cursor--
	if m.bs_cursor < 0 {
		m.bs_cursor = 0
		m.bookPaginator.Page = 0
		return
	}
	start, _ := m.bookPaginator.GetSliceBounds(len(m.bookshelfLibrary.Books))
	if m.bs_cursor < start {
		m.bookPaginator.Page--
	}

}

func (m *DashboardModel) PreviousBookPage() {
	if m.bookshelf_input.Focused() {
		return
	}
	m.bookPaginator.Page--
	_, m.bs_cursor = m.bookPaginator.GetSliceBounds(len(m.bookshelfLibrary.Books))
	m.bs_cursor--
	if m.bookPaginator.Page < 0 {
		m.bookPaginator.Page = 0
		m.bs_cursor = 0
	}
}

func (m *DashboardModel) NextBookPage() {
	if m.bookshelf_input.Focused() {
		return
	}
	if m.bookPaginator.Page >= m.bookPaginator.TotalPages-1 {
		m.bookPaginator.Page = m.bookPaginator.TotalPages - 1
		return
	}
	m.bookPaginator.Page++
	m.bs_cursor, _ = m.bookPaginator.GetSliceBounds(len(m.bookshelfLibrary.Books))

}

func (m *DashboardModel) ResetPaginator() {
	m.bookPaginator.Page = 0
	m.bs_cursor = 0
	m.bookPaginator.TotalPages = len(m.bookshelfLibrary.Books)
}

func (m *DashboardModel) SubmitBookshelfSearch() {
	log.Println("SubmitBookshelfSearch | ")
	if m.bookshelf_input.Value() == "" {
		m.bookshelf_input.Blur()
		m.bookshelfLibrary = m.ps.Reader.Library
		m.ResetPaginator()
		return
	}
	temp_lib, err := DB.FindBooksContainsSingle(m.bookshelf_input.Value(), m.bookshelf_input.Value())
	if err != nil {
		panic(err)
	}
	log.Println(temp_lib)
	var lib Library
	for _, b := range m.ps.Reader.Library.Books {
		if temp_lib.ContainsBook(b) {
			lib.AddBookToLibrary(b)
		}
	}
	m.bookshelfLibrary = lib
	m.ResetPaginator()

	m.bookshelf_input.Blur()
}

func (m *DashboardModel) TrySwitchBook() {
	if m.bookshelf_input.Focused() {
		return
	}
	if m.bs_cursor >= len(m.bookshelfLibrary.Books) {
		return
	}
	cr := m.ps.Reader.CurrentReads
	book_owned := m.bookshelfLibrary.Books
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
			cr.AddBookToLibrary(m.bookshelfLibrary.Books[m.bs_cursor].ID)
		} else {
			if cr.ContainsBook(m.bookshelfLibrary.Books[m.bs_cursor].ID) {
				m.bookChange = false
				m.errorMessage = "You are already reading this book!"
				return m, nil
			}
			cr.ReplaceBook(index, m.bookshelfLibrary.Books[m.bs_cursor].ID)
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

func (m *DashboardModel) UnfocusBookshelfInput() {
	m.bookshelf_input.SetValue("")
	m.bookshelf_input.Blur()
}

func (m *DashboardModel) ClearBookshelfInput() {
	m.bookshelf_input.SetValue("")
	m.bookshelf_input.Blur()
	m.bookshelfLibrary = m.ps.Reader.Library
}
