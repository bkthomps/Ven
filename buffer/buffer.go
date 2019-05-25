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
	"math"
	"os"
	"strings"
)

const resizeRatio = 1.5
const resizeAt = 0.75
const minimumSize = 16

var fileName string
var buffer []rune
var capacity = 0
var length = 0
var cursorIndex = 0    // The empty byte at the start of the gap
var backBlockIndex = 0 // The byte after the end of the gap
var mutated = false

func Init(name string) (arr []rune) {
	fileName = name
	dat, err := ioutil.ReadFile(name)
	data := ""
	if !os.IsNotExist(err) {
		data = string(dat)
	}
	arr = []rune(data)
	if len(arr) == 0 || arr[len(arr)-1] != '\n' {
		arr = append(arr, '\n')
	}
	length = len(arr)
	capacity = int(math.Max(float64(length)*resizeRatio, minimumSize))
	backBlockIndex = capacity - length
	buffer = make([]rune, capacity)
	copy(buffer[backBlockIndex:], arr)
	return arr
}

func Redraw(yCurrent, height int) (arr []rune) {
	// TODO: fix crashes associated with this function
	yTemp := yCurrent
	startIndex := cursorIndex - 1
	for yTemp >= 0 {
		if buffer[startIndex] == '\n' {
			yTemp--
		}
		startIndex--
	}
	yTemp = yCurrent
	endIndex := backBlockIndex
	for yTemp < height-1 {
		if buffer[endIndex] == '\n' {
			yTemp++
		}
		endIndex++
	}
	arr = make([]rune, endIndex-backBlockIndex+1+cursorIndex-startIndex)
	pivot := cursorIndex - startIndex
	copy(arr[:pivot], arr[startIndex:cursorIndex])
	copy(arr[pivot:], arr[backBlockIndex:endIndex+1])
	return arr
}

func Save() (err error) {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	var arr []rune
	arr = make([]rune, length)
	copy(arr[:cursorIndex], buffer[:cursorIndex])
	copy(arr[cursorIndex:], buffer[backBlockIndex:])
	_, err = fmt.Fprintf(file, string(arr))
	if err != nil {
		return err
	}
	mutated = false
	return nil
}

func Log(name, text string) {
	file, err := os.Create(name)
	if err != nil {
		return
	}
	defer file.Close()
	_, _ = fmt.Fprintf(file, text)
}

func CanSafeQuit() (isPossible bool) {
	return !mutated
}

func Add(add rune) (requiredUpdates int) {
	buffer[cursorIndex] = add
	length++
	cursorIndex++
	if length > int(float64(capacity)*resizeAt) {
		capacity = int(float64(capacity) * resizeRatio)
		temp := make([]rune, capacity)
		backLength := length - cursorIndex
		newBackBlockIndex := capacity - backLength
		copy(temp[:cursorIndex], buffer[:cursorIndex])
		copy(temp[newBackBlockIndex:], buffer[backBlockIndex:])
		backBlockIndex = newBackBlockIndex
		buffer = temp
	}
	mutated = true
	return computeRequiredUpdates()
}

func Remove() (possible bool, newX, requiredUpdates int) {
	if cursorIndex == 0 {
		return false, 0, 0
	}
	length--
	cursorIndex--
	mutated = true
	return true, computeNewX(), computeRequiredUpdates()
}

func computeNewX() (newX int) {
	newX = 0
	for i := cursorIndex - 1; i >= 0 && buffer[i] != '\n'; i-- {
		newX++
	}
	return newX
}

func computeRequiredUpdates() (requiredUpdates int) {
	for i := backBlockIndex; i < capacity && buffer[i] != '\n'; i++ {
		requiredUpdates++
	}
	return requiredUpdates
}

func RemoveCurrent() (possible, xBack bool, requiredUpdates int) {
	if backBlockIndex == length || buffer[backBlockIndex] == '\n' {
		return false, false, 0
	}
	length--
	backBlockIndex++
	mutated = true
	return true, buffer[backBlockIndex] == '\n', computeRequiredUpdates() + 1
}

func RemoveLine() (yBack, isEmpty bool) {
	for i := cursorIndex - 1; i >= 0 && buffer[i] != '\n'; i-- {
		cursorIndex--
		length--
	}
	for i := backBlockIndex; buffer[i] != '\n'; i++ {
		backBlockIndex++
		length--
	}
	if length > 1 {
		backBlockIndex++
		length--
	}
	mutated = true
	return backBlockIndex == capacity, length == 1
}

func RemoveRestOfLine() (requiredUpdates int) {
	requiredUpdates = 0
	for i := backBlockIndex; buffer[i] != '\n'; i++ {
		requiredUpdates++
		backBlockIndex++
		length--
	}
	mutated = true
	return requiredUpdates
}

func Left() (possible bool) {
	if cursorIndex == 0 || buffer[cursorIndex-1] == '\n' {
		return false
	}
	cursorIndex--
	backBlockIndex--
	buffer[backBlockIndex] = buffer[cursorIndex]
	return true
}

func Right(isInsert bool) (possible bool) {
	offset := computeOffset(isInsert)
	if backBlockIndex == capacity-offset || buffer[backBlockIndex] == '\n' ||
		buffer[backBlockIndex+offset] == '\n' {
		return false
	}
	buffer[cursorIndex] = buffer[backBlockIndex]
	cursorIndex++
	backBlockIndex++
	return true
}

func Up(oldX int, isInsert bool) (possible bool, newX int) {
	if cursorIndex == 0 {
		return false, oldX
	}
	i := cursorIndex - 1
	count := 0
	for i > 0 && buffer[i] != '\n' {
		i--
		count++
	}
	if i == 0 && buffer[i] != '\n' {
		return false, oldX
	}
	i--
	count++
	if i == -1 || buffer[i] == '\n' {
		temp := buffer[i+1 : cursorIndex]
		copy(buffer[backBlockIndex-count:backBlockIndex], temp)
		cursorIndex = i + 1
		backBlockIndex -= count
		return true, 0
	}
	lineLen := 0
	for j := i; j > 0 && buffer[j-1] != '\n'; j-- {
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
	temp := buffer[i+1 : cursorIndex]
	copy(buffer[backBlockIndex-count:backBlockIndex], temp)
	cursorIndex = i + 1
	backBlockIndex -= count
	return true, newX
}

func Down(oldX int, isInsert bool) (possible bool, newX int) {
	offset := computeOffset(isInsert)
	i := backBlockIndex
	for i < capacity && buffer[i] != '\n' {
		i++
	}
	i++
	if i >= capacity {
		return false, oldX
	}
	if buffer[i] != '\n' {
		for newX = 0; newX < oldX && i < capacity-offset && buffer[i+offset] != '\n'; newX++ {
			i++
		}
	}
	temp := buffer[backBlockIndex:i]
	delta := i - backBlockIndex
	copy(buffer[cursorIndex:cursorIndex+delta], temp)
	cursorIndex += delta
	backBlockIndex = i
	return true, newX
}

func computeOffset(isInsert bool) (offset int) {
	if isInsert {
		return 0
	}
	return 1
}

func GetBottom(currentY, getY int) (bottom string) {
	var i int
	deltaY := getY - currentY
	for i = backBlockIndex; i < capacity; i++ {
		if buffer[i] == '\n' {
			deltaY--
			if deltaY == 0 {
				break
			}
		}
	}
	var sb strings.Builder
	for j := i + 1; j < capacity && buffer[j] != '\n'; j++ {
		sb.WriteRune(buffer[j])
	}
	if sb.Len() == 0 {
		return "~"
	}
	return sb.String()
}

func GetLine() (previous string) {
	var sb strings.Builder
	startIndex := cursorIndex - 1
	for i := startIndex; i >= 0 && buffer[i] != '\n'; i-- {
		startIndex--
	}
	for i := startIndex + 1; i < cursorIndex; i++ {
		sb.WriteRune(buffer[i])
	}
	for i := backBlockIndex; i < capacity && buffer[i] != '\n'; i++ {
		sb.WriteRune(buffer[i])
	}
	return sb.String()
}
