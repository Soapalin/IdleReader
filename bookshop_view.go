package main

import (
	"log"

	"github.com/charmbracelet/lipgloss"
)

func (m *DashboardModel) BookshopView() string {
	s := ""

	styleFunc := func(row, col int) lipgloss.Style {
		if row == 0 {
			return lipgloss.NewStyle().Padding(0)
		}
		if col == 0 || col == 1 {
			if row == m.ps.Shop.TableIndex {
				return lipgloss.NewStyle().Width(28).Padding(1).Foreground(lipgloss.Color("5"))
			}
			return lipgloss.NewStyle().Width(28).Padding(1)
		}
		if row == m.ps.Shop.TableIndex {
			return lipgloss.NewStyle().Width(12).Padding(1).Foreground(lipgloss.Color("5"))
		}
		return lipgloss.NewStyle().Width(12).Padding(1, 1)
	}
	m.ps.Shop.table.StyleFunc(styleFunc)
	s += m.ps.Shop.table.String()

	s += "\n\nLast Modified: " + m.ps.Shop.Modified.Format("02-01-2006 15:04:05")

	return s
}

func (m *DashboardModel) TryBuy() {
	itemID := m.ps.Shop.GetShopItemByIndex()
	if m.ps.AlreadyOwned(itemID) {
		log.Println("AlreadyOwned returned True")
		m.errorMessage = "You already owned this item/book!"
	} else {
		typeItem := m.ps.Shop.Buy(&m.ps)
		switch typeItem {
		case unknownTransaction:
			m.errorMessage = "Unable to purchase item/book."
		case knowledgeMissingTransaction:
			m.errorMessage = "Not enough knowledge."
		case iqMissingTransaction:
			m.errorMessage = "Not enough IQ."
		case bookTransaction, itemTransaction:
			m.errorMessage = "Item/Book bought successfully! Check your Bookshelf."
		}
	}
}
