package buffer

import "testing"

func TestAddEnd(t *testing.T) {
	line := Line{}
	line.Init(nil, nil)
	i := 0
	for c := 'a'; c <= 'z'; c++ {
		line.AddAt(i, c)
		i++
	}
	if len(line.Data) != 26 {
		t.Error("not all characters were added")
	}
	c := 'a'
	for _, r := range line.Data {
		if r != c {
			t.Error("incorrect character order")
		}
		c++
	}
}

func TestAddStart(t *testing.T) {
	line := Line{}
	line.Init(nil, nil)
	for c := 'a'; c <= 'z'; c++ {
		line.AddAt(0, c)
	}
	if len(line.Data) != 26 {
		t.Error("not all characters were added")
	}
	c := 'z'
	for _, r := range line.Data {
		if r != c {
			t.Error("incorrect character order")
		}
		c--
	}
}

func TestRemoveStart(t *testing.T) {
	line := Line{}
	line.Init(nil, nil)
	line.AddAt(0, 'a')
	line.AddAt(1, 'b')
	line.AddAt(2, 'c')
	if len(line.Data) != 3 {
		t.Error("bad size")
	}
	line.RemoveAt(0)
	if len(line.Data) != 2 {
		t.Error("bad size")
	}
	if line.Data[0] != 'b' {
		t.Error("bad first rune")
	}
	if line.Data[1] != 'c' {
		t.Error("bad second rune")
	}
}

func TestRemoveMiddle(t *testing.T) {
	line := Line{}
	line.Init(nil, nil)
	line.AddAt(0, 'a')
	line.AddAt(1, 'b')
	line.AddAt(2, 'c')
	if len(line.Data) != 3 {
		t.Error("bad size")
	}
	line.RemoveAt(1)
	if len(line.Data) != 2 {
		t.Error("bad size")
	}
	if line.Data[0] != 'a' {
		t.Error("bad first rune")
	}
	if line.Data[1] != 'c' {
		t.Error("bad second rune")
	}
}

func TestRemoveEnd(t *testing.T) {
	line := Line{}
	line.Init(nil, nil)
	line.AddAt(0, 'a')
	line.AddAt(1, 'b')
	line.AddAt(2, 'c')
	if len(line.Data) != 3 {
		t.Error("bad size")
	}
	line.RemoveAt(2)
	if len(line.Data) != 2 {
		t.Error("bad size")
	}
	if line.Data[0] != 'a' {
		t.Error("bad first rune")
	}
	if line.Data[1] != 'b' {
		t.Error("bad second rune")
	}
}
