package buffer_old

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	resizeRatio = 1.5
	resizeAt    = 0.75
	minimumSize = 16
)

type File struct {
	fileName string
	buffer   []rune
	capacity int
	length   int
	// cursorIndex is the empty byte at the start of the gap
	cursorIndex int
	// backBlockIndex is the byte after the end of the gap
	backBlockIndex int
	mutated        bool
}

func (buf *File) Init(name string) (arr []rune) {
	buf.fileName = name
	buf.mutated = false
	arr, buf.length = readFile(name)
	buf.capacity = computeCapacity(buf.length)
	buf.backBlockIndex = buf.capacity - buf.length
	buf.buffer = make([]rune, buf.capacity)
	copy(buf.buffer[buf.backBlockIndex:], arr)
	return arr
}

func readFile(name string) (arr []rune, length int) {
	dat, err := ioutil.ReadFile(name)
	data := ""
	if !os.IsNotExist(err) {
		data = string(dat)
	}
	arr = []rune(data)
	if len(arr) == 0 || arr[len(arr)-1] != '\n' {
		arr = append(arr, '\n')
	}
	return arr, len(arr)
}

func computeCapacity(length int) (capacity int) {
	capacity = int(float64(length) * resizeRatio)
	if capacity < minimumSize {
		return minimumSize
	}
	return capacity
}

func (buf *File) Redraw(yCurrent, height int) (arr []rune) {
	start := buf.computeStart(yCurrent)
	end := buf.computeEnd(yCurrent, height)
	startBlock := buf.buffer[start:buf.cursorIndex]
	endBlock := buf.buffer[buf.backBlockIndex:end]
	arr = make([]rune, len(startBlock)+len(endBlock))
	pivot := buf.cursorIndex - start
	copy(arr[:pivot], startBlock)
	copy(arr[pivot:], endBlock)
	return arr
}

func (buf *File) computeStart(yCurrent int) (start int) {
	start = buf.cursorIndex - 1
	for {
		if start < 0 || buf.buffer[start] == '\n' {
			yCurrent--
			if yCurrent < 0 {
				start++
				break
			}
		}
		start--
	}
	return start
}

func (buf *File) computeEnd(yCurrent, height int) (end int) {
	end = buf.backBlockIndex
	for {
		if end >= buf.capacity {
			break
		}
		if buf.buffer[end] == '\n' {
			yCurrent++
			if yCurrent > height-1 {
				break
			}
		}
		end++
	}
	return end
}

func (buf *File) Search(word string, yCurrent, height int) (xPoints, yPoints []int) {
	xPoints, yPoints = make([]int, 0), make([]int, 0)
	x, y := 0, 0
	arr := buf.Redraw(yCurrent, height)
	for i := 0; i < len(arr); i++ {
		if arr[i] == '\n' {
			x = 0
			y++
			continue
		}
		if isMatching(arr, word, i) {
			xPoints = append(xPoints, x)
			yPoints = append(yPoints, y)
		}
		x++
	}
	return xPoints, yPoints
}

func isMatching(arr []rune, word string, index int) (isMatching bool) {
	arrLen := len(arr)
	wordArr := []rune(word)
	for i := 0; i < len(word); i++ {
		if i+index >= arrLen || wordArr[i] != arr[i+index] {
			return false
		}
	}
	return true
}

func (buf *File) Save() (err error) {
	file, err := os.Create(buf.fileName)
	if err != nil {
		return err
	}
	var arr []rune
	arr = make([]rune, buf.length)
	copy(arr[:buf.cursorIndex], buf.buffer[:buf.cursorIndex])
	copy(arr[buf.cursorIndex:], buf.buffer[buf.backBlockIndex:])
	_, err = fmt.Fprintf(file, string(arr))
	if err != nil {
		return err
	}
	_ = file.Close()
	buf.mutated = false
	return nil
}

func (buf *File) CanSafeQuit() (isPossible bool) {
	return !buf.mutated
}

func (buf *File) Add(add rune) (requiredUpdates int) {
	buf.buffer[buf.cursorIndex] = add
	buf.length++
	buf.cursorIndex++
	if buf.length > int(float64(buf.capacity)*resizeAt) {
		buf.capacity = int(float64(buf.capacity) * resizeRatio)
		temp := make([]rune, buf.capacity)
		backLength := buf.length - buf.cursorIndex
		newBackBlockIndex := buf.capacity - backLength
		copy(temp[:buf.cursorIndex], buf.buffer[:buf.cursorIndex])
		copy(temp[newBackBlockIndex:], buf.buffer[buf.backBlockIndex:])
		buf.backBlockIndex = newBackBlockIndex
		buf.buffer = temp
	}
	buf.mutated = true
	return buf.computeRequiredUpdates()
}

func (buf *File) Remove() (possible bool, newX, requiredUpdates, spacing int) {
	if buf.cursorIndex == 0 {
		return false, 0, 0, 0
	}
	if buf.buffer[buf.cursorIndex-1] == '\t' {
		spacing = 4 + 1
	} else {
		spacing = 1
	}
	buf.length--
	buf.cursorIndex--
	buf.mutated = true
	return true, buf.computeNewX(), buf.computeRequiredUpdates(), spacing
}

func (buf *File) computeNewX() (newX int) {
	newX = 0
	for i := buf.cursorIndex - 1; i >= 0 && buf.buffer[i] != '\n'; i-- {
		newX++
	}
	return newX
}

func (buf *File) computeRequiredUpdates() (requiredUpdates int) {
	for i := buf.backBlockIndex; i < buf.capacity && buf.buffer[i] != '\n'; i++ {
		requiredUpdates++
	}
	return requiredUpdates
}

func (buf *File) RemoveCurrent() (possible, xBack bool, requiredUpdates int) {
	if buf.backBlockIndex == buf.length || buf.buffer[buf.backBlockIndex] == '\n' {
		return false, false, 0
	}
	buf.length--
	buf.backBlockIndex++
	buf.mutated = true
	return true, buf.buffer[buf.backBlockIndex] == '\n', buf.computeRequiredUpdates() + 1
}

func (buf *File) RemoveLine() (yBack, isEmpty bool) {
	for i := buf.cursorIndex - 1; i >= 0 && buf.buffer[i] != '\n'; i-- {
		buf.cursorIndex--
		buf.length--
	}
	for i := buf.backBlockIndex; buf.buffer[i] != '\n'; i++ {
		buf.backBlockIndex++
		buf.length--
	}
	if buf.length > 1 {
		buf.backBlockIndex++
		buf.length--
	}
	buf.mutated = true
	return buf.backBlockIndex == buf.capacity, buf.length == 1
}

func (buf *File) RemoveRestOfLine() (requiredUpdates int) {
	requiredUpdates = 0
	for i := buf.backBlockIndex; buf.buffer[i] != '\n'; i++ {
		requiredUpdates++
		buf.backBlockIndex++
		buf.length--
	}
	buf.mutated = true
	return requiredUpdates
}

func (buf *File) Left() (characters int) {
	if buf.cursorIndex == 0 || buf.buffer[buf.cursorIndex-1] == '\n' {
		return 0
	}
	char := buf.buffer[buf.backBlockIndex]
	buf.cursorIndex--
	buf.backBlockIndex--
	buf.buffer[buf.backBlockIndex] = buf.buffer[buf.cursorIndex]
	if char == '\t' {
		return 4 + 1
	}
	return 1
}

func (buf *File) Right(isInsert bool) (characters int) {
	offset := computeOffset(isInsert)
	if buf.backBlockIndex == buf.capacity-offset ||
		buf.buffer[buf.backBlockIndex] == '\n' ||
		buf.buffer[buf.backBlockIndex+offset] == '\n' {
		return 0
	}
	buf.buffer[buf.cursorIndex] = buf.buffer[buf.backBlockIndex]
	buf.cursorIndex++
	buf.backBlockIndex++
	if buf.buffer[buf.backBlockIndex+offset-1] == '\t' {
		return 4 + 1
	}
	return 1
}

func (buf *File) Up(oldX int, isInsert bool) (possible bool, newX int) {
	if buf.cursorIndex == 0 {
		return false, oldX
	}
	i := buf.cursorIndex - 1
	count := 0
	for i > 0 && buf.buffer[i] != '\n' {
		i--
		count++
	}
	if i == 0 && buf.buffer[i] != '\n' {
		return false, oldX
	}
	i--
	count++
	if i == -1 || buf.buffer[i] == '\n' {
		temp := buf.buffer[i+1 : buf.cursorIndex]
		copy(buf.buffer[buf.backBlockIndex-count:buf.backBlockIndex], temp)
		buf.cursorIndex = i + 1
		buf.backBlockIndex -= count
		return true, 0
	}
	lineLen := 0
	for j := i; j > 0 && buf.buffer[j-1] != '\n'; j-- {
		i--
		count++
		lineLen++
	}
	if lineLen < oldX {
		newX = lineLen
		if isInsert {
			newX++
		}
	} else {
		newX = oldX
	}
	i += newX - 1
	count -= newX - 1
	temp := buf.buffer[i+1 : buf.cursorIndex]
	copy(buf.buffer[buf.backBlockIndex-count:buf.backBlockIndex], temp)
	buf.cursorIndex = i + 1
	buf.backBlockIndex -= count
	return true, newX
}

func (buf *File) Down(oldX int, isInsert bool) (possible bool, newX int) {
	offset := computeOffset(isInsert)
	i := buf.backBlockIndex
	for i < buf.capacity && buf.buffer[i] != '\n' {
		i++
	}
	i++
	if i >= buf.capacity {
		return false, oldX
	}
	if buf.buffer[i] != '\n' {
		for newX = 0; newX < oldX && i < buf.capacity-offset && buf.buffer[i+offset] != '\n'; newX++ {
			i++
		}
	}
	temp := buf.buffer[buf.backBlockIndex:i]
	delta := i - buf.backBlockIndex
	copy(buf.buffer[buf.cursorIndex:buf.cursorIndex+delta], temp)
	buf.cursorIndex += delta
	buf.backBlockIndex = i
	return true, newX
}

func computeOffset(isInsert bool) (offset int) {
	if isInsert {
		return 0
	}
	return 1
}

func (buf *File) GetBottom(currentY, getY int) (bottom string) {
	var i int
	deltaY := getY - currentY
	for i = buf.backBlockIndex; i < buf.capacity; i++ {
		if buf.buffer[i] == '\n' {
			deltaY--
			if deltaY == 0 {
				break
			}
		}
	}
	var sb strings.Builder
	for j := i + 1; j < buf.capacity && buf.buffer[j] != '\n'; j++ {
		sb.WriteRune(buf.buffer[j])
	}
	if sb.Len() == 0 {
		return "~"
	}
	return sb.String()
}

func (buf *File) GetLine() (previous string) {
	var sb strings.Builder
	startIndex := buf.cursorIndex - 1
	for i := startIndex; i >= 0 && buf.buffer[i] != '\n'; i-- {
		startIndex--
	}
	for i := startIndex + 1; i < buf.cursorIndex; i++ {
		sb.WriteRune(buf.buffer[i])
	}
	for i := buf.backBlockIndex; i < buf.capacity && buf.buffer[i] != '\n'; i++ {
		sb.WriteRune(buf.buffer[i])
	}
	return sb.String()
}
