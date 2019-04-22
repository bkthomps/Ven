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

func Init(name string) (arr []rune) {
	fileName = name
	dat, err := ioutil.ReadFile(name)
	data := ""
	if !os.IsNotExist(err) {
		data = string(dat)
	}
	arr = []rune(data)
	length = len(arr)
	capacity = int(math.Max(float64(length)*resizeRatio, minimumSize))
	backBlockIndex = capacity - length
	buffer = make([]rune, capacity)
	copy(buffer[backBlockIndex:], arr)
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
	fmt.Fprintf(file, string(arr))
	return nil
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
	return computeRequiredUpdates()
}

func Remove() (possible bool, requiredUpdates int) {
	if cursorIndex == 0 {
		return false, 0
	}
	length--
	cursorIndex--
	return true, computeRequiredUpdates()
}

func computeRequiredUpdates() (requiredUpdates int) {
	for i := backBlockIndex; i < capacity && buffer[i] != '\n'; i++ {
		requiredUpdates++
	}
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
