package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "hello  world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  HeLLo  wORLd  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "hello  world",
			expected: []string{"hello", "world"},
		},
	}

	for _, c := range cases {
		response := CleanInput(c.input)

		if len(response) != len(c.expected) {
			t.Errorf("Expected slice with size: %d, received: %d", len(c.expected), len(response))
			return
		}

		for i := range response {
			receivedWord := response[i]
			expectedWord := c.expected[i]

			if receivedWord != expectedWord {
				t.Errorf("Expected: %s, received: %s", expectedWord, receivedWord)
			}
		}
	}
}
