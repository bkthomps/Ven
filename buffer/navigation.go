package buffer

import "unicode"

func (file *File) Left() (xPosition int) {
	if file.runeOffset == 0 {
		return 0
	}
	file.runeOffset--
	r := file.Current.Data[file.runeOffset]
	file.spacingOffset = file.runeWidthDecrease(r)
	return file.spacingOffset
}

func (file *File) Right(isInsert bool) (xPosition int) {
	insertOffset := 0
	if isInsert {
		insertOffset = 1
	}
	if file.runeOffset >= len(file.Current.Data)+insertOffset-1 {
		return file.spacingOffset
	}
	r := file.Current.Data[file.runeOffset]
	file.spacingOffset = file.runeWidthIncrease(r)
	file.runeOffset++
	return file.spacingOffset
}

func (file *File) Up(isInsert bool) (wasPossible bool, xPosition int) {
	if file.Current == file.First {
		return false, file.spacingOffset
	}
	file.Current = file.Current.Prev
	file.calculateOffset(isInsert)
	return true, file.spacingOffset
}

func (file *File) Down(isInsert bool) (wasPossible bool, xPosition int) {
	if file.Current.Next == nil {
		return false, file.spacingOffset
	}
	file.Current = file.Current.Next
	file.calculateOffset(isInsert)
	return true, file.spacingOffset
}

func (file *File) calculateOffset(isInsert bool) {
	oldSpacingOffset := file.spacingOffset
	file.runeOffset = 0
	file.spacingOffset = 0
	for i, r := range file.Current.Data {
		if !isInsert && i == len(file.Current.Data)-1 {
			return
		}
		if file.runeWidthIncrease(r) > oldSpacingOffset {
			return
		}
		file.spacingOffset = file.runeWidthIncrease(r)
		file.runeOffset++
	}
}

func (file *File) StartOfLine() (xPosition int) {
	file.runeOffset = 0
	file.spacingOffset = 0
	return file.spacingOffset
}

func (file *File) EndOfLine(isInsert bool) (xPosition int) {
	insertOffset := 0
	if isInsert {
		insertOffset = 1
	}
	for i := file.runeOffset; i < len(file.Current.Data)+insertOffset-1; i++ {
		r := file.Current.Data[file.runeOffset]
		file.spacingOffset = file.runeWidthIncrease(r)
		file.runeOffset++
	}
	return file.spacingOffset
}

func (file *File) JumpToTop() (xPosition int) {
	file.Current = file.First
	return file.StartOfLine()
}

func (file *File) JumpToBottom() (xPosition int) {
	file.Current = file.last
	return file.StartOfLine()
}

// NextWordStart will move the cursor to the start of the next word,
// unless there is no next word, in which case the cursor moves to
// the end of the file.
func (file *File) NextWordStart() (xPosition int, linesDown int) {
	linesDown = 0
	for file.isNotWhitespace(file.runeOffset) {
		if file.isOutOfBounds(file.runeOffset) {
			if file.Current.Next == nil {
				return file.spacingOffset, linesDown
			}
			linesDown = file.moveLineDown(linesDown)
			break
		}
		file.moveForward()
	}
	for file.isWhitespace(file.runeOffset) {
		if file.isOutOfBounds(file.runeOffset) {
			if file.Current.Next == nil {
				return file.spacingOffset, linesDown
			}
			linesDown = file.moveLineDown(linesDown)
			continue
		}
		file.moveForward()
	}
	return file.spacingOffset, linesDown
}

// PrevWordStart will move the cursor to the start of the current word,
// unless the cursor is at the start of a word or on whitespace, in
// which case it will move to the start of the previous word, if there
// is one, otherwise it will move the cursor to the start of the file.
func (file *File) PrevWordStart() (xPosition int, linesUp int) {
	linesUp = 0
	file.spacingOffset = 0
	if file.runeOffset >= 0 {
		file.runeOffset--
	}
	for file.isWhitespace(file.runeOffset) {
		if file.isOutOfBounds(file.runeOffset) {
			if file.Current.Prev == nil {
				file.runeOffset = 0
				return file.spacingOffset, linesUp
			}
			linesUp = file.moveLineUp(linesUp)
			continue
		}
		file.runeOffset--
	}
	for file.isNotWhitespace(file.runeOffset - 1) {
		if file.isOutOfBounds(file.runeOffset - 1) {
			file.runeOffset = 0
			return file.spacingOffset, linesUp
		}
		file.runeOffset--
	}
	for i := 0; i < file.runeOffset; i++ {
		file.spacingOffset = file.runeWidthIncrease(file.Current.Data[i])
	}
	return file.spacingOffset, linesUp
}

// NextWordEnd will move the cursor to the end of the current word,
// unless the cursor is at the end of a word or on whitespace, in
// which case it will move to the end of the next word, if there
// is one, otherwise it will move the cursor to the end of the file.
func (file *File) NextWordEnd() (xPosition int, linesDown int) {
	linesDown = 0
	if file.runeOffset < len(file.Current.Data) {
		file.moveForward()
	}
	for file.isWhitespace(file.runeOffset) {
		if file.isOutOfBounds(file.runeOffset) {
			if file.Current.Next == nil {
				return file.spacingOffset, linesDown
			}
			linesDown = file.moveLineDown(linesDown)
			continue
		}
		file.moveForward()
	}
	for file.isNotWhitespace(file.runeOffset + 1) {
		if file.isOutOfBounds(file.runeOffset + 1) {
			return file.spacingOffset, linesDown
		}
		file.moveForward()
	}
	return file.spacingOffset, linesDown
}

func (file *File) isWhitespace(index int) bool {
	return file.isOutOfBounds(index) || unicode.IsSpace(file.Current.Data[index])
}

func (file *File) isNotWhitespace(index int) bool {
	return file.isOutOfBounds(index) || !unicode.IsSpace(file.Current.Data[index])
}

func (file *File) isOutOfBounds(index int) bool {
	return index < 0 || index >= len(file.Current.Data)
}

func (file *File) moveLineDown(linesDown int) int {
	file.spacingOffset = 0
	file.runeOffset = 0
	file.Current = file.Current.Next
	return linesDown + 1
}

func (file *File) moveLineUp(linesUp int) int {
	file.Current = file.Current.Prev
	file.runeOffset = len(file.Current.Data) - 1
	return linesUp + 1
}

func (file *File) moveForward() {
	file.spacingOffset = file.runeWidthIncrease(file.Current.Data[file.runeOffset])
	file.runeOffset++
}
