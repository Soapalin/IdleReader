package main


func (m *DashboardModel) ExitView() string {
	s := "\n"
	for i, ch := range m.exitChoices {
		if i == m.ex_cursor {
			s += "[x] " + ch + "\n"

		} else {
			s += "[ ] "+ ch + "\n"
		}
	}
	return s
}


func (m *DashboardModel) NextExitChoice() {
	m.ex_cursor++
	if m.ex_cursor >= len(m.exitChoices) {
		m.ex_cursor = 0
	}
}

func (m *DashboardModel) PreviousExitChoice() {
	m.ex_cursor-- 
	if m.ex_cursor < 0 {
		m.ex_cursor = 0 
	}
} 