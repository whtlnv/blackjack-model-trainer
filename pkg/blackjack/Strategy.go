package blackjack

import (
	"fmt"
	"strconv"
)

type Strategyish interface {
	Play(playerHand Hand, dealerHand Hand) PlayerAction
	Bet() int
	GetInitialBankroll() float64
	GetEncodedStrategy() []byte
}

type Strategy struct {
	pairMap         map[PlayerHand]map[DealerHand]PlayerAction
	softMap         map[PlayerHand]map[DealerHand]PlayerAction
	hardMap         map[PlayerHand]map[DealerHand]PlayerAction
	mainGameBetMap  map[int]int
	initialBankroll int

	raw []byte
}

// Factory

func NewStrategy(raw []byte) (*Strategy, error) {
	strategy := &Strategy{}

	strategy.raw = raw

	err := validateRawStrategy(raw)
	if err != nil {
		return strategy, err
	}

	err = strategy.parseActionMap(raw)
	if err != nil {
		return strategy, err
	}

	err = strategy.parseBankroll(raw)
	if err != nil {
		return strategy, err
	}

	err = strategy.parseMainBettingStrategy(raw)
	if err != nil {
		return strategy, err
	}

	return strategy, nil
}

// Static methods

func GetSequencing() [][]byte {
	newArrayFilledWith := func(length int, value []byte) [][]byte {
		array := make([][]byte, length)
		for i := range array {
			array[i] = value
		}
		return array
	}

	actions := []byte("HSDP")
	numbers := []byte("0123456789ABCDF")

	hardMapLength := DealerHandCount * PlayerHardHandCount
	softMapLength := DealerHandCount * PlayerSoftHandCount
	pairMapLength := DealerHandCount * PlayerPairHandCount
	bankrollLength := bankrollLength
	bettingLength := mainBetLength

	sequence := make([][]byte, 0)
	sequence = append(sequence, newArrayFilledWith(hardMapLength, actions)...)
	sequence = append(sequence, newArrayFilledWith(softMapLength, actions)...)
	sequence = append(sequence, newArrayFilledWith(pairMapLength, actions)...)
	sequence = append(sequence, newArrayFilledWith(bankrollLength, numbers)...)
	sequence = append(sequence, newArrayFilledWith(bettingLength, numbers)...)

	return sequence
}

// Public methods

func (strategy *Strategy) GetInitialBankroll() float64 {
	return float64(strategy.initialBankroll)
}

func (strategy *Strategy) Play(playerHand Hand, dealerHand Hand) PlayerAction {
	dealerScore, _ := dealerHand.Score()
	playerScore, _ := playerHand.Score()
	dealerHighScore := DealerHand(dealerScore.High)
	playerHighScore := PlayerHand(playerScore.High)
	playerLowScore := PlayerHand(playerScore.Low)

	var action PlayerAction
	if playerHand.IsPair() {
		action = strategy.pairMap[playerLowScore][dealerHighScore]
	} else if playerHand.HasSoftValue() {
		action = strategy.softMap[playerHighScore][dealerHighScore]
	} else {
		action = strategy.hardMap[playerHighScore][dealerHighScore]
	}

	if action == "" {
		action = Stand
	}

	return action
}

func (strategy *Strategy) Bet() int {
	return strategy.mainGameBetMap[0]
}

func (strategy *Strategy) GetEncodedStrategy() []byte {
	return strategy.raw
}

// Private methods

func (strategy *Strategy) parseActionMap(raw []byte) error {
	rawHardMapStartsAt := 0
	rawSoftMapStartsAt := rawHardMapStartsAt + (DealerHandCount * PlayerHardHandCount)
	rawPairMapStartsAt := rawSoftMapStartsAt + (DealerHandCount * PlayerSoftHandCount)

	rawHardMap := raw[rawHardMapStartsAt:rawSoftMapStartsAt]
	rawSoftMap := raw[rawSoftMapStartsAt:rawPairMapStartsAt]
	rawPairMap := raw[rawPairMapStartsAt:]

	strategy.hardMap = strategy.parseHardMap(rawHardMap)
	strategy.softMap = strategy.parseSoftMap(rawSoftMap)
	strategy.pairMap = strategy.parsePairMap(rawPairMap)

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

func (strategy *Strategy) parseBankroll(raw []byte) error {
	rawHardMapStartsAt := 0
	rawBankrollStartsAt := rawHardMapStartsAt + (DealerHandCount * PlayerHardHandCount) + (DealerHandCount * PlayerSoftHandCount) + (DealerHandCount * PlayerPairHandCount)

	rawBankrollHex := raw[rawBankrollStartsAt : rawBankrollStartsAt+bankrollLength]
	parsed, parseErr := strconv.ParseInt(string(rawBankrollHex), 16, 32)

	if parseErr != nil {
		return parseErr
	}

	strategy.initialBankroll = int(parsed)

	return nil
}

func (strategy *Strategy) parseMainBettingStrategy(raw []byte) error {
	parsedMap := make(map[int]int)

	rawHardMapStartsAt := 0
	rawBettingStartsAt :=
		rawHardMapStartsAt +
			(DealerHandCount * PlayerHardHandCount) +
			(DealerHandCount * PlayerSoftHandCount) +
			(DealerHandCount * PlayerPairHandCount) +
			bankrollLength

	rawBetHex := raw[rawBettingStartsAt : rawBettingStartsAt+mainBetLength]
	parsed, parseErr := strconv.ParseInt(string(rawBetHex), 16, 32)

	if parseErr != nil {
		return parseErr
	}

	parsedMap[0] = int(parsed)
	strategy.mainGameBetMap = parsedMap

	return nil
}

// Helper methods

func validateRawStrategy(raw []byte) error {
	expectedLength := HandCount + bankrollLength + mainBetLength
	if len(raw) != expectedLength {
		return fmt.Errorf("expected strategy length to be %d, got %d", expectedLength, len(raw))
	}

	return nil
}

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
