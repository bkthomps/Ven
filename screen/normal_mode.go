package screen

import (
	"github.com/bkthomps/Ven/buffer"
	"github.com/gdamore/tcell"
)

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
