/*
Copyright (c) 2019 Bailey Thompson

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package buffer

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const resizeRatio = 1.5
const resizeAt = 0.75
const minimumSize = 16

type Info struct {
	fileName       string
	buffer         []rune
	capacity       int
	length         int
	cursorIndex    int // The empty byte at the start of the gap
	backBlockIndex int // The byte after the end of the gap
	mutated        bool
}

func Init(name string) (buf *Info, arr []rune) {
	arr, length := readFile(name)
	capacity := computeCapacity(length)
	backBlockIndex := capacity - length
	buffer := make([]rune, capacity)
	copy(buffer[backBlockIndex:], arr)
	buf = &Info{
		fileName:       name,
		buffer:         buffer,
		capacity:       capacity,
		length:         length,
		backBlockIndex: backBlockIndex,
		mutated:        false,
	}
	return buf, arr
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

func Redraw(buf *Info, yCurrent, height int) (arr []rune) {
	start := computeStart(buf, yCurrent)
	end := computeEnd(buf, yCurrent, height)
	startBlock := buf.buffer[start:buf.cursorIndex]
	endBlock := buf.buffer[buf.backBlockIndex:end]
	arr = make([]rune, len(startBlock)+len(endBlock))
	pivot := buf.cursorIndex - start
	copy(arr[:pivot], startBlock)
	copy(arr[pivot:], endBlock)
	return arr
}

func computeStart(buf *Info, yCurrent int) (start int) {
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

func computeEnd(buf *Info, yCurrent, height int) (end int) {
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

func Search(buf *Info, word string, yCurrent, height int) (xPoints, yPoints []int) {
	xPoints, yPoints = make([]int, 0), make([]int, 0)
	x, y := 0, 0
	arr := Redraw(buf, yCurrent, height)
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

func Save(buf *Info) (err error) {
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

func CanSafeQuit(buf *Info) (isPossible bool) {
	return !buf.mutated
}

func Add(buf *Info, add rune) (requiredUpdates int) {
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
	return computeRequiredUpdates(buf)
}

func Remove(buf *Info) (possible bool, newX, requiredUpdates int) {
	if buf.cursorIndex == 0 {
		return false, 0, 0
	}
	buf.length--
	buf.cursorIndex--
	buf.mutated = true
	return true, computeNewX(buf), computeRequiredUpdates(buf)
}

func computeNewX(buf *Info) (newX int) {
	newX = 0
	for i := buf.cursorIndex - 1; i >= 0 && buf.buffer[i] != '\n'; i-- {
		newX++
	}
	return newX
}

func computeRequiredUpdates(buf *Info) (requiredUpdates int) {
	for i := buf.backBlockIndex; i < buf.capacity && buf.buffer[i] != '\n'; i++ {
		requiredUpdates++
	}
	return requiredUpdates
}

func RemoveCurrent(buf *Info) (possible, xBack bool, requiredUpdates int) {
	if buf.backBlockIndex == buf.length || buf.buffer[buf.backBlockIndex] == '\n' {
		return false, false, 0
	}
	buf.length--
	buf.backBlockIndex++
	buf.mutated = true
	return true, buf.buffer[buf.backBlockIndex] == '\n', computeRequiredUpdates(buf) + 1
}

func RemoveLine(buf *Info) (yBack, isEmpty bool) {
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

func RemoveRestOfLine(buf *Info) (requiredUpdates int) {
	requiredUpdates = 0
	for i := buf.backBlockIndex; buf.buffer[i] != '\n'; i++ {
		requiredUpdates++
		buf.backBlockIndex++
		buf.length--
	}
	buf.mutated = true
	return requiredUpdates
}

func Left(buf *Info) (possible bool) {
	if buf.cursorIndex == 0 || buf.buffer[buf.cursorIndex-1] == '\n' {
		return false
	}
	buf.cursorIndex--
	buf.backBlockIndex--
	buf.buffer[buf.backBlockIndex] = buf.buffer[buf.cursorIndex]
	return true
}

func Right(buf *Info, isInsert bool) (possible bool) {
	offset := computeOffset(isInsert)
	if buf.backBlockIndex == buf.capacity-offset ||
		buf.buffer[buf.backBlockIndex] == '\n' ||
		buf.buffer[buf.backBlockIndex+offset] == '\n' {
		return false
	}
	buf.buffer[buf.cursorIndex] = buf.buffer[buf.backBlockIndex]
	buf.cursorIndex++
	buf.backBlockIndex++
	return true
}

func Up(buf *Info, oldX int, isInsert bool) (possible bool, newX int) {
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

func Down(buf *Info, oldX int, isInsert bool) (possible bool, newX int) {
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

func GetBottom(buf *Info, currentY, getY int) (bottom string) {
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

func GetLine(buf *Info) (previous string) {
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
