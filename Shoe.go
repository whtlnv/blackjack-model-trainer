package main

import (
	"math/rand"
)

type Shoe struct {
	decks int
	cards []Card
}

// Factory

func NewShoe(decks int) *Shoe {
	shoe := &Shoe{decks, []Card{}}

	shoe.build()
	shoe.shuffle()

	return shoe
}

// Public methods

func (shoe *Shoe) Size() int {
	return len(shoe.cards)
}

func (shoe *Shoe) Peek(cursor int, count int) []Card {
	if cursor+count > len(shoe.cards) {
		count = len(shoe.cards) - cursor
	}
	return shoe.cards[cursor : cursor+count]
}

// Private methods

func (shoe *Shoe) build() {
	for i := 0; i < shoe.decks; i++ {
		for _, suit := range SUITS {
			for _, rank := range RANKS {
				card := Card{rank, suit}
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
