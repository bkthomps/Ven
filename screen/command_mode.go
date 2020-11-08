package screen

import (
	"github.com/bkthomps/Ven/buffer"
	"github.com/bkthomps/Ven/search"
	"github.com/gdamore/tcell"
)

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
