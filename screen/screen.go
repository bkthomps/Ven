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
const errorSave = "-- Could Not Save File --"

const (
	InsertMode = iota
	NormalMode
	CommandMode
	CommandErrorMode
)

var mode = NormalMode

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

func Init(s tcell.Screen, quit chan struct{}) {
	screen = s
	if e := screen.Init(); e != nil {
		log.Fatal(e)
	}
	arr := buffer.Init("file.txt") // TODO: should set initial screen
	screen.SetStyle(terminalStyle)
	screen.Show()
	updateProperties()
	setInitial(arr)
	setColor(xCursor, yCursor, cursorStyle)
	displayMode()
	go listener(quit)
}

func setInitial(arr []rune) {
	x, y := 0, 0
	for i := 0; i < len(arr) && y != screenHeight; i++ {
		cur := arr[i]
		if cur == '\n' {
			y++
			x = 0
		} else if x < screenWidth {
			putRune(cur, x, y)
			x++
		}
	}
	for i := y; i != screenHeight-1; i++ {
		putRune('~', 0, i)
	}
}

func updateProperties() {
	x, y := screen.Size()
	screenWidth, screenHeight = x, y
	for i := 0; i < x; i++ {
		blankLine += " "
	}
}

func setColor(x, y int, s tcell.Style) {
	r, _, _, _ := screen.GetContent(x, y)
	screen.SetContent(x, y, r, nil, s)
}

func displayError(error string) {
	putCommand(blankLine)
	putCommand(error)
	mode = CommandErrorMode
	displayMode()
}

func displayMode() {
	switch mode {
	case InsertMode:
		setColor(xCursor, yCursor, cursorStyle)
		putCommand("-- INSERT --")
	case NormalMode:
		setColor(xCursor, yCursor, cursorStyle)
		putCommand(blankLine)
	case CommandMode:
		setColor(xCursor, yCursor, terminalStyle)
		putCommand(blankLine)
		putCommand(command)
		setColor(xCommandCursor, screenHeight-1, cursorStyle)
	case CommandErrorMode:
		setColor(xCommandCursor, screenHeight-1, terminalStyle)
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
	setColor(xCursor, yCursor, terminalStyle)
	// TODO: add cases for scrolling
	switch ev.Key() {
	case tcell.KeyDown:
		if yCursor == screenHeight-2 {
			// TODO: down scrolling
		} else {
			possible, x := buffer.Down(xCursor)
			if possible {
				yCursor++
				xCursor = x
			}
		}
	case tcell.KeyUp:
		if yCursor == 0 {
			// TODO: up scrolling
		} else {
			possible, x := buffer.Up(xCursor)
			if possible {
				yCursor--
				xCursor = x
			}
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
		possible, requiredUpdates := buffer.Remove()
		if possible {
			if ev.Rune() != '\n' {
				xCursor--
				shiftLeft(requiredUpdates)
			} else {
				// TODO: xCursor goes to end, count lines
				yCursor--
				// TODO: shift all lines back
			}
		}
	default:
		x, y := xCursor, yCursor
		requiredUpdates := buffer.Add(ev.Rune())
		if ev.Rune() != '\n' {
			xCursor++
			shiftRight(requiredUpdates)
		} else {
			xCursor = 0
			yCursor++
			// TODO: shift all lines over
		}
		putRune(ev.Rune(), x, y)
	}
	setColor(xCursor, yCursor, cursorStyle)
}

func shiftLeft(requiredUpdates int) {
	for i := xCursor; i <= xCursor+requiredUpdates; i++ {
		r, _, _, _ := screen.GetContent(i+1, yCursor)
		screen.SetContent(i, yCursor, r, nil, terminalStyle)
	}
}

func shiftRight(requiredUpdates int) {
	for i := xCursor + requiredUpdates; i >= xCursor; i-- {
		r, _, _, _ := screen.GetContent(i-1, yCursor)
		screen.SetContent(i, yCursor, r, nil, terminalStyle)
	}
}

func executeCommand(quit chan struct{}) {
	switch command {
	case ":q!":
		close(quit)
	case ":w":
		// TODO: implement the version with file name also
		write()
	case ":wq":
		// TODO: implement the version with file name also
		saved := write()
		if saved {
			close(quit)
		}
	default:
		displayError(errorCommand)
	}
}

func write() (saved bool) {
	err := buffer.Save()
	if err != nil {
		displayError(errorSave)
		return false
	}
	mode = NormalMode
	return true
}

func putRune(r rune, x, y int) {
	puts(screen, terminalStyle, x, y, string(r))
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
