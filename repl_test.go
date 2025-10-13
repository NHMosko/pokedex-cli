package main

import "testing"

func TestCleanInput(t *testing.T) {
	 cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  pokeMOn is, Bulbasaur, in this              WORLD  ",
			expected: []string{"pokemon", "is", "bulbasaur", "in", "this", "world"},
		},
		{
			input: " 333343slite O,F  hand ,",
			expected: []string{"333343slite", "o", "f", "hand"},
		},	
		{
			input: " nicolas  fireType   fire,water",
			expected: []string{"nicolas", "firetype", "fire", "water"},
		},
	}  

	for _, c := range cases {
	actual := cleanInput(c.input)
	// Check the length of the actual slice against the expected slice
	// if they don't match, use t.Errorf to print an error message
	// and fail the test

	if len(actual) != len(c.expected) {
		t.Errorf("Expected message differs in length from actual slice.\n-%v- len: %v != -%v- len: %v", c.expected, len(c.expected), actual, len(actual))
	}

	for i := range actual {
		word := actual[i]
		expectedWord := c.expected[i]
		// Check each word in the slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test

		if word != expectedWord {
			t.Errorf("Expected word doesn't match actual word.\n-%v- is not the expected -%v-", word, expectedWord)
		}
	}
}
}
