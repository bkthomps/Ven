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
	line.Data = line.Prev.Data[file.runeOffset:]
	line.Prev.Data = line.Prev.Data[:file.runeOffset]
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

func (file *File) RemoveBefore() (xPosition int) {
	if file.runeOffset == 0 {
		return file.spacingOffset
	}
	if file.Current.Data[file.runeOffset] == '\t' {
		file.spacingOffset -= TabSize
	} else {
		file.spacingOffset--
	}
	file.runeOffset--
	file.Current.RemoveAt(file.runeOffset)
	file.mutated = true
	return file.spacingOffset
}

func (file *File) Backspace() (xPosition int, deletedLine bool) {
	file.mutated = true
	if file.runeOffset == 0 {
		if file.Current == file.First {
			return file.spacingOffset, false
		}
		file.spacingOffset = 1_000_000_000
		current := file.Current
		file.Current = current.Prev
		file.calculateOffset(true)
		current.Prev.Data = append(current.Prev.Data, current.Data...)
		current.Prev.Next = current.Next
		if current.Next != nil {
			current.Next.Prev = current.Prev
		} else {
			file.last = file.Current
		}
		file.lines--
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

func (file *File) RemoveLine(isInsert bool) (xPosition int, lineUp bool) {
	file.mutated = true
	if file.First == file.last {
		line := &Line{}
		line.Init(nil, nil)
		file.First = line
		file.last = line
		file.Current = line
		file.runeOffset = 0
		file.spacingOffset = 0
		return file.spacingOffset, false
	}
	if file.Current == file.First {
		file.Current = file.Current.Next
		file.Current.Prev = nil
		file.First = file.Current
		file.lines--
		file.calculateOffset(isInsert)
		return file.spacingOffset, false
	}
	if file.Current == file.last {
		file.Current = file.Current.Prev
		file.Current.Next = nil
		file.last = file.Current
		file.lines--
		file.calculateOffset(isInsert)
		return file.spacingOffset, true
	}
	deleteNode := file.Current
	deleteNode.Prev.Next = deleteNode.Next
	deleteNode.Next.Prev = deleteNode.Prev
	file.Current = deleteNode.Next
	file.lines--
	file.calculateOffset(isInsert)
	return file.spacingOffset, false
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
