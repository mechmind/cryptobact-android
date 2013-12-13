package evo

type Environment struct {
	MutateProbability   float64
	MutateRate          float64
	RecombinationChance float64
	RecombinationDrop   float64
}

var DefaultEnvironment = &Environment{
	MutateProbability:   0.5,
	MutateRate:          1.0,
	RecombinationChance: 1.0,
	RecombinationDrop:   10,
}
