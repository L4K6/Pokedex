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
			input:    "hello    world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "summer breeze !!! hahAha",
			expected: []string{"summer", "breeze", "!!!", "hahaha"},
		},
		{
			input:    "ta Ble, what ??",
			expected: []string{"ta", "ble,", "what", "??"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("test failed: mismatched length of %s and %s", actual, c.expected)
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("test has failed: %s != %s", word, expectedWord)
			}
		}
	}
}
