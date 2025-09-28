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
			input:    "foo",
			expected: []string{"foo"},
		},
		{
			input:    "   ",
			expected: []string{},
		},
		{
			input:    "bar   baz qux",
			expected: []string{"bar", "baz", "qux"},
		},
		{
			input:    "\tfoo\tbar\nbaz  ",
			expected: []string{"foo", "bar", "baz"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("Failed for input: %s\n\tExpected: %v -> len = %d\n\t Got: %v  -> len = %d", c.input, c.expected, len(c.expected), actual, len(actual))
			}
		}
	}
}
