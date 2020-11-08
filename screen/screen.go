package screen

import (
	"log"

	"github.com/bkthomps/Ven/buffer"
	"github.com/bkthomps/Ven/search"
	"github.com/gdamore/tcell"
)

const (
	insertMode = iota
	normalMode
	commandMode
	commandErrorMode
)

var (
	insertMessage = []rune("-- INSERT --")
	errorCommand  = []rune("-- Invalid Command --")
	errorSave     = []rune("-- Could Not Save File --")
	modifiedFile  = []rune("-- File Has Been Modified Since Last Save --")
	badRegex      = []rune("-- Malformed Regex --")
)

var (
	terminalStyle  = tcell.StyleDefault.Foreground(tcell.ColorBlack)
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
	runeOffset  int
	spaceOffset int
	yPosition   int
	current     buffer.Line
	old         buffer.Line
}

type Screen struct {
	tCell     tcell.Screen
	mode      int
	firstLine *buffer.Line

	height int
	width  int

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
	screen.completeDraw(nil)
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
}

func (screen *Screen) completeDraw(matchLines *[]search.MatchLine) {
	matchIndex := 0
	y := 0
	for traverse := screen.firstLine; traverse != nil && y < screen.file.height; y++ {
		var matchInstances *[]search.MatchInstance = nil
		if matchLines != nil && matchIndex < len(*matchLines) && traverse == (*matchLines)[matchIndex].Line {
			matchInstances = &(*matchLines)[matchIndex].Instances
			matchIndex++
		}
		if matchInstances == nil {
			screen.drawLine(y, traverse.Data)
		} else {
			screen.drawLineHighlight(y, traverse.Data, *matchInstances)
		}
		traverse = traverse.Next
	}
	for y < screen.file.height {
		screen.drawLine(y, []rune{'~'})
		y++
	}
}

func (screen *Screen) displayError(error []rune) {
	screen.clearCommand()
	screen.putCommand(error)
	screen.mode = commandErrorMode
	screen.displayMode()
}

func (screen *Screen) displayMode() {
	switch screen.mode {
	case insertMode:
		screen.clearCommand()
		screen.putCommand(insertMessage)
	case normalMode:
		screen.clearCommand()
	case commandMode:
		screen.clearCommand()
		screen.putCommand(screen.command.current.Data)
		screen.tCell.ShowCursor(screen.command.spaceOffset, screen.command.yPosition)
	case commandErrorMode:
		screen.tCell.HideCursor()
	}
	screen.tCell.Sync()
}

func (screen *Screen) clearCommand() {
	screen.drawBlankLine(screen.command.yPosition)
}

func (screen *Screen) putCommand(runes []rune) {
	screen.drawLine(screen.command.yPosition, runes)
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
			screen.updateProperties()
			screen.completeDraw(nil)
			screen.displayMode()
		}
	}
}

func (screen *Screen) executeInsertMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEsc:
		screen.mode = normalMode
		screen.file.xCursor = screen.file.buffer.Left()
		screen.drawLine(screen.file.yCursor, screen.file.buffer.Current.Data)
	default:
		screen.bufferAction(ev)
	}
}

func (screen *Screen) executeNormalMode(ev *tcell.EventKey) {
	previousCommand := screen.command.old
	screen.command.old = buffer.Line{}
	switch ev.Key() {
	case tcell.KeyDown, tcell.KeyUp, tcell.KeyLeft, tcell.KeyRight:
		screen.bufferAction(ev)
	default:
		switch ev.Rune() {
		case 'j':
			screen.actionDown()
		case 'k':
			screen.actionUp()
		case 'h':
			screen.actionLeft()
		case 'l':
			screen.actionRight()
		case 'i':
			screen.mode = insertMode
		case ':', '/':
			screen.mode = commandMode
			screen.command.current = buffer.Line{Data: []rune{ev.Rune()}}
			screen.command.runeOffset = 1
			screen.command.spaceOffset = buffer.RuneWidthJump(ev.Rune(), 0)
		case 'x':
			screen.file.xCursor = screen.file.buffer.Remove()
			screen.drawLine(screen.file.yCursor, screen.file.buffer.Current.Data)
		case 'X':
			screen.file.xCursor = screen.file.buffer.RemoveBefore()
			screen.drawLine(screen.file.yCursor, screen.file.buffer.Current.Data)
		case 'd':
			if previousCommand.Equals("d") {
				x, wasFirst, wasLast := screen.file.buffer.RemoveLine(screen.mode == insertMode)
				screen.file.xCursor = x
				if wasFirst {
					screen.firstLine = screen.firstLine.Next
				} else if wasLast {
					screen.file.yCursor--
				}
				screen.completeDraw(nil)
			} else {
				screen.command.old = buffer.Line{Data: []rune("d")}
			}
		case 'D':
			screen.file.xCursor = screen.file.buffer.RemoveRestOfLine(screen.mode == insertMode)
			screen.drawLine(screen.file.yCursor, screen.file.buffer.Current.Data)
		}
	}
}

func (screen *Screen) executeCommandMode(ev *tcell.EventKey, quit chan struct{}) {
	switch ev.Key() {
	case tcell.KeyEsc:
		screen.mode = normalMode
	case tcell.KeyEnter:
		screen.executeCommand(quit)
	case tcell.KeyDown, tcell.KeyUp:
		// Do Nothing
	case tcell.KeyLeft:
		if screen.command.runeOffset > 1 {
			screen.command.runeOffset--
			r := screen.command.current.Data[screen.command.runeOffset]
			runes := screen.command.current.Data
			runeOffset := screen.command.runeOffset
			spaceOffset := screen.command.spaceOffset
			screen.command.spaceOffset = buffer.RuneWidthBackJump(r, runes, runeOffset, spaceOffset)
		}
	case tcell.KeyRight:
		if screen.command.runeOffset < len(screen.command.current.Data) {
			r := screen.command.current.Data[screen.command.runeOffset]
			screen.command.spaceOffset = buffer.RuneWidthJump(r, screen.command.spaceOffset)
			screen.command.runeOffset++
		}
	case tcell.KeyDEL:
		if screen.command.runeOffset == 1 && len(screen.command.current.Data) > 1 {
			break
		}
		screen.command.runeOffset--
		r := screen.command.current.Data[screen.command.runeOffset]
		runes := screen.command.current.Data
		runeOffset := screen.command.runeOffset
		spaceOffset := screen.command.spaceOffset
		screen.command.spaceOffset = buffer.RuneWidthBackJump(r, runes, runeOffset, spaceOffset)
		screen.command.current.RemoveAt(screen.command.runeOffset)
		if len(screen.command.current.Data) == 0 {
			screen.mode = normalMode
		}
	default:
		screen.command.current.AddAt(screen.command.runeOffset, ev.Rune())
		screen.command.runeOffset++
		screen.command.spaceOffset = buffer.RuneWidthJump(ev.Rune(), screen.command.spaceOffset)
	}
}

func (screen *Screen) bufferAction(ev *tcell.EventKey) {
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
	screen.drawLine(screen.file.yCursor, screen.file.buffer.Current.Data)
}

func (screen *Screen) actionDown() {
	possible, x := screen.file.buffer.Down(screen.mode == insertMode)
	if !possible {
		return
	}
	screen.file.xCursor = x
	if screen.file.yCursor == screen.file.height-1 {
		screen.firstLine = screen.firstLine.Next
		screen.completeDraw(nil)
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
		screen.completeDraw(nil)
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
		screen.completeDraw(nil)
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
		screen.completeDraw(nil)
	}
}

func (screen *Screen) executeCommand(quit chan struct{}) {
	if len(screen.command.current.Data) > 1 && screen.command.current.Data[0] == '/' {
		pattern := screen.command.current.Data[1:]
		matches, firstLineIndex, err := search.AllMatches(string(pattern), screen.firstLine, screen.file.height)
		if err != nil {
			screen.displayError(badRegex)
			return
		}
		if len(matches) > 0 && firstLineIndex > screen.file.height-screen.file.yCursor {
			for i := 0; i < firstLineIndex-1; i++ {
				_, x := screen.file.buffer.Down(screen.mode == insertMode)
				screen.file.xCursor = x
				screen.firstLine = screen.firstLine.Next
			}
		}
		screen.completeDraw(&matches)
		return
	}
	if screen.command.current.Equals(":q") {
		if screen.file.buffer.CanSafeQuit() {
			close(quit)
		} else {
			screen.displayError(modifiedFile)
		}
		return
	}
	if screen.command.current.Equals(":q!") {
		close(quit)
		return
	}
	if screen.command.current.Equals(":w") {
		screen.write()
		return
	}
	if screen.command.current.Equals(":wq") {
		saved := screen.write()
		if saved {
			close(quit)
		}
		return
	}
	screen.displayError(errorCommand)
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
