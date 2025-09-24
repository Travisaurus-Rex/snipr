package shortener

import "testing"

func TestGenerateCode(t *testing.T) {
	length := 8
	code := GenerateCode()

	if len(code) != length {
		t.Errorf("expected code length %d, but got %d", length, len(code))
	}

	for _, c := range code {
		if !contains(charset, c) {
			t.Errorf("unexpected char: %c", c)
		}
	}
}

func contains(s string, c rune) bool {
	for _, sc := range s {
		if sc == c {
			return true
		}
	}

	return false
}
