package blackjack

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helpers

func openTextFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := fileInfo.Size()
	fileContent := make([]byte, fileSize)

	_, err = file.Read(fileContent)
	if err != nil {
		return nil, err
	}

	return fileContent, nil
}

func getTestStrategy() ([]byte, error) {
	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return nil, cwdErr
	}

	testStrategyFileName := cwd + "/__mock_data__/strategies/ideal.strategy"
	raw, readError := openTextFile(testStrategyFileName)

	if readError != nil {
		return nil, readError
	}

	// Remove line breaks and spaces
	noLineBreaks := bytes.ReplaceAll(raw, []byte("\n"), []byte(""))
	noSpaces := bytes.ReplaceAll(noLineBreaks, []byte(" "), []byte(""))

	return noSpaces, nil
}

// Tests

func TestStrategyActionParsing(t *testing.T) {
	raw, err := getTestStrategy()
	if err != nil {
		t.Fatal(err)
	}
	strategy, _ := NewStrategy(raw)

	testCases := []struct {
		desc           string
		playerHand     Hand
		dealerHand     Hand
		expectedAction PlayerAction
	}{
		{
			desc:           "Should handle hard values",
			playerHand:     Hand{NewCard(King, Clubs), NewCard(Five, Clubs)},
			dealerHand:     Hand{NewCard(Two, Clubs)},
			expectedAction: Stand,
		},
		{
			desc:           "Should handle soft values",
			playerHand:     Hand{NewCard(Ace, Clubs), NewCard(Five, Clubs)},
			dealerHand:     Hand{NewCard(Six, Clubs)},
			expectedAction: Double,
		},
		{
			desc:           "Should handle Ace upcard",
			playerHand:     Hand{NewCard(Ace, Hearts), NewCard(Five, Clubs)},
			dealerHand:     Hand{NewCard(Ace, Diamonds)},
			expectedAction: Hit,
		},
		{
			desc:           "Should handle pairs",
			playerHand:     Hand{NewCard(Four, Clubs), NewCard(Four, Hearts)},
			dealerHand:     Hand{NewCard(Four, Clubs)},
			expectedAction: Hit,
		},
		{
			desc:           "Should handle pairs",
			playerHand:     Hand{NewCard(Nine, Clubs), NewCard(Nine, Hearts)},
			dealerHand:     Hand{NewCard(Seven, Clubs)},
			expectedAction: Stand,
		},
		{
			desc:           "Should handle two aces",
			playerHand:     Hand{NewCard(Ace, Spades), NewCard(Ace, Hearts)},
			dealerHand:     Hand{NewCard(Ten, Clubs)},
			expectedAction: SplitOrHit,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			got := strategy.Play(testCase.playerHand, testCase.dealerHand)
			want := testCase.expectedAction
			assert.Equal(t, want, got)
		})
	}
}

func TestStrategyPlayMapping(t *testing.T) {
	raw, err := getTestStrategy()
	if err != nil {
		t.Fatal(err)
	}
	strategy, _ := NewStrategy(raw)

	t.Run("Should return Stand if mapping does not exist (Bust)", func(t *testing.T) {
		playerHand := Hand{NewCard(King, Spades), NewCard(King, Hearts), NewCard(King, Diamonds)}
		dealerHand := Hand{NewCard(Two, Clubs)}

		got := strategy.Play(playerHand, dealerHand)
		want := Stand
		assert.Equal(t, want, got)
	})

	t.Run("Should return Stand if mapping does not exist (Blackjack)", func(t *testing.T) {
		playerHand := Hand{NewCard(King, Spades), NewCard(Ace, Hearts)}
		dealerHand := Hand{NewCard(Two, Clubs)}

		got := strategy.Play(playerHand, dealerHand)
		want := Stand
		assert.Equal(t, want, got)
	})
}

func TestPlayerInitializationParsing(t *testing.T) {
	raw, err := getTestStrategy()
	if err != nil {
		t.Fatal(err)
	}
	strategy, _ := NewStrategy(raw)

	t.Run("Should parse initial bankroll", func(t *testing.T) {
		got := strategy.GetInitialBankroll()
		want := 1000.0
		assert.Equal(t, want, got)
	})
}

func TestStrategyBetParsing(t *testing.T) {
	raw, err := getTestStrategy()
	if err != nil {
		t.Fatal(err)
	}
	strategy, _ := NewStrategy(raw)

	t.Run("Should handle main game bet", func(t *testing.T) {
		got := strategy.Bet( /* shoe? */ )
		want := 1
		assert.Equal(t, want, got)
	})
}

func TestStrategyParsingErrors(t *testing.T) {
	t.Run("Should handle invalid strategy length", func(t *testing.T) {
		raw := bytes.Repeat([]byte("H"), 1)
		_, err := NewStrategy(raw)
		assert.Error(t, err)
	})
}

func TestStrategyReturnEncodedStrategy(t *testing.T) {
	raw, err := getTestStrategy()
	if err != nil {
		t.Fatal(err)
	}
	strategy, _ := NewStrategy(raw)

	t.Run("Should return encoded strategy", func(t *testing.T) {
		got := strategy.GetEncodedStrategy()
		want := raw
		assert.Equal(t, want, got)
	})
}

func TestStrategySequencing(t *testing.T) {
	t.Run("Should the encoding sequence for strategies", func(t *testing.T) {
		sequence := GetSequencing()
		assert.Len(t, sequence, 348)
	})
}
