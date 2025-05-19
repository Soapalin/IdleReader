package main

import (
	"bytes"
	"encoding/gob"
	calculate "game/engine/utils"
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

// func UpdateAllBooksLibrary() {

// 	existingLibrary := LoadAllBooksLibrary()
// 	fullLibrary := AllBooksLibrary

// 	for _, book := range fullLibrary.Books {
// 		if !existingLibrary.ContainsBookByNameAndAuthor(book) {
// 			existingLibrary.AddBookToLibrary(book)
// 		}
// 	}

// 	f, err := os.Create("AllBooksLibrary.bin")
// 	if err != nil {
// 		log.Println("UpdateAllBooksLibrary | os.Create")
// 		panic(err)
// 	}

// 	defer f.Close()

// 	var buff bytes.Buffer
// 	enc := gob.NewEncoder(&buff)

// 	error_enc := enc.Encode(existingLibrary)
// 	if error_enc != nil {
// 		log.Println(error_enc)
// 		panic(error_enc)
// 	}

// 	if _, err := f.Write(buff.Bytes()); err != nil {
// 		panic(err)
// 	}

// }
// func UpdateAllGameItemDatabase() {

// 	existingItems := LoadAllGameItems()
// 	fullItems := InitAllGameItemDatabase()

// 	for _, item := range fullItems.Items {
// 		if !existingItems.ContainsItemByName(item) {
// 			existingItems.AddItem(item)
// 		}
// 	}

// 	f, err := os.Create("AllGameItems.bin")
// 	if err != nil {
// 		log.Println("CreateAllGameItemBin | os.Create")
// 		panic(err)
// 	}

// 	defer f.Close()

// 	var buff bytes.Buffer
// 	enc := gob.NewEncoder(&buff)

// 	error_enc := enc.Encode(fullItems)
// 	if error_enc != nil {
// 		log.Println(error_enc)
// 		panic(error_enc)
// 	}

// 	if _, err := f.Write(buff.Bytes()); err != nil {
// 		panic(err)
// 	}

// }

func InitAllGameItemDatabase() GameItemDatabase {
	items := []Item{
		{
			ID:            uuid.New(),
			Name:          "Reading Glasses",
			Description:   "Increases Reading speed by 20%",
			Cost:          10000,
			IqRequirement: 1,
			Bought:        false,
			Effect:        "INCREASE_READING_SPEED_20",
		},
		{
			ID:            uuid.New(),
			Name:          "Bookmark",
			Description:   "Know where you left your reading.",
			Cost:          100,
			IqRequirement: 1,
			Bought:        false,
			Effect:        "KEEP_PROGRESS_ON_EXIT",
		},
		{
			ID:            uuid.New(),
			Name:          "Reading Light",
			Description:   "Everything is clearer. Your gain an additional 10% knowledge.",
			Cost:          100,
			IqRequirement: 1,
			Bought:        false,
			Effect:        "INCREASE_KNOWLEDGE_10",
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
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(70, 500),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(70),
			IntelligenceIncrease:    2,
			IntelligenceRequirement: 70,
			Pages:                   500,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Dragon Republic",
			Author:                  "R.F Kuang",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(70, 500),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(70),
			IntelligenceIncrease:    2,
			IntelligenceRequirement: 70,
			Pages:                   500,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "The Burning God",
			Author:                  "R.F Kuang",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(70, 500),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(70),
			IntelligenceIncrease:    2,
			IntelligenceRequirement: 70,
			Pages:                   500,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Yellow Face",
			Author:                  "R.F Kuang",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(70, 500),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(70),
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
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(70, 700),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(70),
			IntelligenceIncrease:    3,
			IntelligenceRequirement: 90,
			Pages:                   700,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Dune",
			Author:                  "Frank Herbert",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(70, 700),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(70),
			IntelligenceIncrease:    3,
			IntelligenceRequirement: 90,
			Pages:                   700,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "A Court of Thorns and Roses",
			Author:                  "Sarah J. Maas",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(70, 700),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(70),
			IntelligenceIncrease:    3,
			IntelligenceRequirement: 90,
			Pages:                   700,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Nineteen Eighty-Four",
			Author:                  "George Orwell",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(70, 700),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(70),
			IntelligenceIncrease:    3,
			IntelligenceRequirement: 90,
			Pages:                   300,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "The Hunger Games",
			Author:                  "Suzanne Collins",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(50, 500),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(50),
			IntelligenceIncrease:    2,
			IntelligenceRequirement: 50,
			Pages:                   500,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Fourth Wing",
			Author:                  "Rebecca Yarros",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(50, 500),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(50),
			IntelligenceIncrease:    2,
			IntelligenceRequirement: 50,
			Pages:                   500,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Iron Flame",
			Author:                  "Rebecca Yarros",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(50, 500),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(50),
			IntelligenceIncrease:    2,
			IntelligenceRequirement: 50,
			Pages:                   500,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Harry Potter and The Ph. Stone",
			Author:                  "J.K Rowling",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(40, 500),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(40),
			IntelligenceIncrease:    2,
			IntelligenceRequirement: 40,
			Pages:                   500,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Harry Potter and The Ch. of Secrets",
			Author:                  "J.K Rowling",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(40, 500),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(40),
			IntelligenceIncrease:    2,
			IntelligenceRequirement: 40,
			Pages:                   500,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Harry Potter and The Pris. of Azkaban",
			Author:                  "J.K Rowling",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(40, 500),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(40),
			IntelligenceIncrease:    2,
			IntelligenceRequirement: 40,
			Pages:                   500,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Harry Potter and The Glob. of Fire",
			Author:                  "J.K Rowling",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(40, 500),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(40),
			IntelligenceIncrease:    2,
			IntelligenceRequirement: 40,
			Pages:                   500,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Romeo and Juliet",
			Author:                  "Shakespeare",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(60, 500),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(60),
			IntelligenceIncrease:    2,
			IntelligenceRequirement: 60,
			Pages:                   200,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Macbeth",
			Author:                  "Shakespeare",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(60, 500),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(60),
			IntelligenceIncrease:    2,
			IntelligenceRequirement: 60,
			Pages:                   500,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "The Legendary Mechanic",
			Author:                  "N/A",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(70, 500),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(70),
			IntelligenceIncrease:    2,
			IntelligenceRequirement: 70,
			Pages:                   500,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "The Author's POV",
			Author:                  "N/A",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(40, 500),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(40),
			IntelligenceIncrease:    2,
			IntelligenceRequirement: 40,
			Pages:                   500,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "I Shall Seal The Heavens",
			Author:                  "Er Gen",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(50, 500),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(50),
			IntelligenceIncrease:    2,
			IntelligenceRequirement: 50,
			Pages:                   500,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Beyond the Timescape",
			Author:                  "Er Gen",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(50, 500),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(50),
			IntelligenceIncrease:    2,
			IntelligenceRequirement: 50,
			Pages:                   500,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Desolate Era",
			Author:                  "I Eat Tomatoes",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(50, 500),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(50),
			IntelligenceIncrease:    2,
			IntelligenceRequirement: 50,
			Pages:                   500,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Tomorrow, and Tomorrow, and Tomorrow",
			Author:                  "Gabrielle Zevin",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(80, 450),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(80),
			IntelligenceIncrease:    2,
			IntelligenceRequirement: 80,
			Pages:                   450,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Elsewhere",
			Author:                  "Gabrielle Zevin",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(80, 450),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(80),
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
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(100, 600),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(100),
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
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(100, 600),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(100),
			IntelligenceIncrease:    5,
			IntelligenceRequirement: 100,
			Pages:                   600,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Signals and Systems",
			Author:                  "Peter McLean",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(100, 600),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(100),
			IntelligenceIncrease:    5,
			IntelligenceRequirement: 100,
			Pages:                   600,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Embedded Software Introduction",
			Author:                  "Peter McLean",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(100, 600),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(100),
			IntelligenceIncrease:    5,
			IntelligenceRequirement: 110,
			Pages:                   600,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Introduction to C Programming",
			Author:                  "Beeshanga",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(100, 600),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(100),
			IntelligenceIncrease:    5,
			IntelligenceRequirement: 110,
			Pages:                   600,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Diary of a Wimpy Kid",
			Author:                  "Jeff Kinney",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(40, 40),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(40),
			IntelligenceIncrease:    1,
			IntelligenceRequirement: 40,
			Pages:                   150,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Child Book I",
			Author:                  "Unknown",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(40, 40),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(40),
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
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(40, 40),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(40),
			IntelligenceIncrease:    1,
			IntelligenceRequirement: 40,
			Pages:                   40,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Child Book III",
			Author:                  "Unknown",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(40, 40),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(40),
			IntelligenceIncrease:    1,
			IntelligenceRequirement: 40,
			Pages:                   40,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Child Book IV",
			Author:                  "Unknown",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(40, 40),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(40),
			IntelligenceIncrease:    1,
			IntelligenceRequirement: 40,
			Pages:                   40,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Atomic Habits",
			Author:                  "James Clear",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(40, 200),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(40),
			IntelligenceIncrease:    1,
			IntelligenceRequirement: 40,
			Pages:                   200,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "The Subtle Art of Not Giving a F*ck",
			Author:                  "Mark Hanson",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(40, 200),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(40),
			IntelligenceIncrease:    1,
			IntelligenceRequirement: 40,
			Pages:                   200,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Rich Dad, Poor Dad",
			Author:                  "Robert T. Kiyosaki",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(40, 200),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(40),
			IntelligenceIncrease:    1,
			IntelligenceRequirement: 40,
			Pages:                   200,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "The Alchemist",
			Author:                  "Paulo Coelho",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(40, 200),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(40),
			IntelligenceIncrease:    1,
			IntelligenceRequirement: 40,
			Pages:                   200,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Poem Book I",
			Author:                  "Unknown",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(40, 20),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(40),
			IntelligenceIncrease:    1,
			IntelligenceRequirement: 40,
			Pages:                   20,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Poem Book II",
			Author:                  "Unknown",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(40, 20),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(40),
			IntelligenceIncrease:    1,
			IntelligenceRequirement: 40,
			Pages:                   20,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Poem Book III",
			Author:                  "Unknown",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(40, 20),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(40),
			IntelligenceIncrease:    1,
			IntelligenceRequirement: 40,
			Pages:                   20,
			Repeat:                  0,
		},
		{
			ID:                      uuid.New(),
			Name:                    "Poem Book IV",
			Author:                  "Unknown",
			Progress:                0,
			KnowledgeIncrease:       calculate.CalculateKnowledgeIncrease(40, 20),
			KnowledgeRequirement:    calculate.CalculateKnowledgeCost(40),
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
