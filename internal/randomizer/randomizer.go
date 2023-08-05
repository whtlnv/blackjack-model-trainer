package randomizer

import (
	"math/rand"
)

type Randomizer struct {
	generator *rand.Rand
}

func NewRandomizer(seed int64) *Randomizer {
	source := rand.NewSource(seed)
	generator := rand.New(source)
	return &Randomizer{generator}
}

func (r *Randomizer) EventDidHappen(probability float64) bool {
	return r.generator.Float64() <= probability
}

func (r *Randomizer) PickOne(options []byte) byte {
	return options[r.generator.Intn(len(options))]
}

func (r *Randomizer) NumberBetween(min, max int) int {
	return r.generator.Intn(max-min) + min
}
