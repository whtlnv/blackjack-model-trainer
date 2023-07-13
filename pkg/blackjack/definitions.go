package blackjack

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

const DealerHandCount = 10

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

const PlayerHardHandCount = 16

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

const PlayerSoftHandCount = 8

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

const PlayerPairHandCount = 10

const PlayerHandCount = PlayerHardHandCount + PlayerSoftHandCount + PlayerPairHandCount
const HandCount = PlayerHandCount * DealerHandCount

type PlayerAction string

const (
	Hit    PlayerAction = "H"
	Stay   PlayerAction = "S"
	Double PlayerAction = "D"
	Split  PlayerAction = "P"
)

const bankrollLength = 4
const mainBetLength = 4
