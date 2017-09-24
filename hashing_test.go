package tinyURL

import "testing"

func TestHashing(t *testing.T) {
	cases := []struct {
		in string
	}{
		{""},
		{"hello"},
		{"123456"},
		{"kdkfjdiej29j29212312312ref"},
	}
	for _, c := range cases {
		got := HashOut(HashIn(c.in))
		if got != c.in {
			t.Errorf("Hashed %q but got %q", c.in, got)
		}
	}
}