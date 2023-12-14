package cards

import (
	"cmp"
	"context"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
)

var faceCards = map[rune]int{
	'A': 14,
	'K': 13,
	'Q': 12,
	'J': 11,
	'T': 10,
}

type hand struct {
	cards      []int
	bid        int
	totalScore int
	ogValue    string
}

func CalculatePartOne(ctx context.Context, lines []string) int {
	var hands []hand
	for _, line := range lines {
		cards, bid := separateCardsBid(line)
		hands = append(hands, createHand(cards, bid, line))
	}
	// multiply by -1 to ensure worst hand is at the top
	slices.SortFunc(hands, func(a, b hand) int { return cmp.Compare(a.totalScore, b.totalScore) })
	var totalWinnings int
	// todo: left off here, type ordering is still not working
	for idx, hand := range hands {
		totalWinnings += ((idx + 1) * hand.bid)
	}
	return totalWinnings
}

func CalculatePartTwo(ctx context.Context, input []string) int {
	return 0
}

func runeToCardValue(cardRune rune) int {
	card := cardRune - '0'
	if card > 9 {
		// it's a face card
		card, ok := faceCards[cardRune]
		if !ok {
			panic(cardRune)
		}
		return card
	}
	return int(card)
}

func separateCardsBid(line string) ([]int, int) {
	cardsBids := strings.Fields(line)
	cardsStr, bidsStr := cardsBids[0], cardsBids[1]
	bid, err := strconv.Atoi(bidsStr)
	if err != nil {
		panic([]any{err, bidsStr})
	}

	cards := make([]int, 0, len(cardsStr))

	for _, cardRune := range cardsStr {
		cards = append(cards, runeToCardValue(cardRune))
	}
	return cards, bid
}

func createHand(cards []int, bid int, ogValue string) hand {

	totalCardOrderScore, cardToOccurrences := createCardOrderScoreAndMap(cards)
	typeScore := createTypeScore(cardToOccurrences)
	return hand{
		cards:      cards,
		bid:        bid,
		totalScore: totalCardOrderScore + typeScore,
		ogValue:    ogValue,
	}
}

// already looping through the cards, so just create the map for the type score later
func createCardOrderScoreAndMap(cards []int) (int, map[int]int) {
	var handCardOrderScore int
	cardCount := make(map[int]int, 5)
	for idx, card := range cards {
		// must reverse the idx,
		// as createOrderScore increases the importance from left to right in the array,
		// and in this case the left most value is the most important value
		cardOrderScore, err := createOrderScore(
			createOrderScoreInput{
				numToScore:         cards[idx],
				idxInArray:         len(cards) - idx - 1,
				minAllowedValue:    2,
				maxAllowedNumValue: 14,
				base:               3,
			})
		if err != nil {
			panic(err)
		}
		handCardOrderScore += cardOrderScore

		if existing, ok := cardCount[card]; ok {
			cardCount[card] = existing + 1
			continue
		}
		cardCount[card] = 1
	}
	return handCardOrderScore, cardCount
}

func createTypeScore(cardCount map[int]int) int {
	var firstPair, secondPair int
	for _, count := range cardCount {
		if count < 2 {
			continue
		}
		if count > 1 && firstPair > 1 {
			secondPair = count
			continue
		}
		firstPair = count
	}
	var typeOrder int
	switch {
	case firstPair == 5:
		typeOrder = 13
	case firstPair == 4:
		typeOrder = 12
	case firstPair == 3 && secondPair == 2:
		typeOrder = 11
	case firstPair == 3:
		typeOrder = 10
	case firstPair == 2 && secondPair == 2:
		typeOrder = 9
	case firstPair == 2:
		typeOrder = 8
	default:
		typeOrder = 7
	}

	// should be at the same scale as the order score for the card values,
	// hence starting at 14
	typeScore, err := createOrderScore(
		createOrderScoreInput{
			numToScore:         typeOrder,
			idxInArray:         5, // max card value place is 4, so one more than that
			maxAllowedNumValue: 12,
			base:               3,
		},
	)
	if err != nil {
		panic(err)
	}

	return typeScore
}

type createOrderScoreInput struct {
	numToScore, idxInArray              int
	minAllowedValue, maxAllowedNumValue int
	base                                int
}

// createOrderScore uses the principle of:
// b^n > b^(n-1)*(n-1) for b^(n-1) < n < (b^n)-1
// and
// b^n = b^(n-1)*(n)   for b^(n-1) < n < b^n
//
// TODO: flesh out explainer
//
// each number in the array must be one of:
//   - a set of consecutive numbers starting from some number,
//     up to the possible number of ordered values (see below)
//   - eg. only possible numbers are 2 through 6, so there are 5 possible numbers
//
// the idx in array assumes the further right in the array, the higher the order score should be
//   - eg. index 5 should be scored higher than index 4
//
// the index for when a value is equal to the next power is always b^(n-1)*b
// this means there are only k-1 valid consecutive values
//
//   - if k were 4 and n were 2, you would only be able to use this powers principle if there were up to 16 valid values
//
//   - if k were 2 and n were 3, you would only be able to use this powers principle if there were 4 valid values
//
//   - this can also be represented as m*b*(n-1)*b
//     -- here m is the multiplier, b is the base, and n again is the i
//
// use the below google sheet to play with calculations / numbers
// https://docs.google.com/spreadsheets/d/1Hhnie7jD6O-1ZY7OxhfSHz_ooTdWfPHSzwI3SN5AiLQ/edit#gid=129553777
func createOrderScore(input createOrderScoreInput) (int, error) {
	if input.base < 3 {
		return 0, fmt.Errorf("base must be greater than 2. Received %d", input.base)
	}

	numPossibleValues := input.maxAllowedNumValue - input.minAllowedValue

	// normalize number where smallest is 0
	numTranslated := input.numToScore - input.minAllowedValue

	scaleFactor := input.base * numPossibleValues / (input.base - 1)

	numToScoreScaled := float64(scaleFactor)*(math.Pow(float64(input.base), float64(input.idxInArray))) +
		(float64(numTranslated) * float64(input.base))

	return int(numToScoreScaled), nil
}
