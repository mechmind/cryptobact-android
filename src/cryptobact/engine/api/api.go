package api

func GetBacts(device_id int) []Bact {
	return []Bact{
		Bact{
			1,
			10.0,
			10.0,
			1.0,
			0.5,
			0.5,
			0.5,
			0.5,
			0.5,
			1000,
			1000,
		},
		Bact{
			2,
			15.0,
			15.0,
			1.0,
			0.5,
			0.5,
			0.5,
			0.5,
			0.5,
			1000,
			1000,
		},
	}
}
