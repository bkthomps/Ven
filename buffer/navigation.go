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
	insertOffset := 0
	if isInsert {
		insertOffset = 1
	}
	if file.runeOffset >= len(file.Current.Data)+insertOffset-1 {
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

func (file *File) calculateOffset(isInsert bool) {
	if len(file.Current.Data) == 0 {
		file.runeOffset = 0
		file.spacingOffset = 0
		return
	}
	oldSpacingOffset := file.spacingOffset
	file.runeOffset = -1
	file.spacingOffset = -1
	for _, r := range file.Current.Data {
		currentSpacing := 1
		if r == '\t' {
			currentSpacing = TabSize
		}
		if file.spacingOffset+currentSpacing > oldSpacingOffset {
			if file.runeOffset < 0 {
				file.runeOffset = 0
				file.spacingOffset = 0
			}
			return
		}
		file.runeOffset++
		file.spacingOffset += currentSpacing
	}
	if file.spacingOffset == oldSpacingOffset {
		return
	}
	if isInsert {
		file.runeOffset++
		file.spacingOffset++
	}
}
