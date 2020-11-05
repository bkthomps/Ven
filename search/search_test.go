package search

import (
	"testing"

	"github.com/bkthomps/Ven/buffer"
)

func TestSingleLineNoMatches(t *testing.T) {
	line := &buffer.Line{}
	i := 0
	for c := 'a'; c <= 'z'; c++ {
		line.AddAt(i, c)
		i++
	}
	pattern := Compile([]rune("zyx"))
	match := pattern.Search(line)
	if match != nil {
		t.Error("expected no matches")
	}
}

func TestSingleLineSingleMatch(t *testing.T) {
	line := &buffer.Line{}
	i := 0
	for c := 'a'; c <= 'z'; c++ {
		line.AddAt(i, c)
		i++
	}
	pattern := Compile([]rune("cde"))
	match := pattern.Search(line)
	if match == nil {
		t.Error("bad match count")
	} else {
		if match.Line != line {
			t.Error("bad match line")
		}
		if match.StartOffset != 2 {
			t.Error("bad match offset")
		}
	}
}
