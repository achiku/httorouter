package tiniyrouter

import "testing"

func TestFind(t *testing.T) {
	cases := []struct {
		Src    string
		Target byte
		Pos    int
	}{
		{Src: "/test/api", Target: '/', Pos: 0},
		{Src: "/test/api/*", Target: '*', Pos: 10},
		{Src: "/test/apiv1/:name", Target: ':', Pos: 12},
	}

	for _, c := range cases {
		pos := find(c.Src, c.Target, 0, len(c.Src))
		if pos != c.Pos {
			t.Errorf("want %d got %d", c.Pos, pos)
		}
	}
}

func TestGetIndexPosition(t *testing.T) {
	node := treeNode{
		indices: "/hello/world",
	}
	str := "w"
	i := node.getIndexPosition(byte(str))
	t.Logf("%d", i)
}
