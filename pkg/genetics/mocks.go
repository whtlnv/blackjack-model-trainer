package genetics

import (
	"github.com/stretchr/testify/mock"
)

type RandomizerMock struct {
	mock.Mock
}

func (r *RandomizerMock) EventDidHappen(probability float64) bool {
	args := r.Called(probability)
	return args.Bool(0)
}

func (r *RandomizerMock) PickOne(options []byte) byte {
	args := r.Called(options)
	return args.Get(0).(byte)
}

func (r *RandomizerMock) NumberBetween(min, max int) int {
	args := r.Called(min, max)
	return args.Int(0)
}

type ChromosomeMock struct {
	mock.Mock
}

func (c *ChromosomeMock) Raw() []byte {
	args := c.Called()
	return args.Get(0).([]byte)
}

func (c *ChromosomeMock) Merge(other *Chromosome, randomizer Randomizerish) *Chromosome {
	args := c.Called(other, randomizer)
	return args.Get(0).(*Chromosome)
}

func (c *ChromosomeMock) Mutate(randomizer Randomizerish) *Chromosome {
	args := c.Called(randomizer)
	return args.Get(0).(*Chromosome)
}
