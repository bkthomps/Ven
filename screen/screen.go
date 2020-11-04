package screen

import (
	"log"

	"github.com/bkthomps/Ven/buffer"
	"github.com/gdamore/tcell"
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

type file struct {
	xCursor int
	yCursor int

	height int
	width  int

	buffer *buffer.File
}

type command struct {
	xCursor   int
	yPosition int
	current   string
	old       string
}

type Screen struct {
	tCell     tcell.Screen
	mode      int
	firstLine *buffer.Line

	height int
	width  int

	blankLine string

	file    *file
	command *command
}

func (screen *Screen) Init(tCellScreen tcell.Screen, quit chan struct{}, fileName string) {
	screen.tCell = tCellScreen
	screen.mode = normalMode
	screen.command = &command{}
	buf := &buffer.File{}
	buf.Init(fileName)
	screen.firstLine = buf.First
	screen.file = &file{}
	screen.file.buffer = buf
	if err := screen.tCell.Init(); err != nil {
		log.Fatal(err)
	}
	screen.tCell.SetStyle(terminalStyle)
	screen.tCell.Show()
	screen.updateProperties()
	screen.completeDraw()
	screen.displayMode()
	go screen.listener(quit)
}

func (screen *Screen) updateProperties() {
	x, y := screen.tCell.Size()
	screen.height = y
	screen.width = x
	screen.file.height = y - 1
	screen.file.width = x
	screen.command.yPosition = y - 1
	screen.blankLine = ""
	for i := 0; i < x; i++ {
		screen.blankLine += " "
	}
	// TODO: update x or y cursor if now out of screen
}

func (screen *Screen) completeDraw() {
	y := 0
	for traverse := screen.firstLine; traverse != nil && y < screen.file.height; y++ {
		screen.drawLine(y, []rune(screen.blankLine), false)
		screen.drawLine(y, traverse.Data, true)
		traverse = traverse.Next
	}
	for y < screen.file.height {
		screen.drawLine(y, []rune(screen.blankLine), false)
		screen.drawLine(y, []rune{'~'}, true)
		y++
	}
}

func (screen *Screen) drawColor(x, y int, s tcell.Style) {
	r, _, _, _ := screen.tCell.GetContent(x, y)
	screen.tCell.SetContent(x, y, r, nil, s)
}

func (screen *Screen) displayError(error string) {
	screen.putCommand(screen.blankLine)
	screen.putCommand(error)
	screen.mode = commandErrorMode
	screen.displayMode()
}

func (screen *Screen) displayMode() {
	switch screen.mode {
	case insertMode:
		screen.putCommand(screen.blankLine)
		screen.putCommand(insertMessage)
	case normalMode:
		screen.putCommand(screen.blankLine)
	case commandMode:
		screen.putCommand(screen.blankLine)
		screen.putCommand(screen.command.current)
		screen.drawColor(screen.command.xCursor, screen.command.yPosition, cursorStyle)
	case commandErrorMode:
		screen.drawColor(screen.command.xCursor, screen.command.yPosition, terminalStyle)
	}
	screen.tCell.Sync()
}

func (screen *Screen) putCommand(str string) {
	screen.drawLine(screen.command.yPosition, []rune(str), true)
}

func (screen *Screen) listener(quit chan struct{}) {
	for {
		ev := screen.tCell.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch screen.mode {
			case insertMode:
				screen.executeInsertMode(ev)
			case normalMode:
				screen.executeNormalMode(ev)
			case commandMode:
				screen.executeCommandMode(ev, quit)
			case commandErrorMode:
				screen.mode = commandMode
			}
			screen.displayMode()
		case *tcell.EventResize:
			/* TODO
			isBigger := state.updateProperties()
			for state.xCursor >= state.screenWidth {
				state.actionLeft()
			}
			for state.yCursor >= state.screenHeight-1 {
				state.shiftUp(-1, state.screenHeight-1)
				state.yCursor--
			}
			if isBigger {
				arr := state.buffer.Redraw(state.yCursor, state.screenHeight)
				state.setInitial(arr)
			}
			state.displayMode()
			*/
		}
	}
}

func (screen *Screen) executeInsertMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEsc:
		screen.mode = normalMode
		screen.file.xCursor = screen.file.buffer.Left()
		screen.drawLine(screen.file.yCursor, screen.file.buffer.Current.Data, true)
	default:
		screen.bufferAction(ev)
	}
}

func (screen *Screen) executeNormalMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyDown, tcell.KeyUp, tcell.KeyLeft, tcell.KeyRight:
		screen.bufferAction(ev)
	default:
		switch ev.Rune() {
		case 'i':
			screen.command.old = ""
			screen.mode = insertMode
		case ':', '/':
			screen.command.old = ""
			screen.mode = commandMode
			screen.command.current = string(ev.Rune())
			screen.command.xCursor = len(screen.command.current)
		case 'x':
			screen.command.old = ""
			screen.file.xCursor = screen.file.buffer.Remove()
			screen.drawLine(screen.file.yCursor, screen.file.buffer.Current.Data, true)
		case 'X':
			screen.command.old = ""
			screen.file.xCursor = screen.file.buffer.RemoveBefore()
			screen.drawLine(screen.file.yCursor, screen.file.buffer.Current.Data, true)
		case 'd':
			if screen.command.old == "d" {
				x, wasFirst, wasLast := screen.file.buffer.RemoveLine(screen.mode == insertMode)
				screen.file.xCursor = x
				if wasFirst {
					screen.firstLine = screen.firstLine.Next
				} else if wasLast {
					screen.file.yCursor--
				}
				screen.completeDraw()
				screen.command.old = ""
			} else {
				screen.command.old = "d"
			}
		case 'D':
			screen.command.old = ""
			screen.file.xCursor = screen.file.buffer.RemoveRestOfLine(screen.mode == insertMode)
			screen.drawLine(screen.file.yCursor, []rune(screen.blankLine), false)
			screen.drawLine(screen.file.yCursor, screen.file.buffer.Current.Data, true)
		}
	}
}

func (screen *Screen) executeCommandMode(ev *tcell.EventKey, quit chan struct{}) {
	switch ev.Key() {
	case tcell.KeyEsc:
		screen.mode = normalMode
	case tcell.KeyEnter:
		screen.executeCommand(quit)
	case tcell.KeyDEL:
		if screen.command.xCursor <= 1 && len(screen.command.current) > 1 {
			break
		}
		runeCopy := []rune(screen.command.current)
		copy(runeCopy[screen.command.xCursor-1:], runeCopy[screen.command.xCursor:])
		shrinkSize := len(runeCopy) - 1
		runeCopy = runeCopy[:shrinkSize]
		screen.command.current = string(runeCopy)
		if shrinkSize == 0 {
			screen.mode = normalMode
		}
		screen.command.xCursor--
	case tcell.KeyDown, tcell.KeyUp:
		// Do Nothing
	case tcell.KeyLeft:
		if screen.command.xCursor > 1 {
			screen.command.xCursor--
		}
	case tcell.KeyRight:
		if screen.command.xCursor < len(screen.command.current) {
			screen.command.xCursor++
		}
	default:
		screen.command.current += string(ev.Rune())
		screen.command.xCursor++
	}
}

func (screen *Screen) bufferAction(ev *tcell.EventKey) {
	screen.drawLine(screen.file.yCursor, screen.file.buffer.Current.Data, false)
	switch ev.Key() {
	case tcell.KeyDown:
		screen.actionDown()
	case tcell.KeyUp:
		screen.actionUp()
	case tcell.KeyLeft:
		screen.actionLeft()
	case tcell.KeyRight:
		screen.actionRight()
	case tcell.KeyDEL:
		screen.actionDelete()
	case tcell.KeyEnter:
		screen.actionKeyPress('\n')
	default:
		screen.actionKeyPress(ev.Rune())
	}
	screen.drawLine(screen.file.yCursor, screen.file.buffer.Current.Data, true)
}

func (screen *Screen) actionDown() {
	possible, x := screen.file.buffer.Down(screen.mode == insertMode)
	if !possible {
		return
	}
	screen.file.xCursor = x
	if screen.file.yCursor == screen.file.height-1 {
		screen.firstLine = screen.firstLine.Next
		screen.completeDraw()
	} else {
		screen.file.yCursor++
	}
}

func (screen *Screen) actionUp() {
	possible, x := screen.file.buffer.Up(screen.mode == insertMode)
	if !possible {
		return
	}
	screen.file.xCursor = x
	if screen.file.yCursor == 0 {
		screen.firstLine = screen.firstLine.Prev
		screen.completeDraw()
	} else {
		screen.file.yCursor--
	}
}

func (screen *Screen) actionLeft() {
	screen.file.xCursor = screen.file.buffer.Left()
}

func (screen *Screen) actionRight() {
	screen.file.xCursor = screen.file.buffer.Right(screen.mode == insertMode)
}

func (screen *Screen) actionDelete() {
	x, deletedLine := screen.file.buffer.Backspace()
	screen.file.xCursor = x
	if deletedLine {
		screen.file.yCursor--
		screen.completeDraw()
	}
}

func (screen *Screen) actionKeyPress(rune rune) {
	x, addedLine := screen.file.buffer.Add(rune)
	screen.file.xCursor = x
	if addedLine {
		if screen.file.yCursor == screen.file.height-1 {
			screen.firstLine = screen.firstLine.Next
		} else {
			screen.file.yCursor++
		}
		screen.completeDraw()
	}
}

func (screen *Screen) executeCommand(quit chan struct{}) {
	if len(screen.command.current) > 1 && screen.command.current[0] == '/' {
		// TODO: searching
		return
	}
	switch screen.command.current {
	case ":q":
		if screen.file.buffer.CanSafeQuit() {
			close(quit)
		} else {
			screen.displayError(modifiedFile)
		}
	case ":q!":
		close(quit)
	case ":w":
		screen.write()
	case ":wq":
		saved := screen.write()
		if saved {
			close(quit)
		}
	default:
		screen.displayError(errorCommand)
	}
}

func (screen *Screen) write() (saved bool) {
	err := screen.file.buffer.Save()
	if err != nil {
		screen.displayError(errorSave)
		return false
	}
	screen.mode = normalMode
	return true
}
