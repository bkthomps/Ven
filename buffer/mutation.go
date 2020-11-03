package buffer

func (file *File) Add(character rune) (xPosition int, addedLine bool) {
	file.mutated = true
	if character == '\n' {
		file.addLine()
		file.runeOffset = 0
		file.spacingOffset = 0
		return file.spacingOffset, true
	}
	if character == '\t' {
		file.spacingOffset += TabSize
	} else {
		file.spacingOffset++
	}
	file.Current.AddAt(file.runeOffset, character)
	file.runeOffset++
	return file.spacingOffset, false
}

func (file *File) addLine() {
	line := &Line{}
	line.Init(file.Current.Next, file.Current)
	if file.Current.Next != nil {
		file.Current.Next.Prev = line
	} else {
		file.last = line
	}
	file.Current.Next = line
	file.Current = line
	file.lines++
}

func (file *File) Remove() (xPosition int) {
	if len(file.Current.Data) > 0 {
		file.mutated = true
		if file.runeOffset == len(file.Current.Data)-1 {
			if file.Current.Data[file.runeOffset] == '\t' {
				file.spacingOffset -= TabSize
			} else {
				file.spacingOffset--
			}
			file.runeOffset--
		}
		file.Current.RemoveAt(file.runeOffset)
	}
	return file.spacingOffset
}

func (file *File) RemoveBefore() (xPosition int, deletedLine bool) {
	file.mutated = true
	if file.runeOffset == 0 {
		// TODO: merge with precedent line
		return file.spacingOffset, true
	}
	file.runeOffset--
	if file.Current.Data[file.runeOffset] == '\t' {
		file.spacingOffset -= TabSize
	} else {
		file.spacingOffset--
	}
	file.Current.RemoveAt(file.runeOffset)
	return file.spacingOffset, false
}

func (file *File) RemoveLine(isInsert bool) (xPosition int) {
	file.mutated = true
	if file.lines == 1 {
		line := &Line{}
		line.Init(nil, nil)
		file.First = line
		file.last = line
		file.Current = line
		file.runeOffset = 0
		file.spacingOffset = 0
		return file.spacingOffset
	}
	if file.Current == file.First {
		file.Current = file.Current.Next
		file.Current.Prev = nil
		file.First = file.Current
		file.lines--
		file.calculateOffset(isInsert)
		return file.spacingOffset
	}
	if file.Current == file.last {
		file.Current = file.Current.Prev
		file.Current.Next = nil
		file.last = file.Current
		file.lines--
		file.calculateOffset(isInsert)
		return file.spacingOffset
	}
	node := file.Current.Next
	file.Current.Next = node.Next
	node.Next.Prev = file.Current
	file.lines--
	file.calculateOffset(isInsert)
	return file.spacingOffset
}

func (file *File) RemoveRestOfLine(isInsert bool) (xPosition int) {
	file.mutated = true
	if file.runeOffset == 0 {
		file.Current.Data = []rune{}
		file.runeOffset = 0
		file.spacingOffset = 0
		return file.spacingOffset
	}
	file.Current.Data = file.Current.Data[:file.runeOffset]
	file.calculateOffset(isInsert)
	return file.spacingOffset
}
