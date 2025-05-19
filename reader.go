package main

import (
	"errors"
	"github.com/google/uuid"
	"log"
	"strconv"
)

var (
	MAX_SIMULT_BOOK = 3
)

type ReadingSpeed struct {
	Word    int
	Line    int
	Page    int
	Chapter int
	Book    int
}

type Reader struct {
	ID               uuid.UUID
	Name             string
	IQ               int
	Fun              int
	Knowledge        int
	FavouriteBook    string
	FavouriteAuthor  string
	Prestige         int
	ReadCapacity     int
	CurrentReads     CurrentReads
	CurrentReadLimit int
	ReadingSpeed     int
	Library          Library
	Inventory        GameItemDatabase
}

func (r *Reader) ReadingIncrease() {

}
func (r *Reader) IncreaseProgress(book *Book) {
	book.Progress += float64(r.ReadingSpeed) / float64(book.Pages)
}

func (r *Reader) SwitchBook(index int, bookID uuid.UUID) error {
	if index > (r.CurrentReadLimit - 1) {
		return errors.New("index > current read limit")
	}

	return nil
}

func (r *Reader) DecreaseKnowledge(amount int) {
	r.Knowledge = r.Knowledge - amount
}

func (r *Reader) FinishedBook(book *Book) {
	if book.Repeat == 0 {
		r.IncreaseIQ(book.IntelligenceIncrease)
	}
	log.Println("FinishedBook | Repeat: " + strconv.Itoa(book.Repeat))
	book.Repeat += 1
	r.IncreaseKnowledge(book.KnowledgeIncrease)
	book.Progress = 0

}

func (r *Reader) IncreaseKnowledge(amount int) {
	log.Println("IncreaseKnowledge | " + strconv.Itoa(amount))
	r.Knowledge += amount
}

func (r *Reader) IncreaseIQ(amount int) {
	log.Println("IncreaseIQ | " + strconv.Itoa(amount))
	old_IQ := r.IQ
	r.IQ += amount
	if old_IQ < 100 && r.IQ >= 100 {
		r.CurrentReadLimit = r.CurrentReadLimit + 1
	}
	if old_IQ < 200 && r.IQ >= 200 {
		r.CurrentReadLimit = r.CurrentReadLimit + 1
	}
	r.CalculateReadingSpeed()
}

func (r *Reader) CalculateReadingSpeed() {
	if r.IQ < 60 {
		r.ReadingSpeed = 5
	} else if r.IQ >= 225 {
		r.ReadingSpeed = 225 / 12
	} else {
		r.ReadingSpeed = r.IQ / 12
	}
}

func (r *Reader) ActivatePrestige() {
	r.Prestige += 1
}

func (r *Reader) IQ_Title() string {
	if r.IQ < 70 {
		return "Extremely Low"
	} else if r.IQ < 80 {
		return "Very Low"
	} else if r.IQ < 90 {
		return "Low Average"
	} else if r.IQ < 110 {
		return "Average"
	} else if r.IQ < 120 {
		return "High Average"
	} else if r.IQ < 130 {
		return "Very High"
	} else if r.IQ < 150 {
		return "Extremely High"
	} else if r.IQ < 200 {
		return "Smarty Pants"
	} else if r.IQ < 400 {
		return "Big Brain"
	} else if r.IQ < 600 {
		return "Planetary Brain"
	} else {
		return "Galaxy Brain"
	}
}
