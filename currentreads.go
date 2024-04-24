package main

import (
	"strconv"

	"github.com/google/uuid"
)

type CurrentReads struct {
	BookIDs []uuid.UUID
}

func (cr *CurrentReads) String(preceding string) string {
	s := ""
	if preceding == "DIGITS" {
		for i, b := range cr.BookIDs {
			b_data, err := AllBooksLibrary.GetBookByID(b)
			if err == nil {
				s += strconv.Itoa(i+1) + ". " + b_data.Name + ", " + b_data.Author + "\n"
			}
		}
		return s
	}
	for _, b := range cr.BookIDs {
		b_data, err := AllBooksLibrary.GetBookByID(b)
		if err == nil {
			s += preceding + " " + b_data.Name + ", " + b_data.Author + "\n"
		}
	}

	return s
}

func (cr *CurrentReads) ReplaceBook(index int, book uuid.UUID) {
	cr.BookIDs[index] = book
}

func (cr *CurrentReads) ContainsBook(bookID uuid.UUID) bool {
	for _, b := range cr.BookIDs {
		if b == bookID {
			return true
		}
	}
	return false
}

func (cr *CurrentReads) AddBookToLibrary(bookID uuid.UUID) {
	cr.BookIDs = append(cr.BookIDs, bookID)
}
