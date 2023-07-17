package blackjack

import "github.com/samber/lo"

type Hand []Card
type HandScore struct {
	Low  int
	High int
}

func (hand *Hand) Deal(card Card) {
	*hand = append(*hand, card)
}

func (hand *Hand) Values() []int {
	var calculateSum func(cards []Card, index int, currentSum int) []int

	calculateSum = func(cards []Card, index int, currentSum int) []int {
		if index >= len(cards) {
			return []int{currentSum}
		}

		card := cards[index]
		values := card.Value()
		var sums []int
		for _, value := range values {
			moreSums := calculateSum(cards, index+1, currentSum+value)
			sums = append(sums, moreSums...)
		}

		return sums
	}

	return lo.Uniq(calculateSum(*hand, 0, 0))
}

func (hand *Hand) Score() (score HandScore, isBusted bool) {
	values := hand.Values()

	notBusted := lo.Filter(values, func(value int, _ int) bool {
		return value <= 21
	})

	// TODO: remove this when confident
	if len(notBusted) > 2 {
		panic("A hand should never have more than 2 values")
	}

	if len(notBusted) > 0 {
		score = HandScore{Low: lo.Min(notBusted), High: lo.Max(notBusted)}
		isBusted = false
	} else {
		score = HandScore{Low: lo.Min(values), High: lo.Max(values)}
		isBusted = true
	}

	return score, isBusted
}

func (hand *Hand) IsPair() bool {
	return (*hand)[0].rank == (*hand)[1].rank
}

func (hand *Hand) HasSoftValue() bool {
	score, _ := hand.Score()
	return score.Low != score.High
}
