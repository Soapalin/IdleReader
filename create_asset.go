package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"

	"github.com/google/uuid"
)

func CreateAllBookLibBin() {
	f, err := os.Create("AllBooksLibrary.bin")
	if err != nil {
		log.Println("CreateAllBookLibBin | os.Create")
		panic(err)
	}

	defer f.Close()

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)

	var abl Library = InitAllBookLibrary()

	error_enc := enc.Encode(abl)
	if error_enc != nil {
		log.Println(error_enc)
		panic(error_enc)
	}

	if _, err := f.Write(buff.Bytes()); err != nil {
		panic(err)
	}

}

func CreateAllGameItemBin() {
	f, err := os.Create("AllGameItems.bin")
	if err != nil {
		log.Println("CreateAllGameItemBin | os.Create")
		panic(err)
	}

	defer f.Close()

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)

	var gdb GameItemDatabase = InitAllGameItemDatabase()

	error_enc := enc.Encode(gdb)
	if error_enc != nil {
		log.Println(error_enc)
		panic(error_enc)
	}

	if _, err := f.Write(buff.Bytes()); err != nil {
		panic(err)
	}
}

func UpdateAllBooksLibrary()     {}
func UpdateAllGameItemDatabase() {}

func InitAllGameItemDatabase() GameItemDatabase {
	items := []Item{
		{
			ID:            uuid.New(),
			Name:          "Reading Glasses",
			Description:   "Increases Reading speed by 20%",
			Cost:          10000,
			IqRequirement: 1,
			Bought:        false,
			Effect:        nil,
		},
		{
			ID:            uuid.New(),
			Name:          "Bookmark",
			Description:   "Know where you left your reading.",
			Cost:          100,
			IqRequirement: 1,
			Bought:        false,
			Effect:        nil,
		},
		{
			ID:            uuid.New(),
			Name:          "Reading Light",
			Description:   "Everything is more clear. Your gain an additional 10% knowledge.",
			Cost:          100,
			IqRequirement: 1,
			Bought:        false,
			Effect:        nil,
		},
	}
	return GameItemDatabase{
		Items: items,
	}
}

func InitAllBookLibrary() Library {
	books := []Book{
		{
			ID:                      uuid.New(),
			Name:                    "The Poppy War",
			Author:                  "R.F Kuang",
			Progress:                0,
			KnowledgeIncrease:       300,
			KnowledgeRequirement:    150,
			IntelligenceIncrease:    2,
			IntelligenceRequirement: 70,
			Pages:                   500,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Babel",
			Author:                  "R.F Kuang",
			Progress:                0,
			KnowledgeIncrease:       300,
			KnowledgeRequirement:    150,
			IntelligenceIncrease:    3,
			IntelligenceRequirement: 90,
			Pages:                   700,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Tomorrow, and Tomorrow, and Tomorrow",
			Author:                  "Gabrielle Zevin",
			Progress:                0,
			KnowledgeIncrease:       200,
			KnowledgeRequirement:    100,
			IntelligenceIncrease:    2,
			IntelligenceRequirement: 80,
			Pages:                   450,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Intro to Electrical Engineering",
			Author:                  "Peter Mclean",
			Progress:                0,
			KnowledgeIncrease:       600,
			KnowledgeRequirement:    300,
			IntelligenceIncrease:    5,
			IntelligenceRequirement: 100,
			Pages:                   600,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Fundamentals of Electrical Engineering",
			Author:                  "Gang Lei",
			Progress:                0,
			KnowledgeIncrease:       600,
			KnowledgeRequirement:    300,
			IntelligenceIncrease:    5,
			IntelligenceRequirement: 100,
			Pages:                   600,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Child Book I",
			Author:                  "Unknown",
			Progress:                0,
			KnowledgeIncrease:       120,
			KnowledgeRequirement:    150,
			IntelligenceIncrease:    1,
			IntelligenceRequirement: 40,
			Pages:                   40,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Child Book II",
			Author:                  "Unknown",
			Progress:                0,
			KnowledgeIncrease:       120,
			KnowledgeRequirement:    150,
			IntelligenceIncrease:    1,
			IntelligenceRequirement: 40,
			Pages:                   40,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Poem Book I",
			Author:                  "Unknown",
			Progress:                0,
			KnowledgeIncrease:       120,
			KnowledgeRequirement:    150,
			IntelligenceIncrease:    1,
			IntelligenceRequirement: 40,
			Pages:                   20,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Knowledge Cheat Book",
			Author:                  "Lucien",
			Progress:                0,
			KnowledgeIncrease:       10000,
			KnowledgeRequirement:    1,
			IntelligenceIncrease:    1,
			IntelligenceRequirement: 1,
			Pages:                   20,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "IQ Cheat Book",
			Author:                  "Lucien",
			Progress:                0,
			KnowledgeIncrease:       10,
			KnowledgeRequirement:    150,
			IntelligenceIncrease:    100,
			IntelligenceRequirement: 40,
			Pages:                   20,
			Repeat:                  0,
		},
	}

	return Library{
		Books: books,
	}
}
