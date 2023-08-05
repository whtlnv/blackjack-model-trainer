package main

type PlayerStatistics struct {
	GamesPlayed int
	GamesWon    int
	GamesLost   int
	GamesPushed int

	Bankroll      float64
	BankrollDelta float64
}

type Playerish interface {
	GetStatistics() PlayerStatistics
}

func BlackjackFitnessFunction(player Playerish) float64 {
	statistics := player.GetStatistics()

	fitness := 0.0

	fitness += statistics.Bankroll * statistics.Bankroll

	winrate := float64(statistics.GamesWon) / float64(statistics.GamesPlayed)
	fitness += winrate * 100

	return fitness
}
