package tiniyrouter

import "testing"

func TestMin(t *testing.T) {
	cases := []struct {
		A        int
		B        int
		Expected int
	}{
		{A: 1, B: 2, Expected: 1},
		{A: -1, B: 2, Expected: -1},
		{A: 1, B: -2, Expected: -2},
		{A: 1, B: 1, Expected: 1},
	}

	for _, c := range cases {
		v := min(c.A, c.B)
		if v != c.Expected {
			t.Errorf("want %d got %d", c.Expected, v)
		}
	}
}

func TestCountParams(t *testing.T) {
	cases := []struct {
		Path  string
		Count uint8
	}{
		{Path: "/users", Count: 0},
		{Path: "/users/:firstName", Count: 1},
		{Path: "/users/:firstName/:lastName", Count: 2},
		{Path: "/users/*", Count: 1},
	}

	for _, c := range cases {
		n := countParams(c.Path)
		if n != c.Count {
			t.Errorf("want %d got %d", c.Count, n)
		}
	}
}

func TestIncrementChildPrio(t *testing.T) {
}
