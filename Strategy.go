package main

import "fmt"

type Strategy struct {
	PairMap map[PlayerHand]map[DealerHand]PlayerAction
	SoftMap map[PlayerHand]map[DealerHand]PlayerAction
	HardMap map[PlayerHand]map[DealerHand]PlayerAction
}

// Factory

func NewStrategy(raw []byte) (strategy Strategy, err error) {
	strategy = Strategy{}

	err = strategy.parseActionMap(raw)

	if err != nil {
		return strategy, err
	}

	return strategy, nil
}

// Public methods

func (strategy *Strategy) Play(playerHand Hand, dealerHand Hand) PlayerAction {
	dealerScore, _ := dealerHand.Score()
	playerScore, _ := playerHand.Score()

	if playerHand.IsPair() {
		return strategy.PairMap[PlayerHand(playerScore.Hard)][DealerHand(dealerScore.Hard)]
	}

	if playerHand.HasSoftValue() {
		return strategy.SoftMap[PlayerHand(playerScore.Hard)][DealerHand(dealerScore.Soft)]
	}

	return strategy.HardMap[PlayerHand(playerScore.Hard)][DealerHand(dealerScore.Hard)]
}

// Private methods

func (strategy *Strategy) parseActionMap(raw []byte) error {
	if len(raw) != HandCount {
		return fmt.Errorf("expected strategy length to be %d, got %d", HandCount, len(raw))
	}

	rawHardMapStartsAt := 0
	rawSoftMapStartsAt := rawHardMapStartsAt + (DealerHandCount * PlayerHardHandCount)
	rawPairMapStartsAt := rawSoftMapStartsAt + (DealerHandCount * PlayerSoftHandCount)

	rawHardMap := raw[rawHardMapStartsAt:rawSoftMapStartsAt]
	rawSoftMap := raw[rawSoftMapStartsAt:rawPairMapStartsAt]
	rawPairMap := raw[rawPairMapStartsAt:]

	strategy.HardMap = strategy.parseHardMap(rawHardMap)
	strategy.SoftMap = strategy.parseSoftMap(rawSoftMap)
	strategy.PairMap = strategy.parsePairMap(rawPairMap)

	return nil
}

func (strategy *Strategy) parseHardMap(slicedRaw []byte) map[PlayerHand]map[DealerHand]PlayerAction {
	columns := [DealerHandCount]DealerHand{DealerTwo, DealerThree, DealerFour, DealerFive, DealerSix, DealerSeven, DealerEight, DealerNine, DealerTen, DealerAce}
	rows := [PlayerHardHandCount]PlayerHand{HardFive, HardSix, HardSeven, HardEight, HardNine, HardTen, HardEleven, HardTwelve, HardThirteen, HardFourteen, HardFifteen, HardSixteen, HardSeventeen, HardEighteen, HardNineteen, HardTwenty}

	return stringToMap(slicedRaw, columns[:], rows[:])
}

func (strategy *Strategy) parseSoftMap(slicedRaw []byte) map[PlayerHand]map[DealerHand]PlayerAction {
	columns := [DealerHandCount]DealerHand{DealerTwo, DealerThree, DealerFour, DealerFive, DealerSix, DealerSeven, DealerEight, DealerNine, DealerTen, DealerAce}
	rows := [PlayerSoftHandCount]PlayerHand{SoftThirteen, SoftFourteen, SoftFifteen, SoftSixteen, SoftSeventeen, SoftEighteen, SoftNineteen, SoftTwenty}

	return stringToMap(slicedRaw, columns[:], rows[:])
}

func (strategy *Strategy) parsePairMap(slicedRaw []byte) map[PlayerHand]map[DealerHand]PlayerAction {
	columns := [DealerHandCount]DealerHand{DealerTwo, DealerThree, DealerFour, DealerFive, DealerSix, DealerSeven, DealerEight, DealerNine, DealerTen, DealerAce}
	rows := [PlayerPairHandCount]PlayerHand{PairTwos, PairThrees, PairFours, PairFives, PairSixes, PairSevens, PairEights, PairNines, PairTens, PairAces}

	return stringToMap(slicedRaw, columns[:], rows[:])
}

// Helper methods

func stringToMap(raw []byte, columns []DealerHand, rows []PlayerHand) map[PlayerHand]map[DealerHand]PlayerAction {
	parsedMap := make(map[PlayerHand]map[DealerHand]PlayerAction)

	cursor := 0
	for _, row := range rows {
		parsedMap[row] = make(map[DealerHand]PlayerAction)

		for _, column := range columns {
			parsedMap[row][column] = PlayerAction(raw[cursor])
			cursor += 1
		}
	}

	return parsedMap
}
