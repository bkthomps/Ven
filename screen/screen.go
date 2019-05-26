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

const insertMessage = "-- INSERT --"
const errorCommand = "-- Invalid Command --"
const errorSave = "-- Could Not Save File --"
const modifiedFile = "-- File Has Been Modified Since Last Save --"

const (
	insertMode = iota
	normalMode
	commandMode
	commandErrorMode
)

var mode = normalMode

var terminalStyle = tcell.StyleDefault.Foreground(tcell.ColorBlack)
var cursorStyle = terminalStyle.Background(tcell.ColorDarkGray)
var highlightStyle = terminalStyle.Background(tcell.ColorYellow)
var screen tcell.Screen

var screenHeight = 0
var screenWidth = 0

var xCursor = 0
var yCursor = 0
var xCommandCursor = 0

var xSearchPoints = []int(nil)
var ySearchPoints = []int(nil)
var searchStringLength = 0

var oldCommand = '_'

var blankLine = ""
var command = ""

func Init(s tcell.Screen, quit chan struct{}, fileName string) {
	screen = s
	if e := screen.Init(); e != nil {
		log.Fatal(e)
	}
	arr := buffer.Init(fileName)
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
	for i := 0; i < len(arr) && y < screenHeight-1; i++ {
		cur := arr[i]
		if cur == '\n' {
			y++
			x = 0
		} else if x < screenWidth {
			putRune(cur, x, y)
			x++
		}
	}
	for i := y; i < screenHeight-1; i++ {
		for j := 2; j < screenWidth; j++ {
			r1, _, _, _ := screen.GetContent(j-2, i)
			r2, _, _, _ := screen.GetContent(j-1, i)
			if r1 == ' ' && r2 == ' ' {
				break
			}
			screen.SetContent(j-2, i, ' ', nil, terminalStyle)
		}
		putRune('~', 0, i)
	}
}

func updateProperties() (isBigger bool) {
	oldWidth, oldHeight := screenWidth, screenHeight
	x, y := screen.Size()
	screenWidth, screenHeight = x, y
	for i := 0; i < x; i++ {
		blankLine += " "
	}
	return oldWidth < screenWidth || oldHeight < screenHeight
}

func setColor(x, y int, s tcell.Style) {
	r, _, _, _ := screen.GetContent(x, y)
	screen.SetContent(x, y, r, nil, s)
}

func displayError(error string) {
	putCommand(blankLine)
	putCommand(error)
	mode = commandErrorMode
	displayMode()
}

func displayMode() {
	switch mode {
	case insertMode:
		setColor(xCursor, yCursor, cursorStyle)
		putCommand(blankLine)
		putCommand(insertMessage)
	case normalMode:
		setColor(xCursor, yCursor, cursorStyle)
		putCommand(blankLine)
	case commandMode:
		putCommand(blankLine)
		putCommand(command)
		setColor(xCommandCursor, screenHeight-1, cursorStyle)
	case commandErrorMode:
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
			case insertMode:
				executeInsertMode(ev)
			case normalMode:
				executeNormalMode(ev)
			case commandMode:
				executeCommandMode(ev, quit)
			case commandErrorMode:
				mode = commandMode
			}
			displayMode()
		case *tcell.EventResize:
			isBigger := updateProperties()
			for xCursor >= screenWidth {
				actionLeft()
			}
			for yCursor >= screenHeight-1 {
				shiftUp(-1)
				yCursor--
			}
			if isBigger {
				arr := buffer.Redraw(yCursor, screenHeight)
				setInitial(arr)
			}
			displayMode()
		}
	}
}

func executeInsertMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEsc:
		setColor(xCursor, yCursor, terminalStyle)
		mode = normalMode
		possible := buffer.Left()
		if possible {
			xCursor--
		}
	default:
		bufferAction(ev)
	}
}

func executeNormalMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyDown, tcell.KeyUp, tcell.KeyLeft, tcell.KeyRight:
		bufferAction(ev)
	default:
		switch ev.Rune() {
		case 'i':
			mode = insertMode
		case ':', '/':
			setColor(xCursor, yCursor, terminalStyle)
			mode = commandMode
			command = string(ev.Rune())
			xCommandCursor = len(command)
		case 'x':
			isPossible, xBack, requiredUpdates := buffer.RemoveCurrent()
			if isPossible {
				shiftLeft(requiredUpdates)
				if xBack {
					actionLeft()
				}
			}
		case 'X':
			if xCursor != 0 {
				actionLeft()
				isPossible, _, requiredUpdates := buffer.RemoveCurrent()
				if isPossible {
					shiftLeft(requiredUpdates)
				}
			}
		case 'd':
			if oldCommand == 'd' {
				xCursor = 0
				shiftUp(yCursor - 1)
				yBack, isEmpty := buffer.RemoveLine()
				if yBack {
					actionUp()
				}
				if isEmpty {
					screen.SetContent(0, 0, ' ', nil, terminalStyle)
				}
				oldCommand = '_'
			} else {
				oldCommand = ev.Rune()
			}
		case 'D':
			requiredUpdates := buffer.RemoveRestOfLine()
			for i := xCursor; i <= xCursor+requiredUpdates; i++ {
				screen.SetContent(i, yCursor, ' ', nil, terminalStyle)
			}
		}
	}
}

func executeCommandMode(ev *tcell.EventKey, quit chan struct{}) {
	switch ev.Key() {
	case tcell.KeyEsc:
		removeHighlighting()
		mode = normalMode
	case tcell.KeyEnter:
		executeCommand(quit)
	case tcell.KeyDEL:
		if xCommandCursor <= 1 && len(command) > 1 {
			break
		}
		runeCopy := []rune(command)
		copy(runeCopy[xCommandCursor-1:], runeCopy[xCommandCursor:])
		shrinkSize := len(runeCopy) - 1
		runeCopy = runeCopy[:shrinkSize]
		command = string(runeCopy)
		if shrinkSize == 0 {
			removeHighlighting()
			mode = normalMode
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

func removeHighlighting() {
	if xSearchPoints != nil {
		for i := 0; i < len(xSearchPoints); i++ {
			startX, y := xSearchPoints[i], ySearchPoints[i]
			for x := startX; x < searchStringLength+startX; x++ {
				r, _, _, _ := screen.GetContent(x, y)
				screen.SetContent(x, y, r, nil, terminalStyle)
			}
		}
		xSearchPoints = nil
		ySearchPoints = nil
		searchStringLength = 0
	}
}

func bufferAction(ev *tcell.EventKey) {
	setColor(xCursor, yCursor, terminalStyle)
	switch ev.Key() {
	case tcell.KeyDown:
		actionDown()
	case tcell.KeyUp:
		actionUp()
	case tcell.KeyLeft:
		actionLeft()
	case tcell.KeyRight:
		actionRight()
	case tcell.KeyDEL:
		actionDelete()
	case tcell.KeyEnter:
		actionEnter()
	default:
		actionKeyPress(ev)
	}
	setColor(xCursor, yCursor, cursorStyle)
}

func actionDown() {
	possible, x := buffer.Down(xCursor, mode == insertMode)
	if possible {
		if yCursor == screenHeight-2 {
			for y := 0; y < screenHeight-2; y++ {
				for x := 0; x < screenWidth; x++ {
					r, _, _, _ := screen.GetContent(x, y+1)
					screen.SetContent(x, y, r, nil, terminalStyle)
				}
			}
			putString(blankLine, 0, screenHeight-2)
			putString(buffer.GetLine(), 0, screenHeight-2)
		} else {
			yCursor++
		}
		xCursor = x
	}
}

func actionUp() {
	possible, x := buffer.Up(xCursor, mode == insertMode)
	if possible {
		if yCursor == 0 {
			for y := screenHeight - 2; y > 0; y-- {
				for x := 0; x < screenWidth; x++ {
					r, _, _, _ := screen.GetContent(x, y-1)
					screen.SetContent(x, y, r, nil, terminalStyle)
				}
			}
			putString(blankLine, 0, 0)
			putString(buffer.GetLine(), 0, 0)
		} else {
			yCursor--
		}
		xCursor = x
	}
}

func actionLeft() {
	possible := buffer.Left()
	if possible {
		xCursor--
	}
}

func actionRight() {
	possible := buffer.Right(mode == insertMode)
	if possible {
		xCursor++
	}
}

func actionDelete() {
	possible, newX, requiredUpdates := buffer.Remove()
	if possible {
		if xCursor != 0 {
			xCursor--
			shiftLeft(requiredUpdates)
		} else {
			xCursor = newX
			yCursor--
			for x := 0; x < requiredUpdates; x++ {
				r, _, _, _ := screen.GetContent(x, yCursor+1)
				screen.SetContent(x+newX, yCursor, r, nil, terminalStyle)
			}
			shiftUp(yCursor)
		}
	}
}

func shiftLeft(requiredUpdates int) {
	for i := xCursor; i <= xCursor+requiredUpdates; i++ {
		r, _, _, _ := screen.GetContent(i+1, yCursor)
		screen.SetContent(i, yCursor, r, nil, terminalStyle)
	}
}

func shiftUp(ontoY int) {
	for y := ontoY + 1; y < screenHeight-2; y++ {
		for x := 0; x < screenWidth; x++ {
			r, _, _, _ := screen.GetContent(x, y+1)
			screen.SetContent(x, y, r, nil, terminalStyle)
		}
	}
	putString(blankLine, 0, screenHeight-2)
	putString(buffer.GetBottom(ontoY, screenHeight-2), 0, screenHeight-2)
}

func actionEnter() {
	buffer.Add('\n')
	for y := screenHeight - 2; y > yCursor+1; y-- {
		for x := 0; x < screenWidth; x++ {
			r, _, _, _ := screen.GetContent(x, y-1)
			screen.SetContent(x, y, r, nil, terminalStyle)
		}
	}
	for x := 0; x < screenWidth; x++ {
		r, _, _, _ := screen.GetContent(x+xCursor, yCursor)
		screen.SetContent(x, yCursor+1, r, nil, terminalStyle)
		screen.SetContent(x+xCursor, yCursor, ' ', nil, terminalStyle)
	}
	xCursor = 0
	if yCursor != screenHeight-2 {
		yCursor++
	} else {
		for y := 0; y < screenHeight-2; y++ {
			for x := 0; x < screenWidth; x++ {
				r, _, _, _ := screen.GetContent(x, y+1)
				screen.SetContent(x, y, r, nil, terminalStyle)
			}
		}
		putString(blankLine, 0, screenHeight-2)
		putString(buffer.GetLine(), 0, screenHeight-2)
	}
}

func actionKeyPress(ev *tcell.EventKey) {
	x, y := xCursor, yCursor
	requiredUpdates := buffer.Add(ev.Rune())
	xCursor++
	shiftRight(requiredUpdates)
	putRune(ev.Rune(), x, y)
}

func shiftRight(requiredUpdates int) {
	for i := xCursor + requiredUpdates; i >= xCursor; i-- {
		r, _, _, _ := screen.GetContent(i-1, yCursor)
		screen.SetContent(i, yCursor, r, nil, terminalStyle)
	}
}

func executeCommand(quit chan struct{}) {
	if len(command) > 1 && command[0] == '/' {
		removeHighlighting()
		search := command[1:]
		xSearchPoints, ySearchPoints = buffer.Search(search, yCursor, screenHeight)
		searchStringLength = len(search)
		for i := 0; i < len(xSearchPoints); i++ {
			startX, y := xSearchPoints[i], ySearchPoints[i]
			for x := startX; x < searchStringLength+startX; x++ {
				r, _, _, _ := screen.GetContent(x, y)
				screen.SetContent(x, y, r, nil, highlightStyle)
			}
		}
		return
	}
	switch command {
	case ":q":
		if buffer.CanSafeQuit() {
			close(quit)
		} else {
			displayError(modifiedFile)
		}
	case ":q!":
		close(quit)
	case ":w":
		write()
	case ":wq":
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
	mode = normalMode
	return true
}

func putRune(r rune, x, y int) {
	puts(screen, terminalStyle, x, y, string(r))
}

func putString(s string, x, y int) {
	puts(screen, terminalStyle, x, y, s)
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
