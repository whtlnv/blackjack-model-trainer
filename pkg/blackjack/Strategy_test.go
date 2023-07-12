package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	if len(noSpaces) != 340 {
		return nil, fmt.Errorf("Expected strategy length to be 340, got %d", len(noSpaces))
	}

	return noSpaces, nil
}

func TestStrategyParsing(t *testing.T) {
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
			playerHand:     Hand{Card{King, Clubs}, Card{Five, Clubs}},
			dealerHand:     Hand{Card{Two, Clubs}},
			expectedAction: Stay,
		},
		{
			desc:           "Should handle soft values",
			playerHand:     Hand{Card{Ace, Clubs}, Card{Five, Clubs}},
			dealerHand:     Hand{Card{Six, Clubs}},
			expectedAction: Double,
		},
		{
			desc:           "Should handle pairs",
			playerHand:     Hand{Card{Four, Clubs}, Card{Four, Hearts}},
			dealerHand:     Hand{Card{Four, Clubs}},
			expectedAction: Hit,
		},
		{
			desc:           "Should handle pairs",
			playerHand:     Hand{Card{Nine, Clubs}, Card{Nine, Hearts}},
			dealerHand:     Hand{Card{Seven, Clubs}},
			expectedAction: Stay,
		},
		{
			desc:           "Should handle two aces",
			playerHand:     Hand{Card{Ace, Spades}, Card{Ace, Hearts}},
			dealerHand:     Hand{Card{Two, Clubs}},
			expectedAction: Split,
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

func TestStrategyParsingErrors(t *testing.T) {
	t.Run("Should handle invalid strategy length", func(t *testing.T) {
		raw := bytes.Repeat([]byte("H"), 1)
		_, err := NewStrategy(raw)
		assert.Error(t, err)
	})
}
