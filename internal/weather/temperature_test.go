package weather

import "testing"

func TestFromCelsius(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		celsius    float64
		fahrenheit float64
		kelvin     float64
	}{
		{name: "sample reading", celsius: 28.5, fahrenheit: 83.3, kelvin: 301.5},
		{name: "freezing point", celsius: 0, fahrenheit: 32, kelvin: 273},
		{name: "minus forty", celsius: -40, fahrenheit: -40, kelvin: 233},
		{name: "rounds celsius before converting", celsius: 10.333, fahrenheit: 50.59, kelvin: 283.33},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := FromCelsius(tt.celsius)
			if got.Celsius != round2(tt.celsius) || got.Fahrenheit != tt.fahrenheit || got.Kelvin != tt.kelvin {
				t.Fatalf("FromCelsius(%v) = %+v, want C=%v F=%v K=%v", tt.celsius, got, round2(tt.celsius), tt.fahrenheit, tt.kelvin)
			}
		})
	}
}

func TestFromCelsiusFormulas(t *testing.T) {
	t.Parallel()

	got := FromCelsius(28.5)
	if got.Fahrenheit != round2(28.5*1.8+32) {
		t.Fatalf("fahrenheit = %v, want formula result %v", got.Fahrenheit, round2(28.5*1.8+32))
	}
	if got.Kelvin != round2(28.5+273) {
		t.Fatalf("kelvin = %v, want formula result %v", got.Kelvin, round2(28.5+273))
	}
}
