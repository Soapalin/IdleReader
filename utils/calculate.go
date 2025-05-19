package utils

const (
	InitialCost = 100
	InitialIQ   = 40
)

func CalculateKnowledgeCost(IQ int) int {
	knowledgeCost := InitialCost
	for i := InitialIQ; i < IQ; i++ {
		knowledgeCost = int(float32(knowledgeCost)*1.05 + 100)
	}
	return knowledgeCost

}

func CalculateKnowledgeIncrease(bookIQ int, pages int) int {
	return (bookIQ * pages) / 10
}
