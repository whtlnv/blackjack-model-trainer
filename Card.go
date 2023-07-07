package main

type Suit int
const (
	Spades Suit = iota
	Hearts
	Diamonds
	Clubs
)

type Rank int
const (
	Ace = 1
	Two = 2
	Three = 3
	Four = 4
	Five = 5
	Six = 6
	Seven = 7
	Eight = 8
	Nine = 9
	Ten = 10
	Jack = 11
	Queen = 12
	King = 13
)

type Card struct {
	rank Rank 
	suit Suit
}

func (c Card) Value() []int {
	if c.rank > Ten {
		return []int{10}
	}
	if c.rank == Ace {
		return []int{1, 11}
	}
	return []int{int(c.rank)}
}