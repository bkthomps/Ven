package buffer

import (
	"io/ioutil"
	"os"
)

const TabSize = 4

type File struct {
	fileName string
	mutated  bool

	First *Line
	last  *Line

	Current *Line
	offset  int
	lines   int
}

func (file *File) Init(fileName string) {
	file.fileName = fileName
	file.mutated = false
	line := &Line{}
	line.Init(nil, nil)
	file.First = line
	file.last = line
	file.Current = line
	file.offset = 0
	file.lines = 1
	arr := readFile(fileName)
	for _, character := range arr {
		file.Add(character)
	}
	file.Current = file.First
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
