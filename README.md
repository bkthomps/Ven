[![Go Report Card](https://goreportcard.com/badge/github.com/bkthomps/Ven)](https://goreportcard.com/report/github.com/bkthomps/Ven)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/bkthomps/Ven)

# Ven
Vi Enhanced. A text editor which is an enhancement to vi, and is similar to vim, but written in Go. Uses a gap buffer to store the text.

## Installation
```
go get github.com/bkthomps/Ven
```

## Commands
There are three modes: normal mode, command mode, and insertion mode.

### Normal Mode
* `i` to go into insertion mode
* `:` to go into command mode
* `x` delete character under the cursor
* `X` delete character before the cursor
* `dd` delete entire line
* `D` delete rest of line

### Command Mode
* `esc` to go into normal mode
* `:w` to save the file
* `:wq` to save and quit
* `:q` to safely quit
* `:q!` to force quit without saving

### Insertion Mode
* `esc` to go into normal mode
* any character press gets inserted
