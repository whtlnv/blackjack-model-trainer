package main

type DealerHand int

const (
	DealerAce   DealerHand = 11
	DealerTen   DealerHand = 10
	DealerNine  DealerHand = 9
	DealerEight DealerHand = 8
	DealerSeven DealerHand = 7
	DealerSix   DealerHand = 6
	DealerFive  DealerHand = 5
	DealerFour  DealerHand = 4
	DealerThree DealerHand = 3
	DealerTwo   DealerHand = 2
)

type PlayerHand int

const (
	HardTwenty    PlayerHand = 20
	HardNineteen  PlayerHand = 19
	HardEighteen  PlayerHand = 18
	HardSeventeen PlayerHand = 17
	HardSixteen   PlayerHand = 16
	HardFifteen   PlayerHand = 15
	HardFourteen  PlayerHand = 14
	HardThirteen  PlayerHand = 13
	HardTwelve    PlayerHand = 12
	HardEleven    PlayerHand = 11
	HardTen       PlayerHand = 10
	HardNine      PlayerHand = 9
	HardEight     PlayerHand = 8
	HardSeven     PlayerHand = 7
	HardSix       PlayerHand = 6
	HardFive      PlayerHand = 5
	// HardFour 		PlayerHand = 4 // This equals {2,2} -> PairTwos, which takes precedence
)

const (
	SoftTwenty    PlayerHand = 20
	SoftNineteen  PlayerHand = 19
	SoftEighteen  PlayerHand = 18
	SoftSeventeen PlayerHand = 17
	SoftSixteen   PlayerHand = 16
	SoftFifteen   PlayerHand = 15
	SoftFourteen  PlayerHand = 14
	SoftThirteen  PlayerHand = 13
	// SoftTwelve    PlayerHand = 12 // This equals {A,A} -> PairAces, which takes precedence
)

const (
	PairTens   PlayerHand = 20
	PairNines  PlayerHand = 18
	PairEights PlayerHand = 16
	PairSevens PlayerHand = 14
	PairSixes  PlayerHand = 12
	PairFives  PlayerHand = 10
	PairFours  PlayerHand = 8
	PairThrees PlayerHand = 6
	PairTwos   PlayerHand = 4
	PairAces   PlayerHand = 11
)

type PlayerAction string

const (
	Hit    PlayerAction = "H"
	Stay   PlayerAction = "S"
	Double PlayerAction = "D"
	Split  PlayerAction = "P"
)

type Strategy struct {
	PairMap map[PlayerHand]map[DealerHand]PlayerAction
	SoftMap map[PlayerHand]map[DealerHand]PlayerAction
	HardMap map[PlayerHand]map[DealerHand]PlayerAction
}

// Factory

func NewStrategy(raw []byte) Strategy {
	strategy := Strategy{}

	strategy.parseActionMap(raw)

	return strategy
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

func (strategy *Strategy) parseActionMap(raw []byte) {
	dealerHandCount := 10
	playerHardHandCount := 16
	playerSoftHandCount := 8
	// playerPairHandCount := 10

	rawHardMapStartsAt := 0
	rawSoftMapStartsAt := rawHardMapStartsAt + (dealerHandCount * playerHardHandCount)
	rawPairMapStartsAt := rawSoftMapStartsAt + (dealerHandCount * playerSoftHandCount)

	rawHardMap := raw[rawHardMapStartsAt:rawSoftMapStartsAt]
	rawSoftMap := raw[rawSoftMapStartsAt:rawPairMapStartsAt]
	rawPairMap := raw[rawPairMapStartsAt:]

	strategy.HardMap = strategy.parseHardMap(rawHardMap)
	strategy.SoftMap = strategy.parseSoftMap(rawSoftMap)
	strategy.PairMap = strategy.parsePairMap(rawPairMap)
}

func (strategy *Strategy) parseHardMap(slicedRaw []byte) map[PlayerHand]map[DealerHand]PlayerAction {
	columns := [10]DealerHand{DealerTwo, DealerThree, DealerFour, DealerFive, DealerSix, DealerSeven, DealerEight, DealerNine, DealerTen, DealerAce}
	rows := [16]PlayerHand{HardFive, HardSix, HardSeven, HardEight, HardNine, HardTen, HardEleven, HardTwelve, HardThirteen, HardFourteen, HardFifteen, HardSixteen, HardSeventeen, HardEighteen, HardNineteen, HardTwenty}

	return stringToMap(slicedRaw, columns[:], rows[:])
}

func (strategy *Strategy) parseSoftMap(slicedRaw []byte) map[PlayerHand]map[DealerHand]PlayerAction {
	columns := [10]DealerHand{DealerTwo, DealerThree, DealerFour, DealerFive, DealerSix, DealerSeven, DealerEight, DealerNine, DealerTen, DealerAce}
	rows := [8]PlayerHand{SoftThirteen, SoftFourteen, SoftFifteen, SoftSixteen, SoftSeventeen, SoftEighteen, SoftNineteen, SoftTwenty}

	return stringToMap(slicedRaw, columns[:], rows[:])
}

func (strategy *Strategy) parsePairMap(slicedRaw []byte) map[PlayerHand]map[DealerHand]PlayerAction {
	columns := [10]DealerHand{DealerTwo, DealerThree, DealerFour, DealerFive, DealerSix, DealerSeven, DealerEight, DealerNine, DealerTen, DealerAce}
	rows := [10]PlayerHand{PairTwos, PairThrees, PairFours, PairFives, PairSixes, PairSevens, PairEights, PairNines, PairTens, PairAces}

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