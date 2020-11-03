package buffer

func (file *File) Left() (xPosition int) {
	if file.runeOffset == 0 {
		return 0
	}
	file.runeOffset--
	if file.Current.Data[file.runeOffset] == '\t' {
		file.spacingOffset -= TabSize
	} else {
		file.spacingOffset--
	}
	return file.spacingOffset
}

func (file *File) Right(isInsert bool) (xPosition int) {
	if file.runeOffset >= len(file.Current.Data)+insertOffset(isInsert)-1 {
		return file.spacingOffset
	}
	if file.Current.Data[file.runeOffset] == '\t' {
		file.spacingOffset += TabSize
	} else {
		file.spacingOffset++
	}
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

func insertOffset(isInsert bool) int {
	if isInsert {
		return 1
	}
	return 0
}

func (file *File) calculateOffset(isInsert bool) {
	oldSpacingOffset := file.spacingOffset
	file.runeOffset = 0
	file.spacingOffset = 0
	for _, r := range file.Current.Data {
		currentSpacing := 1
		if r == '\t' {
			currentSpacing = TabSize
		}
		if file.spacingOffset+currentSpacing > oldSpacingOffset+insertOffset(isInsert) {
			break
		}
		file.runeOffset++
		file.spacingOffset += currentSpacing
	}
}
