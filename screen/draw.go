package screen

import (
	"github.com/bkthomps/Ven/buffer"
	"github.com/bkthomps/Ven/search"
)

func (screen *Screen) drawLine(y int, runes []rune, cursorHighlight bool, matchInstances *[]search.MatchInstance) {
	matchIndex := 0
	x := 0
	spacingOffset := 0
	for _, r := range runes {
		style := terminalStyle
		if matchInstances != nil && matchIndex < len(*matchInstances) {
			offset := (*matchInstances)[matchIndex].StartOffset + spacingOffset
			length := (*matchInstances)[matchIndex].Length
			if x >= offset && x < offset+length {
				style = highlightStyle
			}
			if x >= offset+length-1 {
				matchIndex++
			}
		}
		if cursorHighlight && y == screen.file.yCursor && x == screen.file.xCursor {
			style = cursorStyle
		}
		if r == '\t' {
			screen.tCell.SetContent(x, y, ' ', nil, style)
			for i := x + 1; i < x+buffer.TabSize; i++ {
				screen.tCell.SetContent(i, y, ' ', nil, style)
			}
			spacingOffset += buffer.TabSize - 1
			x += buffer.TabSize
			continue
		}
		screen.tCell.SetContent(x, y, r, nil, style)
		x++
	}
	style := terminalStyle
	if cursorHighlight && y == screen.file.yCursor && x == screen.file.xCursor {
		style = cursorStyle
	}
	screen.tCell.SetContent(x, y, ' ', nil, style)
	for i := x + 1; i < x+buffer.TabSize; i++ {
		screen.tCell.SetContent(i, y, ' ', nil, terminalStyle)
	}
}
