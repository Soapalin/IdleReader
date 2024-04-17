package main

import (
	"game/engine/theme"
)

func (m *DashboardModel) PreviousCurrentReads() {
	m.cr_cursor--
	if m.cr_cursor < 0 {
		m.cr_cursor = 0
	}
}

func (m *DashboardModel) NextCurrentReads() {
	m.cr_cursor++
	if m.cr_cursor >= len(m.ps.Reader.CurrentReads.Books) {
		m.cr_cursor = 0
	}
}

func (m *DashboardModel) CurrentReadsView() string {
	s := ""
	r := &m.ps.Reader
	for i, cr := range r.CurrentReads.Books {
		s += theme.Heading2.Render(cr.Name + ", " + cr.Author)
		if i == m.cr_cursor {
			s += " ←"
		}
		s += "\n"
		// s += strconv.Itoa(cr.Progress) + "/100"
		// s += "\n\n"
		s += "[" + m.progress[i].View() + "]"
		s += "\n\n"

	}

	book_locked := MAX_SIMULT_BOOK - r.CurrentReadLimit
	book_unlocked := MAX_SIMULT_BOOK - book_locked - len(r.CurrentReads.Books)
	n := 0
	for n < book_unlocked {
		s += theme.Heading2.Render("[Select a book to read in your Bookshelf]")
		s += "\n"
		s += "*/100\n\n"
		n++
	}
	n = 0

	for n < book_locked {
		s += theme.Heading2.Render("[Increase your inteligence to read multiple books simultaneously]")
		s += "\n*/100\n\n"
		n++
	}

	return s
}
