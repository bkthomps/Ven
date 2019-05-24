[![Go Report Card](https://goreportcard.com/badge/github.com/bkthomps/Ven)](https://goreportcard.com/report/github.com/bkthomps/Ven)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/bkthomps/Ven)

# Ven
Vi Enhanced. An enhancement to vi which is similar to vim, but written in Go. Uses a gap buffer to store data.

## Installation
```
go get github.com/bkthomps/Ven
```

## Commands
There are three modes: normal mode, command mode, and insertion mode.

### Normal Mode
* `i` to go into insertion mode
* `:` to go into command mode

### Command Mode
* `esc` to go into normal mode
* `:w` to save the file
* `:wq` to save and quit
* `:q!` to force quit without saving

### Insertion Mode
* `esc` to go into normal mode
* any character press gets inserted
