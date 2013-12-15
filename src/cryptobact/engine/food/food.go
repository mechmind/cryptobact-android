package food

type Food struct {
	X        float64
	Y        float64
	Calories float64
	Eaten    bool
}

func (f *Food) Locate() (float64, float64) {
	return f.X, f.Y
}
