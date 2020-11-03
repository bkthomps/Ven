package buffer

import (
	"fmt"
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
	lines   int

	runeOffset    int
	spacingOffset int
}

func (file *File) Init(fileName string) {
	file.fileName = fileName
	file.mutated = false
	line := &Line{}
	line.Init(nil, nil)
	file.First = line
	file.last = line
	file.Current = line
	file.lines = 1
	file.runeOffset = 0
	file.spacingOffset = 0
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
