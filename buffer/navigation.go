package buffer

func (file *File) Left() (spaces int) {
	if file.offset == 0 {
		return 0
	}
	file.offset--
	if file.Current.Data[file.offset] == '\t' {
		return TabSize
	}
	return 1
}

func (file *File) Right(isInsert bool) (spaces int) {
	insertOffset := 0
	if isInsert {
		insertOffset = 1
	}
	if file.offset >= len(file.Current.Data)+insertOffset-1 {
		return 0
	}
	if file.Current.Data[file.offset] == '\t' {
		file.offset++
		return TabSize
	}
	file.offset++
	return 1
}

func (file *File) Up(isInsert bool) (wasPossible bool, xPosition int) {
	if file.Current == file.First {
		return false, 0
	}
	insertOffset := 0
	if isInsert {
		insertOffset = 1
	}
	file.Current = file.Current.Prev
	if file.offset >= len(file.Current.Data)+insertOffset {
		// TODO: this should take into account tabs
	}
	// TODO: don't always return 0
	file.offset = 0
	return true, file.offset
}

func (file *File) Down(isInsert bool) (wasPossible bool, xPosition int) {
	if file.Current == file.last {
		return false, 0
	}
	insertOffset := 0
	if isInsert {
		insertOffset = 1
	}
	file.Current = file.Current.Next
	if file.offset >= len(file.Current.Data)+insertOffset {
		// TODO: this should take into account tabs
	}
	// TODO: don't always return 0
	file.offset = 0
	return true, file.offset
}
