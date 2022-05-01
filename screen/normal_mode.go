package screen

import (
	"github.com/bkthomps/Ven/buffer"
	"github.com/gdamore/tcell/v2"
)

func (screen *Screen) executeNormalMode(ev *tcell.EventKey) {
	previousCommand := screen.command.old
	screen.command.old = ""
	switch ev.Key() {
	case tcell.KeyDown, tcell.KeyUp, tcell.KeyLeft, tcell.KeyRight:
		screen.bufferAction(ev)
	case tcell.KeyCtrlF:
		screen.file.xCursor = screen.file.buffer.StartOfLine()
		for i := 0; i < screen.file.height; i++ {
			if screen.file.buffer.Current.Next == nil {
				break
			}
			screen.firstLine = screen.firstLine.Next
			screen.file.buffer.Current = screen.file.buffer.Current.Next
		}
		screen.completeDraw(nil)
	case tcell.KeyCtrlB:
		screen.file.xCursor = screen.file.buffer.StartOfLine()
		for i := 0; i < screen.file.height; i++ {
			if screen.firstLine.Prev == nil {
				break
			}
			screen.firstLine = screen.firstLine.Prev
			screen.file.buffer.Current = screen.file.buffer.Current.Prev
		}
		screen.completeDraw(nil)
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
		case 'a':
			screen.mode = insertMode
			screen.actionRight()
		case 'A':
			screen.mode = insertMode
			screen.file.xCursor = screen.file.buffer.EndOfLine(screen.mode == insertMode)
		case 'I':
			screen.mode = insertMode
			screen.file.xCursor = screen.file.buffer.StartOfLine()
		case 'o':
			screen.file.xCursor = screen.file.buffer.EndOfLine(screen.mode == insertMode)
			screen.actionKeyPress('\n')
			screen.mode = insertMode
		case 'O':
			screen.file.xCursor = screen.file.buffer.StartOfLine()
			screen.actionKeyPress('\n')
			screen.actionUp()
			screen.mode = insertMode
		case ':', '/':
			screen.mode = commandMode
			screen.command.current = buffer.Line{Data: []rune{ev.Rune()}}
			screen.command.runeOffset = 1
			screen.command.spaceOffset = buffer.RuneWidthJump(ev.Rune(), 0)
		case 'H':
			screen.file.xCursor = screen.file.buffer.StartOfLine()
			screen.navigateLineTop(0)
		case 'M':
			screen.file.xCursor = screen.file.buffer.StartOfLine()
			height := screen.maxHeight()
			screen.navigateLineTop(height / 2)
			screen.navigateLineBottom(height / 2)
		case 'L':
			screen.file.xCursor = screen.file.buffer.StartOfLine()
			height := screen.maxHeight()
			screen.navigateLineBottom(height)
		case '0':
			screen.file.xCursor = screen.file.buffer.StartOfLine()
		case '$':
			screen.file.xCursor = screen.file.buffer.EndOfLine(screen.mode == insertMode)
		case 'g':
			if previousCommand != "g" {
				screen.command.old = "g"
				break
			}
			screen.file.xCursor = screen.file.buffer.JumpToTop()
			screen.file.yCursor = 0
			screen.firstLine = screen.file.buffer.Current
			screen.completeDraw(nil)
		case 'G':
			screen.file.xCursor = screen.file.buffer.JumpToBottom()
			screen.file.yCursor = screen.file.height - 1
			screen.firstLine = screen.file.buffer.Current
			for i := 0; i < screen.file.height-1; i++ {
				if screen.firstLine.Prev == nil {
					break
				}
				screen.firstLine = screen.firstLine.Prev
			}
			screen.completeDraw(nil)
		case 'w':
			x, linesDown := screen.file.buffer.NextWordStart()
			screen.file.xCursor = x
			for i := 0; i < linesDown; i++ {
				if screen.file.yCursor == screen.file.height-1 {
					screen.firstLine = screen.firstLine.Next
				} else {
					screen.file.yCursor++
				}
			}
			if linesDown > 0 {
				screen.completeDraw(nil)
			}
		case 'b':
			x, linesUp := screen.file.buffer.PrevWordStart()
			screen.file.xCursor = x
			for i := 0; i < linesUp; i++ {
				if screen.file.yCursor == 0 {
					screen.firstLine = screen.firstLine.Prev
				} else {
					screen.file.yCursor--
				}
			}
			if linesUp > 0 {
				screen.completeDraw(nil)
			}
		case 'e':
			x, linesDown := screen.file.buffer.NextWordEnd()
			screen.file.xCursor = x
			for i := 0; i < linesDown; i++ {
				if screen.file.yCursor == screen.file.height-1 {
					screen.firstLine = screen.firstLine.Next
				} else {
					screen.file.yCursor++
				}
			}
			if linesDown > 0 {
				screen.completeDraw(nil)
			}
		case 'x':
			screen.file.xCursor = screen.file.buffer.Remove()
			screen.drawLine(screen.file.yCursor, screen.file.buffer.Current.Data)
		case 'X':
			screen.file.xCursor = screen.file.buffer.RemoveBefore()
			screen.drawLine(screen.file.yCursor, screen.file.buffer.Current.Data)
		case 'd':
			if previousCommand != "d" {
				screen.command.old = "d"
				break
			}
			x, wasFirst, wasLast := screen.file.buffer.RemoveLine(screen.mode == insertMode)
			screen.file.xCursor = x
			if wasFirst {
				screen.firstLine = screen.firstLine.Next
			} else if wasLast {
				if screen.file.yCursor == 0 {
					screen.firstLine = screen.firstLine.Prev
				} else {
					screen.file.yCursor--
				}
			}
			screen.completeDraw(nil)
		case 'D':
			screen.file.xCursor = screen.file.buffer.RemoveRestOfLine(screen.mode == insertMode)
			screen.drawLine(screen.file.yCursor, screen.file.buffer.Current.Data)
		}
	}
}

func (screen *Screen) navigateLineTop(lineIndex int) {
	for screen.file.yCursor > lineIndex {
		isPossible, _ := screen.file.buffer.Up(screen.mode == insertMode)
		if !isPossible {
			break
		}
		screen.file.yCursor--
	}
}

func (screen *Screen) navigateLineBottom(lineIndex int) {
	for screen.file.yCursor < lineIndex {
		isPossible, _ := screen.file.buffer.Down(screen.mode == insertMode)
		if !isPossible {
			break
		}
		screen.file.yCursor++
	}
}

func (screen *Screen) maxHeight() int {
	height := screen.file.height - 1
	if screen.file.buffer.Lines < height {
		height = screen.file.buffer.Lines
	}
	return height
}
