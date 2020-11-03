package buffer

func (file *File) Add(character rune) {
	if character == '\n' {
		file.addLine()
		return
	}
	file.Current.AddAt(file.offset, character)
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
