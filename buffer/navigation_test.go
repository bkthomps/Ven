package buffer

import "testing"

func TestLeftNoMove(t *testing.T) {
	file := File{}
	file.Init("")
	x := file.Left()
	if x != 0 {
		t.Error("should be no movement")
	}
}

func TestLeftMove(t *testing.T) {
	file := File{}
	file.Init("")
	x, _ := file.Add('a')
	if x != 1 {
		t.Error("should auto-move right")
	}
	x = file.Left()
	if x != 0 {
		t.Error("should have moved left")
	}
}

func TestRightNoMove(t *testing.T) {
	file := File{}
	file.Init("")
	x := file.Right(false)
	if x != 0 {
		t.Error("should be no movement")
	}
}

func TestRightMove(t *testing.T) {
	file := File{}
	file.Init("")
	file.Add('a')
	file.Add('b')
	file.Left()
	file.Left()
	x := file.Right(false)
	if x != 1 {
		t.Error("should have moved right")
	}
}

func TestRightMoveInsert(t *testing.T) {
	file := File{}
	file.Init("")
	file.Add('a')
	file.Left()
	x := file.Right(false)
	if x != 0 {
		t.Error("should not have moved")
	}
}

func TestRightCannotMove(t *testing.T) {
	file := File{}
	file.Init("")
	file.Add('a')
	file.Left()
	x := file.Right(true)
	if x != 1 {
		t.Error("should have moved right")
	}
}

func TestUp(t *testing.T) {
	file := File{}
	file.Init("")
	runes := []rune{'a', 'b', 'c', '\n', 'd', 'e', 'f', '\n', 'g', 'h', 'i'}
	for _, r := range runes {
		file.Add(r)
	}
	possible, _ := file.Up(false)
	if !possible {
		t.Error("should have moved up")
	}
	possible, _ = file.Up(false)
	if !possible {
		t.Error("should have moved up")
	}
	possible, _ = file.Up(false)
	if possible {
		t.Error("should not have moved up")
	}
}

func TestDown(t *testing.T) {
	file := File{}
	file.Init("")
	runes := []rune{'a', 'b', 'c', '\n', 'd', 'e', 'f', '\n', 'g', 'h', 'i'}
	for _, r := range runes {
		file.Add(r)
	}
	file.Up(false)
	file.Up(false)
	possible, _ := file.Down(false)
	if !possible {
		t.Error("should have moved down")
	}
	possible, _ = file.Down(false)
	if !possible {
		t.Error("should have moved down")
	}
	possible, _ = file.Down(false)
	if possible {
		t.Error("should not have moved down")
	}
}

func TestNextWordStart(t *testing.T) {
	file := File{}
	file.Init("")
	runes := []rune{'a', 'b', 'c', ' ', 'd', 'e', 'f', ' ', 'g', 'h', 'i'}
	for _, r := range runes {
		file.Add(r)
	}
	file.StartOfLine()
	x, _ := file.NextWordStart()
	if x != 4 {
		t.Error("did not go to next word start")
	}
	x, _ = file.NextWordStart()
	if x != 8 {
		t.Error("did not go to next word start")
	}
}

func TestPrevWordStart(t *testing.T) {
	file := File{}
	file.Init("")
	runes := []rune{'a', 'b', 'c', ' ', 'd', 'e', 'f', ' ', 'g', 'h', 'i'}
	for _, r := range runes {
		file.Add(r)
	}
	x, _ := file.PrevWordStart()
	if x != 8 {
		t.Error("did not go to prev word start")
	}
	x, _ = file.PrevWordStart()
	if x != 4 {
		t.Error("did not go to prev word start")
	}
	x, _ = file.PrevWordStart()
	if x != 0 {
		t.Error("did not go to prev word start")
	}
}

func TestNextWordEnd(t *testing.T) {
	file := File{}
	file.Init("")
	runes := []rune{'a', 'b', 'c', ' ', 'd', 'e', 'f', ' ', 'g', 'h', 'i'}
	for _, r := range runes {
		file.Add(r)
	}
	file.StartOfLine()
	x, _ := file.NextWordEnd()
	if x != 2 {
		t.Error("did not go to next word end")
	}
	x, _ = file.NextWordEnd()
	if x != 6 {
		t.Error("did not go to next word end")
	}
	x, _ = file.NextWordEnd()
	if x != 10 {
		t.Error("did not go to next word end")
	}
}
