package main

import "testing"

func TestCleanInput(t *testing.T) {
  cases := []struct {
    input    string
    expected []string
  }{
    {
      input:    " hello world ",
      expected: []string{"hello", "world"},
    },
    {
      input:    "",
      expected: []string{},
    },
    // Need more cases
  }

  for _, c := range cases {
    actual := cleanInput(c.input)

    // Check the length of the actual slice
    if len(actual) != len(c.expected) {
      t.Errorf("No len match yo")
    }

    for i := range actual {
      word     := actual[i]
      expected := c.expected[i]
      if word != expected {
        t.Errorf("Wtf is this shit: (%s) I wanted this shit: (%s)", word, expected)
      }
    }
  }
}
