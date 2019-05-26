[![Go Report Card](https://goreportcard.com/badge/github.com/bkthomps/Ven)](https://goreportcard.com/report/github.com/bkthomps/Ven)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/bkthomps/Ven)

# Ven
Vi Enhanced. A text editor which is an enhancement to vi, and is similar to vim, but written in Go. Uses a gap buffer to store the text.

## Installation
* Download golang if you do not yet have it
* Get Ven:
  ```
  go get github.com/bkthomps/Ven
  ```
* Go to the Ven directory:
  ```
  cd ~/go/src/github.com/bkthomps/Ven
  ```
* Build Ven:
  ```
  go build
  ```
* Add this to your bashrc:
  ```
  alias ven='~/go/src/github.com/bkthomps/Ven/Ven'
  ```
* Now you can run Ven anywhere with `ven <filename>`

## Commands
There are three modes: normal mode, command mode, and insertion mode.

### Normal Mode
* `i` to go into insertion mode
* `:` to go into command mode
* `/` to go into command (search) mode
* `x` delete character under the cursor
* `X` delete character before the cursor
* `dd` delete entire line
* `D` delete rest of line

### Command Mode
* `esc` to go into normal mode
* `/<search>` to search for a string (currently doesn't support regex)
* `:w` to save the file
* `:wq` to save and quit
* `:q` to safely quit
* `:q!` to force quit without saving

### Insertion Mode
* `esc` to go into normal mode
* any character press gets inserted
