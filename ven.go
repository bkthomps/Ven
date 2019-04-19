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

package main

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/mattn/go-runewidth"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
	"log"
)

var terminalStyle = tcell.StyleDefault.Foreground(tcell.ColorBlack)
var screen tcell.Screen

const (
	InsertMode = iota
	NormalMode
	CommandMode
	CommandErrorMode
)

var mode = NormalMode

var sidebar = 0
var screenHeight = 0
var screenWidth = 0
var blankLine = ""
var command = ""

func putln(vertical int, str string) {
	puts(screen, terminalStyle, sidebar, vertical, str)
}

// This function is from: https://github.com/gdamore/tcell/blob/master/_demos/unicode.go
func puts(s tcell.Screen, style tcell.Style, x, y int, str string) {
	i := 0
	var deferred []rune
	dwidth := 0
	zwj := false
	for _, r := range str {
		if r == '\u200d' {
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
			deferred = append(deferred, r)
			zwj = true
			continue
		}
		if zwj {
			deferred = append(deferred, r)
			zwj = false
			continue
		}
		switch runewidth.RuneWidth(r) {
		case 0:
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
		case 1:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 1
		case 2:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 2
		}
		deferred = append(deferred, r)
	}
	if len(deferred) != 0 {
		s.SetContent(x+i, y, deferred[0], deferred[1:], style)
		i += dwidth
	}
}

func displayMode() {
	switch mode {
	case InsertMode:
		putln(screenHeight-1, "-- INSERT --")
	case NormalMode:
		putln(screenHeight-1, blankLine)
	case CommandMode:
		putln(screenHeight-1, blankLine)
		putln(screenHeight-1, command)
	}
	screen.Sync()
}

func insertMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEsc:
		mode = NormalMode
	}
	displayMode()
}

func normalMode(ev *tcell.EventKey) {
	switch ev.Rune() {
	case 'i':
		mode = InsertMode
	case ':':
		mode = CommandMode
		command = ":"
	}
	displayMode()
}

func executeCommand(quit chan struct{}) {
	switch command {
	case ":q!":
		close(quit)
	}
}

func commandMode(ev *tcell.EventKey, quit chan struct{}) {
	switch ev.Key() {
	case tcell.KeyEsc:
		mode = NormalMode
	case tcell.KeyEnter:
		executeCommand(quit)
	case tcell.KeyDEL:
		sz := len(command)
		command = command[:sz-1]
		if sz == 1 {
			mode = NormalMode
		}
	default:
		command += string(ev.Rune())
	}
	displayMode()
}

func computeScreenProperties() {
	y, _ := terminal.Height()
	screenHeight = int(y)
	x, _ := terminal.Width()
	screenWidth = int(x)
	for i := 0; i < screenWidth; i++ {
		blankLine += " "
	}
}

func listener(quit chan struct{}) {
	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch mode {
			case InsertMode:
				insertMode(ev)
			case NormalMode:
				normalMode(ev)
			case CommandMode:
				commandMode(ev, quit)
			}
		case *tcell.EventResize:
			computeScreenProperties()
			displayMode()
		}
	}
}

func main() {
	s, e := tcell.NewScreen()
	screen = s
	if e != nil {
		log.Fatal(e)
	}
	encoding.Register()
	if e = screen.Init(); e != nil {
		log.Fatal(e)
	}
	screen.SetStyle(terminalStyle)
	quit := make(chan struct{})
	screen.Show()
	computeScreenProperties()
	displayMode()
	go listener(quit)
	<-quit
	screen.Fini()
}
