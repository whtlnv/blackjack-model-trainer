package blackjack

import (
	"fmt"
	"math/rand"
)

type Shoeish interface {
	Size() int
	Shuffle()
	Peek(count int) []Card
	AdvanceCursor(offset int) (int, error)
	SetPenetration(deckPercentage float64)
	NeedsReshuffle() bool
}

type Shoe struct {
	decks            int
	cards            []Card
	cursor           int
	penetrationIndex int
	needsReshuffle   bool
}

type CursorOutOfBoundsError struct {
	cursor int
	offset int
	size   int
}

func (e *CursorOutOfBoundsError) Error() string {
	return fmt.Sprintf(
		"Shoe cursor out of bounds. Shoe of size %d had cursor at %d, and received offset of %d",
		e.size,
		e.cursor,
		e.offset,
	)
}

// Factory

func NewShoe(decks int) *Shoe {
	shoe := &Shoe{
		decks:            decks,
		cards:            []Card{},
		cursor:           0,
		penetrationIndex: decks * 52,
		needsReshuffle:   false,
	}

	shoe.build()
	shoe.Shuffle()

	return shoe
}

// Public methods

func (shoe *Shoe) Size() int {
	return len(shoe.cards)
}

func (shoe *Shoe) Shuffle() {
	// Fisher-Yates shuffle
	for i := range shoe.cards {
		j := rand.Intn(i + 1)
		shoe.cards[i], shoe.cards[j] = shoe.cards[j], shoe.cards[i]
	}

	// Cut the deck
	firstHalf := shoe.cards[:len(shoe.cards)/2]
	secondHalf := shoe.cards[len(shoe.cards)/2:]
	shoe.cards = append(secondHalf, firstHalf...)

	// shuffle
	for i := range shoe.cards {
		j := rand.Intn(shoe.Size() - 1)
		shoe.cards[i], shoe.cards[j] = shoe.cards[j], shoe.cards[i]
	}

	shoe.needsReshuffle = false
	shoe.cursor = 0
}

func (shoe *Shoe) Peek(count int) []Card {
	cursor := shoe.cursor

	if cursor+count > len(shoe.cards) {
		count = len(shoe.cards) - cursor
	}
	return shoe.cards[cursor : cursor+count]
}

func (shoe *Shoe) AdvanceCursor(offset int) (int, error) {
	advanceTo := shoe.cursor + offset
	if advanceTo > len(shoe.cards) {
		return shoe.cursor, &CursorOutOfBoundsError{shoe.cursor, offset, len(shoe.cards)}
	}

	shoe.cursor = advanceTo

	if shoe.cursor >= shoe.penetrationIndex {
		shoe.needsReshuffle = true
	}

	return shoe.cursor, nil
}

func (shoe *Shoe) SetPenetration(deckPercentage float64) {
	shoe.penetrationIndex = int(float64(shoe.Size()) * deckPercentage)
}

func (shoe *Shoe) NeedsReshuffle() bool {
	return shoe.needsReshuffle
}

// Private methods

func (shoe *Shoe) build() {
	for i := 0; i < shoe.decks; i++ {
		for _, suit := range SUITS {
			for _, rank := range RANKS {
				card := NewCard(rank, suit)
				shoe.cards = append(shoe.cards, card)
			}
		}
	}
}
