package main

import "game/engine/theme"

func (m *DashboardModel) InventoryView() string {
	s := ""
	s += theme.Heading1.Render("Inventory") + "\n"
	for index, i := range m.ps.Reader.Inventory.Items {
		s += "• " + i.Name
		if m.i_cursor == index {
			s += theme.CursorArrow
		}
		s += "\n"
	}

	return s
}

func (m *DashboardModel) NextItemInventory() {
	m.i_cursor++
	if m.i_cursor >= len(m.ps.Reader.Inventory.Items) {
		m.i_cursor = 0
	}
}

func (m *DashboardModel) PreviousItemInventory() {
	m.i_cursor--
	if m.i_cursor < 0 {
		m.i_cursor = 0
	}
}
