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
	if m.cr_cursor >= len(m.ps.Reader.CurrentReads.BookIDs) {
		m.cr_cursor = 0
	}
}

func (m *DashboardModel) CurrentReadsView() string {
	s := ""
	r := &m.ps.Reader
	for i, cr := range r.CurrentReads.BookIDs {
		cr_data, err := m.ps.Reader.Library.GetBookPointerByID(cr)
		if err != nil {
			panic(err)
		}
		s += theme.Heading2.Render(cr_data.Name + ", " + cr_data.Author)
		if i == m.cr_cursor {
			s += " ←"
		}
		s += "\n"
		s += "\n"
		// s += strconv.Itoa(cr.Progress) + "/100"
		// s += "\n\n"
		m.progress[i].Width = m.width - 3
		s += "[" + m.progress[i].View() + "]"
		s += "\n\n"
		s += "\n"

	}

	book_locked := MAX_SIMULT_BOOK - r.CurrentReadLimit
	book_unlocked := MAX_SIMULT_BOOK - book_locked - len(r.CurrentReads.BookIDs)
	n := 0
	for n < book_unlocked {
		s += theme.Heading2.Render("[Select a book to read in your Bookshelf]")
		s += "\n"
		s += "*/100\n\n"
		n++
	}
	n = 0

	for n < book_locked {
		s += theme.Heading2.Render("[Increase your IQ to read multiple books simultaneously]")
		s += "\n*/100\n\n"
		n++
	}

	return s
}
