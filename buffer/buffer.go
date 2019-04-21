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

import "math"

// TODO: buffer should increase by x1.5 when more than x0.75 of space is taken
// TODO: for pasting, check space before pasting
// TODO: if buffer goes to less than x0.25, decrease size by x2 ? (not sure on numbers)
// TODO: when opening a file, make the initial size = file_sz * 1.5
// TODO: the cursor is equivalent to the blank space in this memory model
// TODO: if no file exists, start size at 16

const resizeRatio = 1.5
const resizeAt = 0.75
const minimumSize = 16

var buffer []rune
var capacity = 0
var length = 0
var cursorIndex = 0    // The empty byte at the start of the gap
var backBlockIndex = 0 // The byte after the end of the gap

func Init() {
	// TODO: Do some file io, get the file and its size
	fileSize := 0
	length = fileSize
	capacity = int(float64(length) * resizeRatio)
	backBlockIndex = capacity
	buffer = make([]rune, int(math.Max(float64(capacity), minimumSize)))
	// TODO: copy the file to the back of the block since cursor starts at 0
}

func Add(add rune) {
	buffer[cursorIndex] = add
	length++
	cursorIndex++
	if length > int(float64(capacity)*resizeAt) {
		// TODO: implement resizing
	}
}

func Remove() (possible bool) {
	if cursorIndex == 0 {
		return false
	}
	length--
	cursorIndex--
	return true
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

func Right() (possible bool) {
	if backBlockIndex == capacity || buffer[backBlockIndex] == '\n' {
		return false
	}
	buffer[cursorIndex] = buffer[backBlockIndex]
	cursorIndex++
	backBlockIndex++
	return true
}

func Up() (possible bool) {
	// TODO: implement
	return false
}

func Down() (possible bool) {
	// TODO: implement
	return false
}
