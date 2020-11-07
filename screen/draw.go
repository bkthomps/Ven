package screen

import (
	"github.com/bkthomps/Ven/buffer"
	"github.com/bkthomps/Ven/search"
)

var style = terminalStyle

func (screen *Screen) drawLine(y int, runes []rune, cursorHighlight bool, matchInstances *[]search.MatchInstance) {
	screen.tCell.HideCursor()
	for i := 0; i < screen.width; i++ {
		screen.tCell.SetContent(i, y, ' ', nil, style)
	}
	x := 0
	for _, r := range runes {
		screen.tCell.SetContent(x, y, r, nil, style)
		x = buffer.RuneWidthJump(r, x)
	}
	if cursorHighlight {
		screen.tCell.ShowCursor(screen.file.xCursor, screen.file.yCursor)
	}
}
