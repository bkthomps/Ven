package screen

import (
	"github.com/bkthomps/Ven/buffer"
	"github.com/gdamore/tcell"
)

func (screen *Screen) drawLine(y int, runes []rune, cursorHighlight bool) {
	x := 0
	for _, r := range runes {
		if r == '\t' {
			screen.tCell.SetContent(x, y, ' ', nil, screen.currentStyle(cursorHighlight, x, y))
			for i := x + 1; i < x+buffer.TabSize; i++ {
				screen.tCell.SetContent(i, y, ' ', nil, terminalStyle)
			}
			x += buffer.TabSize
			continue
		}
		screen.tCell.SetContent(x, y, r, nil, screen.currentStyle(cursorHighlight, x, y))
		x++
	}
	screen.tCell.SetContent(x, y, ' ', nil, screen.currentStyle(cursorHighlight, x, y))
	for i := x + 1; i < x+buffer.TabSize; i++ {
		screen.tCell.SetContent(i, y, ' ', nil, terminalStyle)
	}
}

func (screen *Screen) currentStyle(cursorHighlight bool, x, y int) tcell.Style {
	style := terminalStyle
	if cursorHighlight && y == screen.file.yCursor && x == screen.file.xCursor {
		style = cursorStyle
	}
	return style
}
