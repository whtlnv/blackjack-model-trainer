package main

type PlayerStatistics struct {
	Strategy []byte

	GamesSeen   int
	GamesPlayed int
	GamesWon    int
	GamesLost   int
	GamesPushed int

	InitialBankroll float64
	Bankroll        float64
	BankrollDelta   float64
}

func BlackjackFitnessFunction(statistics PlayerStatistics) float64 {
	fitness := 0.0

	if statistics.Bankroll > 0 {
		fitness += statistics.BankrollDelta * statistics.BankrollDelta

		// if player made money, reward doing it in fewer games
		gamesNotPlayed := statistics.GamesSeen - statistics.GamesPlayed
		fitness += float64(gamesNotPlayed)
	} else {
		// if player lost money, reward doing so in more games
		fitness += float64(statistics.GamesPlayed)
	}

	winrate := float64(statistics.GamesWon) / float64(statistics.GamesPlayed)
	fitness += winrate * 100

	return fitness
}
