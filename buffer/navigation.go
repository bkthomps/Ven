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
	if file.Current == file.last {
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
	for len(file.Current.Data) == 0 || !unicode.IsSpace(file.Current.Data[file.runeOffset]) {
		if file.runeOffset >= len(file.Current.Data)-1 {
			if file.Current.Next == nil {
				return file.spacingOffset, linesDown
			}
			file.spacingOffset = 0
			file.runeOffset = 0
			linesDown++
			file.Current = file.Current.Next
			break
		}
		file.spacingOffset = file.runeWidthIncrease(file.Current.Data[file.runeOffset])
		file.runeOffset++
	}
	for len(file.Current.Data) == 0 || unicode.IsSpace(file.Current.Data[file.runeOffset]) {
		if file.runeOffset >= len(file.Current.Data)-1 {
			if file.Current.Next == nil {
				return file.spacingOffset, linesDown
			}
			file.spacingOffset = 0
			file.runeOffset = 0
			linesDown++
			file.Current = file.Current.Next
			continue
		}
		file.spacingOffset = file.runeWidthIncrease(file.Current.Data[file.runeOffset])
		file.runeOffset++
	}
	return file.spacingOffset, linesDown
}

// PrevWordStart will move the cursor to the start of the current word,
// unless the cursor is at the start of a word or on whitespace, in
// which case it will move to the start of the previous word, if there
// is one, otherwise it will move the cursor to the end of the file.
func (file *File) PrevWordStart() (xPosition int, linesUp int) {
	return 0, 0
}

// NextWordEnd will move the cursor to the end of the current word,
// unless the cursor is at the end of a word or on whitespace, in
// which case it will move to the end of the next word, if there
// is one, otherwise it will move the cursor to the end of the file.
func (file *File) NextWordEnd() (xPosition int, linesDown int) {
	linesDown = 0
	if file.runeOffset < len(file.Current.Data) {
		file.spacingOffset = file.runeWidthIncrease(file.Current.Data[file.runeOffset])
		file.runeOffset++
	}
	for file.runeOffset >= len(file.Current.Data) || unicode.IsSpace(file.Current.Data[file.runeOffset]) {
		if file.runeOffset >= len(file.Current.Data) {
			if file.Current.Next == nil {
				return file.spacingOffset, linesDown
			}
			file.spacingOffset = 0
			file.runeOffset = 0
			linesDown++
			file.Current = file.Current.Next
			continue
		}
		file.spacingOffset = file.runeWidthIncrease(file.Current.Data[file.runeOffset])
		file.runeOffset++
	}
	for file.runeOffset >= len(file.Current.Data)-1 || !unicode.IsSpace(file.Current.Data[file.runeOffset+1]) {
		if file.runeOffset >= len(file.Current.Data)-1 {
			return file.spacingOffset, linesDown
		}
		file.spacingOffset = file.runeWidthIncrease(file.Current.Data[file.runeOffset])
		file.runeOffset++
	}
	return file.spacingOffset, linesDown
}
