package screen_old

import "github.com/gdamore/tcell"

func (state *State) actionDelete() {
	possible, newX, requiredUpdates, spacing := state.buffer.Remove()
	if possible {
		if state.xCursor != 0 {
			state.xCursor -= spacing
			state.shiftLeft(requiredUpdates)
		} else {
			state.xCursor = newX
			state.yCursor--
			for x := 0; x < requiredUpdates; x++ {
				r, _, _, _ := state.screen.GetContent(x, state.yCursor+1)
				state.screen.SetContent(x+newX, state.yCursor, r, nil, terminalStyle)
			}
			state.shiftUp(state.yCursor, state.screenHeight-2)
		}
	}
}

func (state *State) shiftLeft(requiredUpdates int) {
	for i := state.xCursor; i <= state.xCursor+requiredUpdates; i++ {
		r, _, _, _ := state.screen.GetContent(i+1, state.yCursor)
		state.screen.SetContent(i, state.yCursor, r, nil, terminalStyle)
	}
}

func (state *State) shiftUp(ontoY, bottomY int) {
	for y := ontoY + 1; y < bottomY; y++ {
		for x := 0; x < state.screenWidth; x++ {
			r, _, _, _ := state.screen.GetContent(x, y+1)
			state.screen.SetContent(x, y, r, nil, terminalStyle)
		}
	}
	state.putString(state.blankLine, 0, bottomY)
	state.putString(state.buffer.GetBottom(ontoY, bottomY), 0, bottomY)
}

func (state *State) actionEnter() {
	state.buffer.Add('\n')
	for y := state.screenHeight - 2; y > state.yCursor+1; y-- {
		for x := 0; x < state.screenWidth; x++ {
			r, _, _, _ := state.screen.GetContent(x, y-1)
			state.screen.SetContent(x, y, r, nil, terminalStyle)
		}
	}
	for x := 0; x < state.screenWidth; x++ {
		r, _, _, _ := state.screen.GetContent(x+state.xCursor, state.yCursor)
		state.screen.SetContent(x, state.yCursor+1, r, nil, terminalStyle)
		state.screen.SetContent(x+state.xCursor, state.yCursor, ' ', nil, terminalStyle)
	}
	state.xCursor = 0
	if state.yCursor != state.screenHeight-2 {
		state.yCursor++
	} else {
		for y := 0; y < state.screenHeight-2; y++ {
			for x := 0; x < state.screenWidth; x++ {
				r, _, _, _ := state.screen.GetContent(x, y+1)
				state.screen.SetContent(x, y, r, nil, terminalStyle)
			}
		}
		state.putString(state.blankLine, 0, state.screenHeight-2)
		state.putString(state.buffer.GetLine(), 0, state.screenHeight-2)
	}
}

func (state *State) actionKeyPress(ev *tcell.EventKey) {
	x, y := state.xCursor, state.yCursor
	requiredUpdates := state.buffer.Add(ev.Rune())
	if ev.Rune() == '\t' {
		state.xCursor += 4
	} else {
		state.xCursor++
	}
	state.shiftRight(requiredUpdates)
	state.putRune(ev.Rune(), x, y)
}

func (state *State) shiftRight(requiredUpdates int) {
	for i := state.xCursor + requiredUpdates; i >= state.xCursor; i-- {
		r, _, _, _ := state.screen.GetContent(i-1, state.yCursor)
		state.screen.SetContent(i, state.yCursor, r, nil, terminalStyle)
	}
}

func (state *State) executeCommand(quit chan struct{}) {
	if len(state.command) > 1 && state.command[0] == '/' {
		state.removeHighlighting()
		searchText := state.command[1:]
		state.highlight(searchText)
		return
	}
	switch state.command {
	case ":q":
		if state.buffer.CanSafeQuit() {
			close(quit)
		} else {
			state.displayError(modifiedFile)
		}
	case ":q!":
		close(quit)
	case ":w":
		state.write()
	case ":wq":
		saved := state.write()
		if saved {
			close(quit)
		}
	default:
		state.displayError(errorCommand)
	}
}

func (state *State) highlight(searchText string) {
	xPoints, yPoints := state.buffer.Search(searchText, state.yCursor, state.screenHeight)
	length := len(searchText)
	for i := 0; i < len(xPoints); i++ {
		startX, y := xPoints[i], yPoints[i]
		for x := startX; x < length+startX; x++ {
			r, _, _, _ := state.screen.GetContent(x, y)
			state.screen.SetContent(x, y, r, nil, highlightStyle)
		}
	}
	state.search = &search{
		xPoints: xPoints,
		yPoints: yPoints,
		length:  length,
	}
}

func (state *State) write() (saved bool) {
	err := state.buffer.Save()
	if err != nil {
		state.displayError(errorSave)
		return false
	}
	state.mode = normalMode
	return true
}

func (state *State) putRune(r rune, x, y int) {
	if r == '\t' {
		for i := 0; i < 4; i++ {
			state.putRune(' ', x, y)
		}
		return
	}
	puts(state.screen, terminalStyle, x, y, string(r))
}

func (state *State) putString(s string, x, y int) {
	puts(state.screen, terminalStyle, x, y, s)
}
