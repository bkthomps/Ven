[![GitHubBuild](https://github.com/bkthomps/Ven/workflows/build/badge.svg)](https://github.com/bkthomps/Ven)
[![Go Report Card](https://goreportcard.com/badge/github.com/bkthomps/Ven)](https://goreportcard.com/report/github.com/bkthomps/Ven)
[![codecov](https://codecov.io/gh/bkthomps/Ven/branch/main/graph/badge.svg)](https://codecov.io/gh/bkthomps/Ven)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/bkthomps/Ven?tab=overview)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/bkthomps/Ven/blob/main/LICENSE)

# Ven
Vi Enhanced. A text editor which is an enhancement to vi, and is similar to vim, but written in Go.

## Installation
* Run: `go install github.com/bkthomps/Ven@latest`
* Then add this to your bashrc or zshrc: `alias ven='~/go/bin/Ven'`
* You can now run Ven from anywhere using `ven` or `ven <filename>`

## Commands
There are three modes: normal mode, command mode, and insertion mode.

### Normal Mode
* `:` to go into command mode
* `/` to go into command (search) mode
* `i` to go into insertion mode at the cursor
* `a` to go into insertion mode after the cursor
* `A` to go into insertion mode at the end of the line
* `I` to go into insertion mode at the beginning of the line
* `o` to open a new line under the cursor and go into insertion mode
* `O` to open a new line above the cursor and go into insertion mode
* `j` or down arrow to go down
* `k` or up arrow to go up
* `h` or left arrow to go left
* `l` or right arrow to go right
* `H` to move the cursor to the top of the screen
* `M` to move the cursor to the middle of the screen
* `L` to move the cursor to the bottom of the screen
* `0` to move the cursor to the start of the line
* `$` to move the cursor to the end of the line
* `gg` to move the cursor to the start of the file
* `G` to move the cursor to the end of the file
* `w` to move the cursor to the start of the next word
* `b` to move the cursor to the start of the current word
* `e` to move the cursor to the end of the current word
* `ctrl-f` to go a page forward
* `ctrl-b` to go a page backward
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
