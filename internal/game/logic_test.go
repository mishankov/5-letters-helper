package game

import "testing"

func TestGetWordStatus(t *testing.T) {
	tests := []struct {
		guess  string
		target string
		status string
	}{
		{
			guess:  "транш",
			target: "транш",
			status: "22222",
		},
		{
			guess:  "кроат",
			target: "транш",
			status: "02011",
		},
		{
			guess:  "белок",
			target: "транш",
			status: "00000",
		},
	}

	for _, test := range tests {
		result := GetWordStatus(test.guess, test.target)

		if result != test.status {
			t.Fatalf("Guess: %q. Target: %q. Want: %q. Got: %v", test.guess, test.target, test.status, result)
		}
	}
}
