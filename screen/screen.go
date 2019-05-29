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

const (
	insertMessage = "-- INSERT --"
	errorCommand  = "-- Invalid Command --"
	errorSave     = "-- Could Not Save File --"
	modifiedFile  = "-- File Has Been Modified Since Last Save --"
)

const (
	insertMode = iota
	normalMode
	commandMode
	commandErrorMode
)

var (
	terminalStyle  = tcell.StyleDefault.Foreground(tcell.ColorBlack)
	cursorStyle    = terminalStyle.Background(tcell.ColorDarkGray)
	highlightStyle = terminalStyle.Background(tcell.ColorYellow)
)

type info struct {
	screen         tcell.Screen
	mode           int
	screenHeight   int
	screenWidth    int
	xCursor        int
	yCursor        int
	xCommandCursor int
	command        string
	oldCommand     rune
	blankLine      string
	search         *search
	fileBuffer     *buffer.Info
}

type search struct {
	xPoints []int
	yPoints []int
	length  int
}

func Init(screen tcell.Screen, quit chan struct{}, fileName string) {
	fileBuffer, arr := buffer.Init(fileName)
	screenInfo := &info{
		screen:     screen,
		mode:       normalMode,
		fileBuffer: fileBuffer,
	}
	if e := screenInfo.screen.Init(); e != nil {
		log.Fatal(e)
	}
	screenInfo.screen.SetStyle(terminalStyle)
	screenInfo.screen.Show()
	updateProperties(screenInfo)
	setInitial(screenInfo, arr)
	setColor(screenInfo.screen, screenInfo.xCursor, screenInfo.yCursor, cursorStyle)
	displayMode(screenInfo)
	go listener(quit, screenInfo)
}

func updateProperties(info *info) (isBigger bool) {
	oldWidth, oldHeight := info.screenWidth, info.screenHeight
	x, y := info.screen.Size()
	info.screenWidth, info.screenHeight = x, y
	for i := 0; i < x; i++ {
		info.blankLine += " "
	}
	return oldWidth < info.screenWidth || oldHeight < info.screenHeight
}

func setInitial(info *info, arr []rune) {
	x, y := 0, 0
	for i := 0; i < len(arr) && y < info.screenHeight-1; i++ {
		cur := arr[i]
		if cur == '\n' {
			y++
			x = 0
		} else if x < info.screenWidth {
			putRune(info.screen, cur, x, y)
			x++
		}
	}
	for i := y; i < info.screenHeight-1; i++ {
		for j := 2; j < info.screenWidth; j++ {
			r1, _, _, _ := info.screen.GetContent(j-2, i)
			r2, _, _, _ := info.screen.GetContent(j-1, i)
			if r1 == ' ' && r2 == ' ' {
				break
			}
			info.screen.SetContent(j-2, i, ' ', nil, terminalStyle)
		}
		putRune(info.screen, '~', 0, i)
	}
}

func setColor(screen tcell.Screen, x, y int, s tcell.Style) {
	r, _, _, _ := screen.GetContent(x, y)
	screen.SetContent(x, y, r, nil, s)
}

func displayError(info *info, error string) {
	putCommand(info, info.blankLine)
	putCommand(info, error)
	info.mode = commandErrorMode
	displayMode(info)
}

func displayMode(info *info) {
	switch info.mode {
	case insertMode:
		setColor(info.screen, info.xCursor, info.yCursor, cursorStyle)
		putCommand(info, info.blankLine)
		putCommand(info, insertMessage)
	case normalMode:
		setColor(info.screen, info.xCursor, info.yCursor, cursorStyle)
		putCommand(info, info.blankLine)
	case commandMode:
		putCommand(info, info.blankLine)
		putCommand(info, info.command)
		setColor(info.screen, info.xCommandCursor, info.screenHeight-1, cursorStyle)
	case commandErrorMode:
		setColor(info.screen, info.xCommandCursor, info.screenHeight-1, terminalStyle)
	}
	info.screen.Sync()
}

func putCommand(info *info, str string) {
	puts(info.screen, terminalStyle, 0, info.screenHeight-1, str)
}

func listener(quit chan struct{}, info *info) {
	for {
		ev := info.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch info.mode {
			case insertMode:
				executeInsertMode(info, ev)
			case normalMode:
				executeNormalMode(info, ev)
			case commandMode:
				executeCommandMode(info, ev, quit)
			case commandErrorMode:
				info.mode = commandMode
			}
			displayMode(info)
		case *tcell.EventResize:
			isBigger := updateProperties(info)
			for info.xCursor >= info.screenWidth {
				actionLeft(info)
			}
			for info.yCursor >= info.screenHeight-1 {
				shiftUp(info, -1, info.screenHeight-1)
				info.yCursor--
			}
			if isBigger {
				arr := buffer.Redraw(info.fileBuffer, info.yCursor, info.screenHeight)
				setInitial(info, arr)
			}
			displayMode(info)
		}
	}
}

func executeInsertMode(info *info, ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEsc:
		setColor(info.screen, info.xCursor, info.yCursor, terminalStyle)
		info.mode = normalMode
		possible := buffer.Left(info.fileBuffer)
		if possible {
			info.xCursor--
		}
	default:
		bufferAction(info, ev)
	}
}

func executeNormalMode(info *info, ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyDown, tcell.KeyUp, tcell.KeyLeft, tcell.KeyRight:
		bufferAction(info, ev)
	default:
		switch ev.Rune() {
		case 'i':
			info.mode = insertMode
		case ':', '/':
			setColor(info.screen, info.xCursor, info.yCursor, terminalStyle)
			info.mode = commandMode
			info.command = string(ev.Rune())
			info.xCommandCursor = len(info.command)
		case 'x':
			isPossible, xBack, requiredUpdates := buffer.RemoveCurrent(info.fileBuffer)
			if isPossible {
				shiftLeft(info, requiredUpdates)
				if xBack {
					actionLeft(info)
				}
			}
		case 'X':
			if info.xCursor != 0 {
				actionLeft(info)
				isPossible, _, requiredUpdates := buffer.RemoveCurrent(info.fileBuffer)
				if isPossible {
					shiftLeft(info, requiredUpdates)
				}
			}
		case 'd':
			if info.oldCommand == 'd' {
				info.xCursor = 0
				shiftUp(info, info.yCursor-1, info.screenHeight-2)
				yBack, isEmpty := buffer.RemoveLine(info.fileBuffer)
				if yBack {
					actionUp(info)
				}
				if isEmpty {
					info.screen.SetContent(0, 0, ' ', nil, terminalStyle)
				}
				info.oldCommand = '_'
			} else {
				info.oldCommand = ev.Rune()
			}
		case 'D':
			requiredUpdates := buffer.RemoveRestOfLine(info.fileBuffer)
			for i := info.xCursor; i <= info.xCursor+requiredUpdates; i++ {
				info.screen.SetContent(i, info.yCursor, ' ', nil, terminalStyle)
			}
		}
	}
}

func executeCommandMode(info *info, ev *tcell.EventKey, quit chan struct{}) {
	switch ev.Key() {
	case tcell.KeyEsc:
		removeHighlighting(info)
		info.mode = normalMode
	case tcell.KeyEnter:
		executeCommand(quit, info)
	case tcell.KeyDEL:
		if info.xCommandCursor <= 1 && len(info.command) > 1 {
			break
		}
		runeCopy := []rune(info.command)
		copy(runeCopy[info.xCommandCursor-1:], runeCopy[info.xCommandCursor:])
		shrinkSize := len(runeCopy) - 1
		runeCopy = runeCopy[:shrinkSize]
		info.command = string(runeCopy)
		if shrinkSize == 0 {
			removeHighlighting(info)
			info.mode = normalMode
		}
		info.xCommandCursor--
	case tcell.KeyDown, tcell.KeyUp:
		// Do Nothing
	case tcell.KeyLeft:
		if info.xCommandCursor > 1 {
			info.xCommandCursor--
		}
	case tcell.KeyRight:
		if info.xCommandCursor < len(info.command) {
			info.xCommandCursor++
		}
	default:
		info.command += string(ev.Rune())
		info.xCommandCursor++
	}
}

func removeHighlighting(info *info) {
	search := info.search
	if search == nil {
		return
	}
	for i := 0; i < len(search.xPoints); i++ {
		startX, y := search.xPoints[i], search.yPoints[i]
		for x := startX; x < search.length+startX; x++ {
			r, _, _, _ := info.screen.GetContent(x, y)
			info.screen.SetContent(x, y, r, nil, terminalStyle)
		}
	}
	info.search = nil
}

func bufferAction(info *info, ev *tcell.EventKey) {
	setColor(info.screen, info.xCursor, info.yCursor, terminalStyle)
	switch ev.Key() {
	case tcell.KeyDown:
		actionDown(info)
	case tcell.KeyUp:
		actionUp(info)
	case tcell.KeyLeft:
		actionLeft(info)
	case tcell.KeyRight:
		actionRight(info)
	case tcell.KeyDEL:
		actionDelete(info)
	case tcell.KeyEnter:
		actionEnter(info)
	default:
		actionKeyPress(info, ev)
	}
	setColor(info.screen, info.xCursor, info.yCursor, cursorStyle)
}

func actionDown(info *info) {
	possible, x := buffer.Down(info.fileBuffer, info.xCursor, info.mode == insertMode)
	if possible {
		if info.yCursor == info.screenHeight-2 {
			for y := 0; y < info.screenHeight-2; y++ {
				for x := 0; x < info.screenWidth; x++ {
					r, _, _, _ := info.screen.GetContent(x, y+1)
					info.screen.SetContent(x, y, r, nil, terminalStyle)
				}
			}
			putString(info.screen, info.blankLine, 0, info.screenHeight-2)
			putString(info.screen, buffer.GetLine(info.fileBuffer), 0, info.screenHeight-2)
		} else {
			info.yCursor++
		}
		info.xCursor = x
	}
}

func actionUp(info *info) {
	possible, x := buffer.Up(info.fileBuffer, info.xCursor, info.mode == insertMode)
	if possible {
		if info.yCursor == 0 {
			for y := info.screenHeight - 2; y > 0; y-- {
				for x := 0; x < info.screenWidth; x++ {
					r, _, _, _ := info.screen.GetContent(x, y-1)
					info.screen.SetContent(x, y, r, nil, terminalStyle)
				}
			}
			putString(info.screen, info.blankLine, 0, 0)
			putString(info.screen, buffer.GetLine(info.fileBuffer), 0, 0)
		} else {
			info.yCursor--
		}
		info.xCursor = x
	}
}

func actionLeft(info *info) {
	possible := buffer.Left(info.fileBuffer)
	if possible {
		info.xCursor--
	}
}

func actionRight(info *info) {
	possible := buffer.Right(info.fileBuffer, info.mode == insertMode)
	if possible {
		info.xCursor++
	}
}

func actionDelete(info *info) {
	possible, newX, requiredUpdates := buffer.Remove(info.fileBuffer)
	if possible {
		if info.xCursor != 0 {
			info.xCursor--
			shiftLeft(info, requiredUpdates)
		} else {
			info.xCursor = newX
			info.yCursor--
			for x := 0; x < requiredUpdates; x++ {
				r, _, _, _ := info.screen.GetContent(x, info.yCursor+1)
				info.screen.SetContent(x+newX, info.yCursor, r, nil, terminalStyle)
			}
			shiftUp(info, info.yCursor, info.screenHeight-2)
		}
	}
}

func shiftLeft(info *info, requiredUpdates int) {
	for i := info.xCursor; i <= info.xCursor+requiredUpdates; i++ {
		r, _, _, _ := info.screen.GetContent(i+1, info.yCursor)
		info.screen.SetContent(i, info.yCursor, r, nil, terminalStyle)
	}
}

func shiftUp(info *info, ontoY, bottomY int) {
	for y := ontoY + 1; y < bottomY; y++ {
		for x := 0; x < info.screenWidth; x++ {
			r, _, _, _ := info.screen.GetContent(x, y+1)
			info.screen.SetContent(x, y, r, nil, terminalStyle)
		}
	}
	putString(info.screen, info.blankLine, 0, bottomY)
	putString(info.screen, buffer.GetBottom(info.fileBuffer, ontoY, bottomY), 0, bottomY)
}

func actionEnter(info *info) {
	buffer.Add(info.fileBuffer, '\n')
	for y := info.screenHeight - 2; y > info.yCursor+1; y-- {
		for x := 0; x < info.screenWidth; x++ {
			r, _, _, _ := info.screen.GetContent(x, y-1)
			info.screen.SetContent(x, y, r, nil, terminalStyle)
		}
	}
	for x := 0; x < info.screenWidth; x++ {
		r, _, _, _ := info.screen.GetContent(x+info.xCursor, info.yCursor)
		info.screen.SetContent(x, info.yCursor+1, r, nil, terminalStyle)
		info.screen.SetContent(x+info.xCursor, info.yCursor, ' ', nil, terminalStyle)
	}
	info.xCursor = 0
	if info.yCursor != info.screenHeight-2 {
		info.yCursor++
	} else {
		for y := 0; y < info.screenHeight-2; y++ {
			for x := 0; x < info.screenWidth; x++ {
				r, _, _, _ := info.screen.GetContent(x, y+1)
				info.screen.SetContent(x, y, r, nil, terminalStyle)
			}
		}
		putString(info.screen, info.blankLine, 0, info.screenHeight-2)
		putString(info.screen, buffer.GetLine(info.fileBuffer), 0, info.screenHeight-2)
	}
}

func actionKeyPress(info *info, ev *tcell.EventKey) {
	x, y := info.xCursor, info.yCursor
	requiredUpdates := buffer.Add(info.fileBuffer, ev.Rune())
	info.xCursor++
	shiftRight(info, requiredUpdates)
	putRune(info.screen, ev.Rune(), x, y)
}

func shiftRight(info *info, requiredUpdates int) {
	for i := info.xCursor + requiredUpdates; i >= info.xCursor; i-- {
		r, _, _, _ := info.screen.GetContent(i-1, info.yCursor)
		info.screen.SetContent(i, info.yCursor, r, nil, terminalStyle)
	}
}

func executeCommand(quit chan struct{}, info *info) {
	if len(info.command) > 1 && info.command[0] == '/' {
		removeHighlighting(info)
		searchText := info.command[1:]
		highlight(info, searchText)
		return
	}
	switch info.command {
	case ":q":
		if buffer.CanSafeQuit(info.fileBuffer) {
			close(quit)
		} else {
			displayError(info, modifiedFile)
		}
	case ":q!":
		close(quit)
	case ":w":
		write(info)
	case ":wq":
		saved := write(info)
		if saved {
			close(quit)
		}
	default:
		displayError(info, errorCommand)
	}
}

func highlight(info *info, searchText string) {
	xPoints, yPoints := buffer.Search(info.fileBuffer, searchText, info.yCursor, info.screenHeight)
	length := len(searchText)
	for i := 0; i < len(xPoints); i++ {
		startX, y := xPoints[i], yPoints[i]
		for x := startX; x < length+startX; x++ {
			r, _, _, _ := info.screen.GetContent(x, y)
			info.screen.SetContent(x, y, r, nil, highlightStyle)
		}
	}
	info.search = &search{
		xPoints: xPoints,
		yPoints: yPoints,
		length:  length,
	}
}

func write(info *info) (saved bool) {
	err := buffer.Save(info.fileBuffer)
	if err != nil {
		displayError(info, errorSave)
		return false
	}
	info.mode = normalMode
	return true
}

func putRune(screen tcell.Screen, r rune, x, y int) {
	puts(screen, terminalStyle, x, y, string(r))
}

func putString(screen tcell.Screen, s string, x, y int) {
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
