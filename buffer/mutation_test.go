package buffer

import (
	"testing"
)

func TestAddCharacters(t *testing.T) {
	file := File{}
	file.Init("")
	runes := []rune{'a', '\t', 'b', 'c', ' ', ' ', '\t', '汉', '字', '🦀', '😂'}
	if !file.CanSafeQuit() {
		t.Error("should be able to safe quit")
	}
	for i, r := range runes {
		_, addedLine := file.Add(r)
		if addedLine {
			t.Error("line should not have been added")
		}
		if file.runeOffset != i+1 {
			t.Error("bad offset")
		}
	}
	for i, r := range file.First.Data {
		if r != runes[i] {
			t.Error("did not add runes correctly")
		}
	}
	if file.Lines != 1 {
		t.Error("expected one line")
	}
	if file.First.Next != nil {
		t.Error("did not set nil correctly")
	}
	if file.CanSafeQuit() {
		t.Error("should not be able to safe quit")
	}
}

func TestAddNewlines(t *testing.T) {
	file := File{}
	file.Init("")
	runes := []rune{'a', '\t', 'b', 'c', ' ', ' ', '\t', '汉', '字', '🦀', '😂'}
	if !file.CanSafeQuit() {
		t.Error("should be able to safe quit")
	}
	for i, r := range runes {
		_, addedLine := file.Add(r)
		if addedLine {
			t.Error("line should not have been added")
		}
		if file.runeOffset != i+1 {
			t.Error("bad offset")
		}
	}
	_, addedLine := file.Add('\n')
	if !addedLine {
		t.Error("line should have been added")
	}
	if file.runeOffset != 0 {
		t.Error("offset should be reset")
	}
	for i, r := range runes {
		_, addedLine := file.Add(r)
		if addedLine {
			t.Error("line should not have been added")
		}
		if file.runeOffset != i+1 {
			t.Error("bad offset")
		}
	}
	for i, r := range file.First.Data {
		if r != runes[i] {
			t.Error("did not add runes correctly")
		}
	}
	for i, r := range file.First.Next.Data {
		if r != runes[i] {
			t.Error("did not add runes correctly")
		}
	}
	if file.Lines != 2 {
		t.Error("expected two lines")
	}
	if file.First.Next.Next != nil {
		t.Error("did not set nil correctly")
	}
	if file.CanSafeQuit() {
		t.Error("should not be able to safe quit")
	}
}

func TestRemoveImpossible(t *testing.T) {
	file := File{}
	file.Init("")
	x := file.Remove()
	if x != 0 {
		t.Error("invalid no-op remove")
	}
}

func TestRemove(t *testing.T) {
	file := File{}
	file.Init("")
	for c := 'a'; c <= 'z'; c++ {
		file.Add(c)
	}
	file.Left()
	size := 26 - 1
	for c := 'z'; c > 'a'; c-- {
		size--
		x := file.Remove()
		if x != size {
			t.Error("bad offset")
		}
		for i := 0; i < size; i++ {
			u := rune('a' + i)
			if u != file.First.Data[i] {
				t.Error("bad removal")
			}
		}
	}
	x := file.Remove()
	if x != 0 {
		t.Error("should be zero index")
	}
	if len(file.First.Data) != 0 {
		t.Error("should be empty")
	}
}

func TestRemoveBeforeImpossible(t *testing.T) {
	file := File{}
	file.Init("")
	x := file.RemoveBefore()
	if x != 0 {
		t.Error("invalid no-op remove")
	}
}

func TestRemoveBefore(t *testing.T) {
	file := File{}
	file.Init("")
	for c := 'a'; c <= 'z'; c++ {
		file.Add(c)
	}
	size := 26
	for c := 'z'; c >= 'a'; c-- {
		size--
		x := file.RemoveBefore()
		if x != size {
			t.Error("bad offset")
		}
		for i := 0; i < size; i++ {
			u := rune('a' + i)
			if u != file.First.Data[i] {
				t.Error("bad removal")
			}
		}
	}
	x := file.Remove()
	if x != 0 {
		t.Error("should go to zero index")
	}
	if len(file.First.Data) != 0 {
		t.Error("should be empty")
	}
}

func TestBackspaceImpossible(t *testing.T) {
	file := File{}
	file.Init("")
	x, deleted := file.Backspace()
	if x != 0 {
		t.Error("should stay at zero index")
	}
	if deleted {
		t.Error("no line should have been deleted")
	}
}

func TestBackspace(t *testing.T) {
	file := File{}
	file.Init("")
	size := 0
	for c := 'a'; c <= 'z'; c++ {
		file.Add(c)
		size++
	}
	for c := 'z'; c >= 'a'; c-- {
		size--
		x, deleted := file.Backspace()
		if deleted {
			t.Error("no line should have been deleted")
		}
		if size != len(file.First.Data) {
			t.Error("bad data size")
		}
		if x != size {
			t.Error("bad position")
		}
		for r := 'a'; r < c; r++ {
			if r != file.First.Data[r-'a'] {
				t.Error("bad data")
			}
		}
	}
}

func TestBackspaceMultipleLines(t *testing.T) {
	file := File{}
	file.Init("")
	for c := 'a'; c <= 'z'; c++ {
		file.Add(c)
	}
	file.Add('\n')
	for c := 'a'; c <= 'z'; c++ {
		file.Add(c)
	}
	size := 26
	for c := 'z'; c >= 'a'; c-- {
		size--
		x, deleted := file.Backspace()
		if deleted {
			t.Error("no line should have been deleted")
		}
		if size != len(file.First.Next.Data) {
			t.Error("bad data size")
		}
		if x != size {
			t.Error("bad position")
		}
		for r := 'a'; r < c; r++ {
			if r != file.First.Next.Data[r-'a'] {
				t.Error("bad data")
			}
		}
	}
	size = 26
	x, deleted := file.Backspace()
	if x != size {
		t.Error("backspace should return to right side offset")
	}
	if !deleted {
		t.Error("backspace should cause deleted line")
	}
	for c := 'z'; c >= 'a'; c-- {
		size--
		x, deleted := file.Backspace()
		if deleted {
			t.Error("no line should have been deleted")
		}
		if size != len(file.First.Data) {
			t.Error("bad data size")
		}
		if x != size {
			t.Error("bad position")
		}
		for r := 'a'; r < c; r++ {
			if r != file.First.Data[r-'a'] {
				t.Error("bad data")
			}
		}
	}
}

func TestRemoveLineNoData(t *testing.T) {
	file := File{}
	file.Init("")
	x, _, _ := file.RemoveLine(false)
	if x != 0 {
		t.Error("position should be zero")
	}
}

func TestRemoveLineFirstLine(t *testing.T) {
	file := File{}
	file.Init("")
	runes := []rune{'a', 'b', 'c', '\n', 'd', 'e', 'f', '\n', 'g', 'h', 'i'}
	for _, r := range runes {
		file.Add(r)
	}
	if file.Lines != 3 {
		t.Error("should have three lines")
	}
	file.Up(false)
	file.Up(false)
	file.RemoveLine(false)
	for i, r := range file.First.Data {
		if r != runes[i+4] {
			t.Error("bad rune")
		}
	}
	for i, r := range file.First.Next.Data {
		if r != runes[i+8] {
			t.Error("bad rune")
		}
	}
}

func TestRemoveLineLastLine(t *testing.T) {
	file := File{}
	file.Init("")
	runes := []rune{'a', 'b', 'c', '\n', 'd', 'e', 'f', '\n', 'g', 'h', 'i'}
	for _, r := range runes {
		file.Add(r)
	}
	if file.Lines != 3 {
		t.Error("should have three lines")
	}
	file.RemoveLine(false)
	for i, r := range file.First.Data {
		if r != runes[i] {
			t.Error("bad rune")
		}
	}
	for i, r := range file.First.Next.Data {
		if r != runes[i+4] {
			t.Error("bad rune")
		}
	}
}

func TestRemoveLineMiddleLine(t *testing.T) {
	file := File{}
	file.Init("")
	runes := []rune{'a', 'b', 'c', '\n', 'd', 'e', 'f', '\n', 'g', 'h', 'i'}
	for _, r := range runes {
		file.Add(r)
	}
	if file.Lines != 3 {
		t.Error("should have three lines")
	}
	file.Up(false)
	file.RemoveLine(false)
	for i, r := range file.First.Data {
		if r != runes[i] {
			t.Error("bad rune")
		}
	}
	for i, r := range file.First.Next.Data {
		if r != runes[i+8] {
			t.Error("bad rune")
		}
	}
}

func TestRemoveRestOfLineNoData(t *testing.T) {
	file := File{}
	file.Init("")
	x := file.RemoveRestOfLine(false)
	if x != 0 {
		t.Error("should be zero index")
	}
}

func TestRemoveRestOfLineZeroOffset(t *testing.T) {
	file := File{}
	file.Init("")
	for c := 'a'; c <= 'z'; c++ {
		file.Add(c)
	}
	for c := 'z'; c >= 'a'; c-- {
		file.Left()
	}
	x := file.RemoveRestOfLine(false)
	if x != 0 {
		t.Error("should be zero index")
	}
	if len(file.First.Data) != 0 {
		t.Error("should have deleted everything")
	}
}

func TestRemoveRestOfLineNonZeroOffset(t *testing.T) {
	file := File{}
	file.Init("")
	for c := 'a'; c <= 'z'; c++ {
		file.Add(c)
	}
	for c := 'z'; c >= 'h'; c-- {
		file.Left()
	}
	x := file.RemoveRestOfLine(false)
	if x != 'h'-'a'-1 {
		t.Errorf("incorrect index: %d", x)
	}
	if len(file.First.Data) != 'h'-'a' {
		t.Errorf("should have deleted to the right: got size %d", len(file.First.Data))
	}
}
