package genetics

type Randomizerish interface {
	EventDidHappen(probability float64) bool
	PickOne(options []byte) byte
}

// TODO
