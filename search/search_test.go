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
	matches, _ := AllMatches("zyx", line, 40)
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
	matches, _ := AllMatches("cde", line, 40)
	if len(matches) != 1 {
		t.Error("bad match count")
	}
	for _, match := range matches {
		if match.Line != line {
			t.Error("bad match line")
		}
		for _, instance := range match.Instances {
			if instance.StartOffset != 2 {
				t.Error("bad match offset")
			}
			if instance.Length != 3 {
				t.Error("bad match length")
			}
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
	matches, _ := AllMatches("cde", line, 40)
	if len(matches) != 1 {
		t.Error("bad match count")
	}
	charset := 26
	for _, match := range matches {
		if len(match.Instances) != repetitions {
			t.Error("bad match count")
		}
		if match.Line != line {
			t.Error("bad match line")
		}
		for i, instance := range match.Instances {
			if instance.StartOffset != charset*i+2 {
				t.Error("bad match offset")
			}
			if instance.Length != 3 {
				t.Error("bad match length")
			}
		}
	}
}

func TestMultipleLinesMultipleMatches(t *testing.T) {
	file := buffer.File{}
	file.Init("TestMultipleLinesMultipleMatches.txt")
	repetitions := 3
	for i := 0; i < repetitions; i++ {
		for c := 'a'; c <= 'z'; c++ {
			file.Add(c)
		}
		file.Add('\n')
	}
	matches, _ := AllMatches("cde", file.First, 40)
	if len(matches) != repetitions {
		t.Error("bad match count")
	}
	line := file.First
	for _, match := range matches {
		if len(match.Instances) != 1 {
			t.Error("bad match count")
		}
		if match.Line != line {
			t.Error("bad match line")
		}
		for _, instance := range match.Instances {
			if instance.StartOffset != 2 {
				t.Error("bad match offset")
			}
			if instance.Length != 3 {
				t.Error("bad match length")
			}
		}
		line = line.Next
	}
}

func TestMultipleLinesMultipleMatchesCropped(t *testing.T) {
	file := buffer.File{}
	file.Init("TestMultipleLinesMultipleMatchesCropped.txt")
	for i := 0; i < 3; i++ {
		for c := 'a'; c <= 'z'; c++ {
			file.Add(c)
		}
		file.Add('\n')
	}
	matches, _ := AllMatches("cde", file.First, 2)
	if len(matches) != 2 {
		t.Error("bad match count")
	}
	line := file.First
	for _, match := range matches {
		if len(match.Instances) != 1 {
			t.Error("bad match count")
		}
		if match.Line != line {
			t.Error("bad match line")
		}
		for _, instance := range match.Instances {
			if instance.StartOffset != 2 {
				t.Error("bad match offset")
			}
			if instance.Length != 3 {
				t.Error("bad match length")
			}
		}
		line = line.Next
	}
}

func TestRegexDot(t *testing.T) {
	line := &buffer.Line{}
	i := 0
	for c := 'a'; c <= 'z'; c++ {
		line.AddAt(i, c)
		i++
	}
	matches, _ := AllMatches("c.e", line, 40)
	if len(matches) != 1 {
		t.Error("bad match count")
	}
	for _, match := range matches {
		if match.Line != line {
			t.Error("bad match line")
		}
		for _, instance := range match.Instances {
			if instance.StartOffset != 2 {
				t.Error("bad match offset")
			}
			if instance.Length != 3 {
				t.Error("bad match length")
			}
		}
	}
}

func TestRegexStar(t *testing.T) {
	line := &buffer.Line{}
	i := 0
	for c := 'a'; c <= 'z'; c++ {
		line.AddAt(i, c)
		i++
	}
	matches, _ := AllMatches("a.*z", line, 40)
	if len(matches) != 1 {
		t.Error("bad match count")
	}
	for _, match := range matches {
		if match.Line != line {
			t.Error("bad match line")
		}
		for _, instance := range match.Instances {
			if instance.StartOffset != 0 {
				t.Error("bad match offset")
			}
			if instance.Length != 26 {
				t.Error("bad match length")
			}
		}
	}
}

func TestRegexPipe(t *testing.T) {
	line := &buffer.Line{}
	i := 0
	for c := 'a'; c <= 'z'; c++ {
		line.AddAt(i, c)
		i++
	}
	matches, _ := AllMatches("(c.e)|(f.h)", line, 40)
	if len(matches) != 1 {
		t.Error("bad match count")
	}
	for _, match := range matches {
		if len(match.Instances) != 2 {
			t.Error("bad match count")
		}
		if match.Line != line {
			t.Error("bad match line")
		}
		for i, instance := range match.Instances {
			if instance.StartOffset != 2+i*3 {
				t.Error("bad match offset")
			}
			if instance.Length != 3 {
				t.Error("bad match length")
			}
		}
	}
}

func TestRegexBrackets(t *testing.T) {
	line := &buffer.Line{}
	i := 0
	for c := 'a'; c <= 'z'; c++ {
		line.AddAt(i, c)
		i++
	}
	matches, _ := AllMatches("[a-d]", line, 40)
	if len(matches) != 1 {
		t.Error("bad match count")
	}
	for _, match := range matches {
		if len(match.Instances) != 4 {
			t.Error("bad match count")
		}
		if match.Line != line {
			t.Error("bad match line")
		}
		for i, instance := range match.Instances {
			if instance.StartOffset != i {
				t.Error("bad match offset")
			}
			if instance.Length != 1 {
				t.Error("bad match length")
			}
		}
	}
}

func TestUnicode(t *testing.T) {
	line := &buffer.Line{}
	i := 0
	repetitions := 3
	for j := 0; j < repetitions; j++ {
		line.AddAt(i, '象')
		i++
		line.AddAt(i, '形')
		i++
		line.AddAt(i, '字')
		i++
		line.AddAt(i, '㫃')
		i++
		line.AddAt(i, '池')
		i++
	}
	matches, _ := AllMatches("形字", line, 40)
	if len(matches) != 1 {
		t.Error("bad match count")
	}
	for _, match := range matches {
		if len(match.Instances) != repetitions {
			t.Error("bad match count")
		}
		if match.Line != line {
			t.Error("bad match line")
		}
		for i, instance := range match.Instances {
			if instance.StartOffset != 5*i+1 {
				t.Errorf("bad match offset: %d", instance.StartOffset)
			}
			if instance.Length != 2 {
				t.Errorf("bad match length: %d", instance.Length)
			}
		}
	}
}
