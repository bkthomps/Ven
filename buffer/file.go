package buffer

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"

	"github.com/mattn/go-runewidth"
)

const TabSize = 8

type File struct {
	fileName string
	mutated  bool

	First *Line
	last  *Line

	Current *Line
	lines   int

	runeOffset    int
	spacingOffset int
}

func (file *File) Init(fileName string) {
	file.fileName = fileName
	line := &Line{}
	line.Init(nil, nil)
	file.First = line
	file.last = line
	file.Current = line
	file.lines = 1
	arr := readFile(fileName)
	for _, character := range arr[:len(arr)-1] {
		file.Add(character)
	}
	file.Current = file.First
	file.runeOffset = 0
	file.spacingOffset = 0
	file.mutated = false
}

func readFile(fileName string) (arr []rune) {
	dat, err := ioutil.ReadFile(fileName)
	data := ""
	if !os.IsNotExist(err) {
		data = string(dat)
	}
	arr = []rune(data)
	if len(arr) == 0 || arr[len(arr)-1] != '\n' {
		arr = append(arr, '\n')
	}
	return arr
}

func (file *File) CanSafeQuit() bool {
	return !file.mutated
}

func (file *File) Save() error {
	osFile, err := os.Create(file.fileName)
	if err != nil {
		return err
	}
	arr := make([]rune, 0)
	for traverse := file.First; traverse != nil; traverse = traverse.Next {
		arr = append(arr, traverse.Data...)
		arr = append(arr, '\n')
	}
	_, err = fmt.Fprintf(osFile, string(arr))
	if err != nil {
		return err
	}
	_ = osFile.Close()
	file.mutated = false
	return nil
}

func (file *File) runeWidthIncrease(r rune) int {
	return runeWidthJump(r, file.spacingOffset)
}

func (file *File) runeWidthDecrease(r rune) int {
	if r == '\t' {
		offset := 0
		for i := 0; i < file.runeOffset; i++ {
			offset = runeWidthJump(file.Current.Data[i], offset)
		}
		return offset
	}
	return file.spacingOffset - runewidth.RuneWidth(r)
}

func runeWidthJump(r rune, offset int) int {
	if r == '\t' {
		return int(math.Ceil(float64(offset+1)/float64(TabSize)) * TabSize)
	}
	return offset + runewidth.RuneWidth(r)
}
