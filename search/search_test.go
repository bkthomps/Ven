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
	matches := AllMatches("zyx", line)
	if len(matches) != 0 {
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
	matches := AllMatches("cde", line)
	if len(matches) != 1 {
		t.Error("bad match count")
	}
	for _, match := range matches {
		if match.Line != line {
			t.Error("bad match line")
		}
		if match.StartOffset != 2 {
			t.Error("bad match offset")
		}
		if match.Length != 3 {
			t.Error("bad match length")
		}
	}
}

func TestSingleLineMultipleMatches(t *testing.T) {
	line := &buffer.Line{}
	repetitions := 3
	i := 0
	for j := 0; j < repetitions; j++ {
		for c := 'a'; c <= 'z'; c++ {
			line.AddAt(i, c)
			i++
		}
	}
	matches := AllMatches("cde", line)
	if len(matches) != repetitions {
		t.Error("bad match count")
	}
	charset := 26
	for j, match := range matches {
		if match.Line != line {
			t.Error("bad match line")
		}
		if match.StartOffset != charset*j+2 {
			t.Error("bad match offset")
		}
		if match.Length != 3 {
			t.Error("bad match length")
		}
	}
}
