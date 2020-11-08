[![GitHubBuild](https://github.com/bkthomps/Ven/workflows/build/badge.svg)](https://github.com/bkthomps/Ven)
[![Go Report Card](https://goreportcard.com/badge/github.com/bkthomps/Ven)](https://goreportcard.com/report/github.com/bkthomps/Ven)
[![codecov](https://codecov.io/gh/bkthomps/Ven/branch/master/graph/badge.svg)](https://codecov.io/gh/bkthomps/Ven)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/bkthomps/Ven?tab=overview)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/bkthomps/Ven/blob/master/LICENSE)

# Ven
Vi Enhanced. A text editor which is an enhancement to vi, and is similar to vim, but written in Go. Uses a gap buffer to store the text.

## Installation
* Download golang if you have not yet done so
* Run:
  ```
  go get github.com/bkthomps/Ven
  go build $GOPATH/src/github.com/bkthomps/Ven
  ```
* Then add this to your bashrc:
  ```
  alias ven='$GOPATH/src/github.com/bkthomps/Ven/Ven'
  ```
* You can now run Ven from anywhere using `ven <filename>`

## Commands
There are three modes: normal mode, command mode, and insertion mode.

### Normal Mode
* `i` to go into insertion mode
* `:` to go into command mode
* `/` to go into command (search) mode
* `j` or down arrow to go down
* `k` or up arrow to go up
* `h` or left arrow to go left
* `l` or right arrow to go right
* `H` to move the cursor to the top of the screen
* `M` to move the cursor to the middle of the screen
* `L` to move the cursor to the bottom of the screen
* `x` delete character under the cursor
* `X` delete character before the cursor
* `dd` delete entire line
* `D` delete rest of line

### Command Mode
* `esc` to go into normal mode
* `/<search>` to search for a string (supports regex)
* `:w` to save the file
* `:wq` to save and quit
* `:q` to safely quit
* `:q!` to force quit without saving

### Insertion Mode
* `esc` to go into normal mode
* any character press gets inserted
