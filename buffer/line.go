package buffer

type Line struct {
	Data []rune
	Next *Line
	Prev *Line
}

func (line *Line) Init(next, prev *Line) {
	line.Data = []rune{}
	line.Next = next
	line.Prev = prev
}

func (line *Line) AddAt(index int, character rune) {
	line.Data = append(line.Data, 0)
	copy(line.Data[index+1:], line.Data[index:])
	line.Data[index] = character
}

func (line *Line) RemoveAt(index int) {
	line.Data = append(line.Data[:index], line.Data[index+1:]...)
}
