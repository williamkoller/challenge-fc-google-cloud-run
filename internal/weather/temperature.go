package weather

import "math"

// Temperatures holds the same reading in the three scales required by the API.
type Temperatures struct {
	Celsius    float64
	Fahrenheit float64
	Kelvin     float64
}

// FromCelsius converts a Celsius reading.
// Fahrenheit is C × 1.8 + 32. Kelvin is C + 273.
// The input is rounded to two decimal places first, then converted, so the
// three fields stay consistent with each other in the JSON payload.
func FromCelsius(celsius float64) Temperatures {
	celsius = round2(celsius)
	return Temperatures{
		Celsius:    celsius,
		Fahrenheit: round2(celsius*1.8 + 32),
		Kelvin:     round2(celsius + 273),
	}
}

func round2(value float64) float64 {
	return math.Round(value*100) / 100
}
