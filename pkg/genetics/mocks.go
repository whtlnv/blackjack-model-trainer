package genetics

import (
	"github.com/stretchr/testify/mock"
)

type randomizerMock struct {
	mock.Mock
}

func (r *randomizerMock) EventDidHappen(probability float64) bool {
	args := r.Called(probability)
	return args.Bool(0)
}

func (r *randomizerMock) PickOne(options []byte) byte {
	args := r.Called(options)
	return args.Get(0).(byte)
}
