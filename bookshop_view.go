package main

import (
	"game/engine/theme"
	"log"

	"github.com/charmbracelet/lipgloss"
)

func (m *DashboardModel) BookshopView() string {
	s := ""

	styleFunc := func(row, col int) lipgloss.Style {
		if row == 0 {
			return lipgloss.NewStyle().Margin(0)
		}
		if col == 0 || col == 1 {
			if row == m.ps.Shop.TableIndex {
				if row == m.ps.Shop.TableLen {
					return lipgloss.NewStyle().Width(25).Foreground(lipgloss.Color("5"))
				}
				return lipgloss.NewStyle().Width(25).MarginBottom(1).Foreground(lipgloss.Color("5"))

			}
			if row == m.ps.Shop.TableLen {
				return lipgloss.NewStyle().Width(25)
			}
			return lipgloss.NewStyle().Width(25).MarginBottom(1)
		}
		if col == 3 {
			if row == m.ps.Shop.TableIndex {
				return lipgloss.NewStyle().Width(18).MarginBottom(1).Foreground(lipgloss.Color("5"))
			}
			return lipgloss.NewStyle().Width(18).MarginBottom(1)
		}
		if row == m.ps.Shop.TableIndex {
			return lipgloss.NewStyle().Width(12).MarginBottom(1).Foreground(lipgloss.Color("5"))
		}
		return lipgloss.NewStyle().Width(12).MarginBottom(1)
	}
	m.ps.Shop.table.StyleFunc(styleFunc)
	s += m.ps.Shop.table.String()

	s += "\nLast Modified: " + m.ps.Shop.Modified.Format("02-01-2006 15:04:05")
	s += "\n* - owned"

	s += "\n\n" + theme.HelpIcon.Render("enter") + theme.HelpText.Render(" buy • ")
	s += theme.HelpIcon.Render("tab/shift+tab") + theme.HelpText.Render(" switch tabs • ")
	s += theme.HelpIcon.Render("i") + theme.HelpText.Render(" book/item info • ")
	s += theme.HelpIcon.Render("esc / q") + theme.HelpText.Render(" quit")

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
