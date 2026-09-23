package weather

import "testing"

func TestValidateZipcode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		zipcode string
		wantErr error
	}{
		{name: "eight digits", zipcode: "01001000"},
		{name: "empty", zipcode: "", wantErr: ErrInvalidZipcode},
		{name: "too short", zipcode: "1234567", wantErr: ErrInvalidZipcode},
		{name: "too long", zipcode: "123456789", wantErr: ErrInvalidZipcode},
		{name: "letters", zipcode: "abcdefgh", wantErr: ErrInvalidZipcode},
		{name: "hyphen", zipcode: "01001-000", wantErr: ErrInvalidZipcode},
		{name: "digit and letter", zipcode: "0100100a", wantErr: ErrInvalidZipcode},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateZipcode(tt.zipcode)
			if err != tt.wantErr {
				t.Fatalf("ValidateZipcode(%q) error = %v, want %v", tt.zipcode, err, tt.wantErr)
			}
		})
	}
}

func TestLocationQuery(t *testing.T) {
	t.Parallel()

	withState := Location{City: "São Paulo", State: "SP"}.Query()
	if withState != "São Paulo, SP, Brazil" {
		t.Fatalf("query = %q", withState)
	}

	cityOnly := Location{City: "Santos"}.Query()
	if cityOnly != "Santos, Brazil" {
		t.Fatalf("query = %q", cityOnly)
	}
}
