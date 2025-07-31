package main

import (
	"testing"
)

func TestCleanInput( t *testing.T) {
	cases := []struct{
		input string
		expected []string
	}{
		{
		input: "hello world",
		expected: []string{"hello", "world"},
		},
		{
		input: "zoro luffy",
		expected: []string{"zoro", "luffy"},
		},
	}
	

	for _,c := range cases {

		actual := cleanInput(c.input)
		expected := c.expected

		if len(actual) != len(expected) {
			t.Errorf("The length of the strings do not match")
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			
			if word != expectedWord {
				t.Errorf("For input %q, expected %q but got %q", c.input, expectedWord, word)
			}
		}
	}
}