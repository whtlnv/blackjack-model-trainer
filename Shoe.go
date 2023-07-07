package main

type Shoe struct {
	Decks int
	cards []Card
}

func NewShoe(decks int) *Shoe {
	s := &Shoe{decks, []Card{}}
	for i := 0; i < decks; i++ {
		for _, suit := range SUITS {
			for _, rank := range RANKS {
				s.cards = append(s.cards, Card{rank, suit})
			}
		}
	}
	return s
}

func (s *Shoe) Size() int {
	return len(s.cards)
}