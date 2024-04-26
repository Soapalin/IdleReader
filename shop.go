package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"

	_ "embed"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type Item struct {
	ID            uuid.UUID
	Name          string
	Description   string
	Cost          int
	IqRequirement int
	Bought        bool
	Effect        string
}

type GameItemDatabase struct {
	Items []Item
}

func (g *GameItemDatabase) ContainsItem(item Item) bool {
	for _, i := range g.Items {
		if i.ID == item.ID {
			return true
		}
	}
	return false
}

func (g *GameItemDatabase) ContainsItemByName(item Item) bool {
	for _, i := range g.Items {
		if i.Name == item.Name {
			return true
		}
	}
	return false
}

func (g *GameItemDatabase) AddItem(item Item) {
	g.Items = append(g.Items, item)
}

// var AllGameItems = GameItemDatabase{
// 	Items: []Item{
// 		{uuid.New(), "Reading Glasses", "Increases Reading speed by 20%", 10000, 1, false, nil},
// 		{uuid.New(), "Bookmark", "Know where you left your reading.", 100, 1, false, nil},
// 		{uuid.New(), "Reading Light", "Everything is more clear. Your gain an additional 10% knowledge.", 100, 1, false, nil},
// 	},
// }

var AllGameItems GameItemDatabase = LoadAllGameItems()

type Shop struct {
	Books      Library
	Items      GameItemDatabase
	ShopSize   int
	Modified   time.Time
	table      table.Table
	TableIndex int
	TableLen   int
}

type TransactionResult int

const (
	bookTransaction TransactionResult = iota
	itemTransaction
	knowledgeMissingTransaction
	iqMissingTransaction
	unknownTransaction
)

func (s *Shop) Buy(ps *PlayerSave) TransactionResult {
	log.Println("Buy | ")
	index := s.TableIndex - 1
	if index < len(s.Books.Books) {
		log.Println(s.Books.Books[index].Name)
		if s.Books.Books[index].IntelligenceRequirement > ps.Reader.IQ {
			return iqMissingTransaction
		}
		if s.Books.Books[index].KnowledgeRequirement > ps.Reader.Knowledge {
			return knowledgeMissingTransaction
		}
		ps.Reader.DecreaseKnowledge(s.Books.Books[index].KnowledgeRequirement)
		ps.Reader.Library.AddBookToLibrary(s.Books.Books[index])
		return bookTransaction
	} else if (index - len(s.Books.Books)) < len(s.Items.Items) {
		log.Println(s.Items.Items[index-len(s.Books.Books)].Name)
		if s.Items.Items[index-len(s.Books.Books)].IqRequirement > ps.Reader.IQ {
			return iqMissingTransaction
		}
		if s.Items.Items[index-len(s.Books.Books)].Cost > ps.Reader.Knowledge {
			return knowledgeMissingTransaction
		}
		ps.Reader.DecreaseKnowledge(s.Items.Items[index-len(s.Books.Books)].Cost)
		ps.Reader.Inventory.AddItem(s.Items.Items[index-len(s.Books.Books)])
		return itemTransaction
	}
	return unknownTransaction

}

func GetItemByID(id uuid.UUID) Item {
	for _, i := range AllGameItems.Items {
		if i.ID == id {
			return i
		}
	}
	return Item{}
}

func (s *Shop) GetShopItemByIndex() uuid.UUID {
	log.Println("GetShopItemByIndex | ")
	shopIndex := s.TableIndex - 1
	if shopIndex < len(s.Books.Books) {
		log.Println(s.Books.Books[shopIndex].Name)
		return s.Books.Books[shopIndex].ID
	} else if (shopIndex - len(s.Books.Books)) < len(s.Items.Items) {
		log.Println(s.Items.Items[shopIndex-len(s.Books.Books)].Name)
		return s.Items.Items[shopIndex-len(s.Books.Books)].ID
	}
	return uuid.New()
}

func (s *Shop) NextRow() {
	s.TableIndex++
	if s.TableIndex > s.TableLen {
		s.TableIndex = 1
	}
}

func (s *Shop) PreviousRow() {
	s.TableIndex--
	if s.TableIndex < 1 {
		s.TableIndex = 1
	}
}

func (s *Shop) Update() {
	if s.Modified.Add(time.Minute * 1).Before(time.Now()) {
		log.Println("Shop | Update()")
		*s = InitShop()
		s.Modified = time.Now()
	}
}

func (s *Shop) LoadShopTable() {
	columns := []string{"Name", "Description", "IQ Required", "Cost"}
	var rows [][]string
	n := 0
	for n < len(s.Books.Books) {
		rows = append(rows, []string{s.Books.Books[n].Name, s.Books.Books[n].Author, strconv.Itoa(s.Books.Books[n].IntelligenceRequirement), strconv.Itoa(s.Books.Books[n].KnowledgeRequirement)})
		n += 1
	}
	n = 0
	for n < len(s.Items.Items) {
		rows = append(rows, []string{s.Items.Items[n].Name, s.Items.Items[n].Description, strconv.Itoa(s.Items.Items[n].IqRequirement), strconv.Itoa(s.Items.Items[n].Cost)})
		n += 1
	}

	styleFunc := func(row, col int) lipgloss.Style {
		if row == 0 {
			return lipgloss.NewStyle().Padding(0, 1)
		}
		if col == 0 {
			return lipgloss.NewStyle().Width(18).Padding(1)
		}
		return lipgloss.NewStyle().Width(20).Padding(1, 1)
	}
	s.table = *table.New().Headers(columns...).Rows(rows...).StyleFunc(styleFunc)
	s.TableIndex = 1
	s.TableLen = len(rows)
}

//go:embed AllGameItems.bin
var agi []byte

func LoadAllGameItems() GameItemDatabase {
	buff := bytes.NewBuffer(agi)
	dec := gob.NewDecoder(buff)

	gameItems := GameItemDatabase{}

	dec_err := dec.Decode(&gameItems)
	if dec_err != nil {
		log.Println("LoadAllGameItems | dec.Decode")
		log.Println(gameItems)
		panic(dec_err)

	}
	return gameItems
}

func InitShop() Shop {
	n := 0
	var books Library
	var items GameItemDatabase

	columns := []string{"Name", "Description", "IQ Required", "Cost"}
	var rows [][]string
	for n < 4 {
		randomIndex := rand.Intn(len(AllBooksLibrary.Books) - 1)
		if !books.ContainsBook(AllBooksLibrary.Books[randomIndex]) {
			books.AddBookToLibrary(AllBooksLibrary.Books[randomIndex])
			b := AllBooksLibrary.Books[randomIndex]
			rows = append(rows, []string{b.Name, b.Author, strconv.Itoa(b.IntelligenceRequirement), strconv.Itoa(b.KnowledgeRequirement)})
			n += 1
		}
	}

	n = 0
	for n < 1 {
		randomIndex := rand.Intn(len(AllGameItems.Items) - 1)
		if !items.ContainsItem(AllGameItems.Items[randomIndex]) {
			items.AddItem(AllGameItems.Items[randomIndex])
			i := AllGameItems.Items[randomIndex]
			rows = append(rows, []string{i.Name, i.Description, strconv.Itoa(i.IqRequirement), strconv.Itoa(i.Cost)})
			n += 1
		}

	}

	styleFunc := func(row, col int) lipgloss.Style {
		if row == 0 {
			return lipgloss.NewStyle().Padding(0, 1)
		}
		if col == 0 || col == 1 {
			return lipgloss.NewStyle().Width(25).Padding(0)
		}
		return lipgloss.NewStyle().Width(5).Padding(0)
	}

	return Shop{
		Books:      books,
		Items:      items,
		ShopSize:   len(books.Books) + len(items.Items),
		Modified:   time.Now(),
		table:      *table.New().Headers(columns...).Rows(rows...).StyleFunc(styleFunc),
		TableIndex: 1,
		TableLen:   len(rows),
	}
}
