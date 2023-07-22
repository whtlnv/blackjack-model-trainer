package blackjack

import (
	"fmt"
	"math/rand"
)

type Shoeish interface {
	Size() int
	Peek(count int) []Card
}

type Shoe struct {
	decks  int
	cards  []Card
	cursor int
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
	shoe := &Shoe{decks, []Card{}, 0}

	shoe.build()
	shoe.shuffle()

	return shoe
}

// Public methods

func (shoe *Shoe) Size() int {
	return len(shoe.cards)
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
	return shoe.cursor, nil
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

func (shoe *Shoe) shuffle() {
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
}
