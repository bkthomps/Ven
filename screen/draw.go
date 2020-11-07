package screen

import (
	"github.com/bkthomps/Ven/search"
	"github.com/mattn/go-runewidth"
)

const zeroWidthJoiner = '\u200d'

var style = terminalStyle

func (screen *Screen) drawLine(y int, runes []rune, cursorHighlight bool, matchInstances *[]search.MatchInstance) {
	for i, r := range screen.blankLine {
		screen.tCell.SetContent(i, y, r, nil, style)
	}
	i := 0
	width := 0
	deferred := make([]rune, 0)
	isJoiner := false
	for _, r := range runes {
		if r == zeroWidthJoiner {
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				width = 1
			}
			deferred = append(deferred, r)
			isJoiner = true
			continue
		}
		if isJoiner {
			deferred = append(deferred, r)
			isJoiner = false
			continue
		}
		switch runewidth.RuneWidth(r) {
		case 0:
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				width = 1
			}
		case 1:
			if len(deferred) != 0 {
				screen.tCell.SetContent(i, y, deferred[0], deferred[1:], style)
				i += width
			}
			deferred = nil
			width = 1
		case 2:
			if len(deferred) != 0 {
				screen.tCell.SetContent(i, y, deferred[0], deferred[1:], style)
				i += width
			}
			deferred = nil
			width = 2
		}
		deferred = append(deferred, r)
	}
	if len(deferred) != 0 {
		screen.tCell.SetContent(i, y, deferred[0], deferred[1:], style)
		i += width
	}
	if cursorHighlight {
		screen.tCell.ShowCursor(screen.file.xCursor, screen.file.yCursor)
	} else {
		screen.tCell.HideCursor()
	}
}
