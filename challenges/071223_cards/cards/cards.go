package cards

import (
	"cmp"
	"context"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var faceCards = map[rune]int{
	'A': 14,
	'K': 13,
	'Q': 12,
	'J': 11,
	'T': 10,
}

type hand struct {
	cards                 []int
	bid                   int
	totalScore            int
	typeFoundStr, ogValue string
}

func (h hand) ToCSVRow() []string {
	cardIntStr := make([]string, 0, 4)
	for _, c := range h.cards {
		cardIntStr = append(cardIntStr, strconv.Itoa(c))
	}
	return []string{
		strings.Join(cardIntStr, " "),
		strconv.Itoa(h.bid),
		h.ogValue,
		h.typeFoundStr,
		strconv.Itoa(h.totalScore),
	}
}

func (h hand) MarshalZerologObject(e *zerolog.Event) {
	e.
		Ints("int_cards", h.cards).
		Int("bid", h.bid).
		Int("total_score", h.totalScore).
		Str("type_found", h.typeFoundStr).
		Str("og_value", h.ogValue)
}
func CalculatePartOne(ctx context.Context, lines []string) int {
	l := log.Ctx(ctx).With().Logger()

	var parsedHands []hand
	for _, line := range lines {
		cards, bid := separateCardsBid(line)
		parsedHands = append(parsedHands, createHand(cards, bid, line))
	}

	slices.SortStableFunc(parsedHands, func(a, b hand) int { return cmp.Compare(a.totalScore, b.totalScore) })
	var totalWinnings int

	for idx, hand := range parsedHands {
		// l.Info().Object("hand", hand).Msg("hands ordered")
		totalWinnings += ((idx + 1) * hand.bid)
	}
	rows := make([][]string, len(parsedHands)+1)

	rows[0] = []string{
		"og_value",
		"int_cards",
		"bid",
		"total_score",
		"type_found_desc",
	}
	l.Info().Msg(strings.Join(rows[0], ","))

	for i := len(parsedHands) - 1; i >= 0; i-- {
		rows[i+1] = parsedHands[i].ToCSVRow()
		l.Info().Msg(strings.Join(rows[i+1], ","))
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

func separateCardsBid(line string) ([]rune, int) {
	cardsBids := strings.Fields(line)
	cardsStr, bidsStr := cardsBids[0], cardsBids[1]
	bid, err := strconv.Atoi(bidsStr)
	if err != nil {
		panic([]any{err, bidsStr})
	}

	return []rune(cardsStr), bid
}

func createHand(cards []rune, bid int, ogValue string) hand {
	// TODO: Left off here, must summarize the key factor in the score
	totalCardOrderScore, cardToOccurrences := createCardOrderScoreAndMap(cards)
	typeScore, typeFound := createTypeScore(cardToOccurrences)
	cardInts := make([]int, 0, len(cards))
	for _, card := range cards {
		cardInts = append(cardInts, runeToCardValue(card))
	}
	return hand{
		cards:        cardInts,
		bid:          bid,
		totalScore:   totalCardOrderScore + typeScore,
		ogValue:      ogValue,
		typeFoundStr: typeFound,
	}
}

// already looping through the cards, so just create the map for the type score later
func createCardOrderScoreAndMap(cards []rune) (int, map[rune]int) {
	var handCardOrderScore int
	cardCount := make(map[rune]int, 5)
	for idx, card := range cards {
		idxInArr := len(cards) - idx - 1

		cardValue := runeToCardValue(card)
		// must reverse the idx,
		// as createOrderScore increases the importance from left to right in the array,
		// and in this case the left most value is the most important value
		cardOrderScore, err := createOrderScore(
			createOrderScoreInput{
				numToScore:         cardValue,
				idxInArray:         idxInArr,
				minAllowedValue:    2,
				maxAllowedNumValue: 14,
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

func createTypeScore(cardCount map[rune]int) (int, string) {
	var firstPair, secondPair int
	var firstPairCard, secondPairCard rune
	keys := strings.Builder{}
	for card, count := range cardCount {
		keys.WriteRune(card)
		if count < 2 {
			continue
		}
		if firstPair > 1 {
			secondPair = count
			secondPairCard = card
			continue
		}
		firstPairCard = card
		firstPair = count
	}
	var typeOrder int
	var typeFound string
	switch {
	case max(firstPair, secondPair) == 5:
		typeOrder = 13
		typeFound = fmt.Sprintf("five of a kind - %s", string(firstPairCard))
	case max(firstPair, secondPair) == 4:
		typeOrder = 12
		typeFound = fmt.Sprintf("four of a kind - %s", string(firstPairCard))
	case max(firstPair, secondPair) == 3 && min(firstPair, secondPair) == 2:
		typeOrder = 11
		typeFound = fmt.Sprintf("full house - %s %s", string(firstPairCard), string(secondPairCard))
	case max(firstPair, secondPair) == 3:
		typeOrder = 10
		typeFound = fmt.Sprintf("three of a kind - %s", string(firstPairCard))
	case firstPair == 2 && secondPair == 2:
		typeOrder = 9
		typeFound = fmt.Sprintf("two pair - %s %s", string(firstPairCard), string(secondPairCard))
	case max(firstPair, secondPair) == 2:
		typeOrder = 8
		typeFound = fmt.Sprintf("pair - %s", string(firstPairCard))
	default:
		typeFound = fmt.Sprintf("card high - %s", keys.String())
		typeOrder = 0
	}

	// should be at the same scale as the order score for the card values,
	// hence starting at 14
	typeScore, err := createOrderScore(
		createOrderScoreInput{
			numToScore:         typeOrder,
			idxInArray:         5, // max card value place is 4, so one more than that
			maxAllowedNumValue: 12,
		},
	)
	if err != nil {
		panic(err)
	}

	return typeScore, typeFound
}

type createOrderScoreInput struct {
	numToScore, idxInArray              int
	minAllowedValue, maxAllowedNumValue int
}

func createOrderScore(input createOrderScoreInput) (int, error) {

	numPossibleValues := input.maxAllowedNumValue - input.minAllowedValue

	// normalize number where smallest is 0
	numTranslated := input.numToScore - input.minAllowedValue
	base := float64(numPossibleValues + 1)

	score := math.Pow(
		base, float64(input.idxInArray),
	)
	idxArrMult := math.Pow(base, float64(input.idxInArray))
	score += idxArrMult * (float64(numTranslated) - 1)

	return int(score), nil
}
