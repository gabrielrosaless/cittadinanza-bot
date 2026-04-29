package detector

import "testing"

func TestHasKeyword(t *testing.T) {
	keywords := []string{"turnos", "cittadinanza", "prenotazioni", "apertura", "ciudadania"}

	tests := []struct {
		name     string
		text     string
		expected bool
	}{
		{"Exact match", "Habilitación nuevos turnos", true},
		{"Case diff", "Nuevos TURNOS disponibles", true},
		{"Accent in text", "nueva ciudadanía", true}, // Target kw is "ciudadania", text has accent
		{"Accent in both but diff case", "CIUDADANÍA", true},
		{"No match", "cambio de horario consular", false},
		{"Partial word match", "aperturas", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasKeyword(tt.text, keywords)
			if got != tt.expected {
				t.Errorf("HasKeyword(%q) = %v; want %v", tt.text, got, tt.expected)
			}
		})
	}
}
