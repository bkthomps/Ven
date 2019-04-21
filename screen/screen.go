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

package screen

import (
	"github.com/bkthomps/Ven/buffer"
	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
	"log"
)

const errorCommand = "-- Invalid Command --"

var terminalStyle = tcell.StyleDefault.Foreground(tcell.ColorBlack)
var cursorStyle = terminalStyle.Background(tcell.ColorDarkGray)
var screen tcell.Screen

var screenHeight = 0
var screenWidth = 0

var xCursor = 0
var yCursor = 0
var xCommandCursor = 0

var blankLine = ""
var command = ""

const (
	InsertMode = iota
	NormalMode
	CommandMode
	CommandErrorMode
)

var mode = NormalMode

func Init(s tcell.Screen, quit chan struct{}) {
	screen = s
	if e := screen.Init(); e != nil {
		log.Fatal(e)
	}
	buffer.Init() // TODO: should set initial screen
	screen.SetStyle(terminalStyle)
	screen.Show()
	updateProperties()
	screen.SetContent(xCursor, yCursor, ' ', nil, cursorStyle)
	displayMode()
	go listener(quit)
}

func updateProperties() {
	x, y := screen.Size()
	screenWidth, screenHeight = x, y
	for i := 0; i < x; i++ {
		blankLine += " "
	}
}

func displayMode() {
	switch mode {
	case InsertMode:
		screen.SetContent(xCursor, yCursor, ' ', nil, cursorStyle)
		putCommand("-- INSERT --")
	case NormalMode:
		screen.SetContent(xCursor, yCursor, ' ', nil, cursorStyle)
		putCommand(blankLine)
	case CommandMode:
		screen.SetContent(xCursor, yCursor, ' ', nil, terminalStyle)
		putCommand(blankLine)
		putCommand(command)
		screen.SetContent(xCommandCursor, screenHeight-1, ' ', nil, cursorStyle)
	case CommandErrorMode:
		screen.SetContent(xCommandCursor, screenHeight-1, ' ', nil, terminalStyle)
		putCommand(blankLine)
		putCommand(errorCommand)
	}
	screen.Sync()
}

func putCommand(str string) {
	puts(screen, terminalStyle, 0, screenHeight-1, str)
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
			case CommandErrorMode:
				mode = CommandMode
			}
			displayMode()
		case *tcell.EventResize:
			updateProperties()
			displayMode()
		}
	}
}

func insertMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEsc:
		mode = NormalMode
	default:
		bufferAction(ev)
	}
}

func normalMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyDown, tcell.KeyUp, tcell.KeyLeft, tcell.KeyRight:
		bufferAction(ev)
	default:
		switch ev.Rune() {
		case 'i':
			mode = InsertMode
		case ':':
			mode = CommandMode
			command = ":"
			xCommandCursor = 1
		}
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
		xCommandCursor--
	case tcell.KeyDown, tcell.KeyUp:
		// Do Nothing
	case tcell.KeyLeft:
		if xCommandCursor > 1 {
			xCommandCursor--
		}
	case tcell.KeyRight:
		if xCommandCursor < len(command) {
			xCommandCursor++
		}
	default:
		command += string(ev.Rune())
		xCommandCursor++
	}
}

func bufferAction(ev *tcell.EventKey) {
	screen.SetContent(xCursor, yCursor, ' ', nil, terminalStyle)
	// TODO: add cases for scrolling
	switch ev.Key() {
	case tcell.KeyDown:
		possible := buffer.Down()
		if possible {
			yCursor++
		}
	case tcell.KeyUp:
		possible := buffer.Up()
		if possible {
			yCursor--
		}
	case tcell.KeyLeft:
		possible := buffer.Left()
		if possible {
			xCursor--
		}
	case tcell.KeyRight:
		possible := buffer.Right()
		if possible {
			xCursor++
		}
	case tcell.KeyDEL:
		possible := buffer.Remove()
		putRune(' ')
		if possible {
			if ev.Rune() != '\n' {
				xCursor--
			} else {
				// TODO: xCursor goes to end, count lines
				yCursor--
				// TODO: shift everything back
			}
		}
	default:
		buffer.Add(ev.Rune())
		putRune(ev.Rune())
		if ev.Rune() != '\n' {
			xCursor++
		} else {
			xCursor = 0
			yCursor++
			// TODO: shift everything over
		}
	}
	screen.SetContent(xCursor, yCursor, ' ', nil, cursorStyle)
}

func executeCommand(quit chan struct{}) {
	switch command {
	case ":q!":
		close(quit)
	default:
		mode = CommandErrorMode
		displayMode()
	}
}

func putRune(r rune) {
	puts(screen, terminalStyle, xCursor, yCursor, string(r))
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
