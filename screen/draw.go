package screen

import (
	"github.com/bkthomps/Ven/buffer"
	"github.com/bkthomps/Ven/search"
)

func (screen *Screen) drawLineHighlight(y int, runes []rune, instances []search.MatchInstance) {
	screen.drawBlankLine(y)
	matchIndex := 0
	x := 0
	for i, r := range runes {
		if matchIndex < len(instances) && i >= instances[matchIndex].StartOffset {
			xUpdated := buffer.RuneWidthJump(r, x)
			if r == '\t' {
				for j := x; j < xUpdated; j++ {
					screen.tCell.SetContent(j, y, ' ', nil, highlightStyle)
				}
			} else {
				screen.tCell.SetContent(x, y, r, nil, highlightStyle)
			}
			if i == instances[matchIndex].StartOffset+instances[matchIndex].Length-1 {
				matchIndex++
			}
			x = xUpdated
			continue
		}
		screen.tCell.SetContent(x, y, r, nil, terminalStyle)
		x = buffer.RuneWidthJump(r, x)
	}
	screen.tCell.ShowCursor(screen.file.xCursor, screen.file.yCursor)
}

func (screen *Screen) drawLine(y int, runes []rune) {
	screen.drawBlankLine(y)
	x := 0
	for _, r := range runes {
		screen.tCell.SetContent(x, y, r, nil, terminalStyle)
		x = buffer.RuneWidthJump(r, x)
	}
	screen.tCell.ShowCursor(screen.file.xCursor, screen.file.yCursor)
}

func (screen *Screen) drawBlankLine(y int) {
	screen.tCell.HideCursor()
	for i := 0; i < screen.width; i++ {
		screen.tCell.SetContent(i, y, ' ', nil, terminalStyle)
	}
	screen.tCell.ShowCursor(screen.file.xCursor, screen.file.yCursor)
}
