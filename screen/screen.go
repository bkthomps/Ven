package screen

import (
	"log"

	"github.com/bkthomps/Ven/buffer"
	"github.com/gdamore/tcell"
)

const (
	insertMessage = "-- INSERT --"
	errorCommand  = "-- Invalid Command --"
	errorSave     = "-- Could Not Save File --"
	modifiedFile  = "-- File Has Been Modified Since Last Save --"
)

const (
	insertMode = iota
	normalMode
	commandMode
	commandErrorMode
)

var (
	terminalStyle  = tcell.StyleDefault.Foreground(tcell.ColorBlack)
	cursorStyle    = terminalStyle.Background(tcell.ColorDarkGray)
	highlightStyle = terminalStyle.Background(tcell.ColorYellow)
)

type State struct {
	screen         tcell.Screen
	mode           int
	screenHeight   int
	screenWidth    int
	xCursor        int
	yCursor        int
	xCommandCursor int
	command        string
	oldCommand     rune
	blankLine      string
	search         *search
	buffer         *buffer.File
}

type search struct {
	xPoints []int
	yPoints []int
	length  int
}

func (state *State) Init(screen tcell.Screen, quit chan struct{}, fileName string) {
	state.screen = screen
	state.mode = normalMode
	state.buffer = &buffer.File{}
	arr := state.buffer.Init(fileName)
	if e := state.screen.Init(); e != nil {
		log.Fatal(e)
	}
	state.screen.SetStyle(terminalStyle)
	state.screen.Show()
	state.updateProperties()
	state.setInitial(arr)
	state.setColor(state.xCursor, state.yCursor, cursorStyle)
	state.displayMode()
	go state.listener(quit)
}

func (state *State) updateProperties() (isBigger bool) {
	oldWidth, oldHeight := state.screenWidth, state.screenHeight
	x, y := state.screen.Size()
	state.screenWidth, state.screenHeight = x, y
	for i := 0; i < x; i++ {
		state.blankLine += " "
	}
	return oldWidth < state.screenWidth || oldHeight < state.screenHeight
}

func (state *State) setInitial(arr []rune) {
	x, y := 0, 0
	for i := 0; i < len(arr) && y < state.screenHeight-1; i++ {
		cur := arr[i]
		if cur == '\n' {
			y++
			x = 0
		} else if x < state.screenWidth {
			state.putRune(cur, x, y)
			x++
		}
	}
	for i := y; i < state.screenHeight-1; i++ {
		for j := 2; j < state.screenWidth; j++ {
			r1, _, _, _ := state.screen.GetContent(j-2, i)
			r2, _, _, _ := state.screen.GetContent(j-1, i)
			if r1 == ' ' && r2 == ' ' {
				break
			}
			state.screen.SetContent(j-2, i, ' ', nil, terminalStyle)
		}
		state.putRune('~', 0, i)
	}
}

func (state *State) setColor(x, y int, s tcell.Style) {
	r, _, _, _ := state.screen.GetContent(x, y)
	state.screen.SetContent(x, y, r, nil, s)
}

func (state *State) displayError(error string) {
	state.putCommand(state.blankLine)
	state.putCommand(error)
	state.mode = commandErrorMode
	state.displayMode()
}

func (state *State) displayMode() {
	switch state.mode {
	case insertMode:
		state.setColor(state.xCursor, state.yCursor, cursorStyle)
		state.putCommand(state.blankLine)
		state.putCommand(insertMessage)
	case normalMode:
		state.setColor(state.xCursor, state.yCursor, cursorStyle)
		state.putCommand(state.blankLine)
	case commandMode:
		state.putCommand(state.blankLine)
		state.putCommand(state.command)
		state.setColor(state.xCommandCursor, state.screenHeight-1, cursorStyle)
	case commandErrorMode:
		state.setColor(state.xCommandCursor, state.screenHeight-1, terminalStyle)
	}
	state.screen.Sync()
}

func (state *State) putCommand(str string) {
	puts(state.screen, terminalStyle, 0, state.screenHeight-1, str)
}

func (state *State) listener(quit chan struct{}) {
	for {
		ev := state.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch state.mode {
			case insertMode:
				state.executeInsertMode(ev)
			case normalMode:
				state.executeNormalMode(ev)
			case commandMode:
				state.executeCommandMode(ev, quit)
			case commandErrorMode:
				state.mode = commandMode
			}
			state.displayMode()
		case *tcell.EventResize:
			isBigger := state.updateProperties()
			for state.xCursor >= state.screenWidth {
				state.actionLeft()
			}
			for state.yCursor >= state.screenHeight-1 {
				state.shiftUp(-1, state.screenHeight-1)
				state.yCursor--
			}
			if isBigger {
				arr := state.buffer.Redraw(state.yCursor, state.screenHeight)
				state.setInitial(arr)
			}
			state.displayMode()
		}
	}
}

func (state *State) executeInsertMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEsc:
		state.setColor(state.xCursor, state.yCursor, terminalStyle)
		state.mode = normalMode
		possible := state.buffer.Left()
		if possible {
			state.xCursor--
		}
	default:
		state.bufferAction(ev)
	}
}

func (state *State) executeNormalMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyDown, tcell.KeyUp, tcell.KeyLeft, tcell.KeyRight:
		state.bufferAction(ev)
	default:
		switch ev.Rune() {
		case 'i':
			state.mode = insertMode
		case ':', '/':
			state.setColor(state.xCursor, state.yCursor, terminalStyle)
			state.mode = commandMode
			state.command = string(ev.Rune())
			state.xCommandCursor = len(state.command)
		case 'x':
			isPossible, xBack, requiredUpdates := state.buffer.RemoveCurrent()
			if isPossible {
				state.shiftLeft(requiredUpdates)
				if xBack {
					state.actionLeft()
				}
			}
		case 'X':
			if state.xCursor != 0 {
				state.actionLeft()
				isPossible, _, requiredUpdates := state.buffer.RemoveCurrent()
				if isPossible {
					state.shiftLeft(requiredUpdates)
				}
			}
		case 'd':
			if state.oldCommand == 'd' {
				state.xCursor = 0
				state.shiftUp(state.yCursor-1, state.screenHeight-2)
				yBack, isEmpty := state.buffer.RemoveLine()
				if yBack {
					state.actionUp()
				}
				if isEmpty {
					state.screen.SetContent(0, 0, ' ', nil, terminalStyle)
				}
				state.oldCommand = '_'
			} else {
				state.oldCommand = ev.Rune()
			}
		case 'D':
			requiredUpdates := state.buffer.RemoveRestOfLine()
			for i := state.xCursor; i <= state.xCursor+requiredUpdates; i++ {
				state.screen.SetContent(i, state.yCursor, ' ', nil, terminalStyle)
			}
			if state.xCursor > 0 {
				state.actionLeft()
			}
		}
	}
}

func (state *State) executeCommandMode(ev *tcell.EventKey, quit chan struct{}) {
	switch ev.Key() {
	case tcell.KeyEsc:
		state.removeHighlighting()
		state.mode = normalMode
	case tcell.KeyEnter:
		state.executeCommand(quit)
	case tcell.KeyDEL:
		if state.xCommandCursor <= 1 && len(state.command) > 1 {
			break
		}
		runeCopy := []rune(state.command)
		copy(runeCopy[state.xCommandCursor-1:], runeCopy[state.xCommandCursor:])
		shrinkSize := len(runeCopy) - 1
		runeCopy = runeCopy[:shrinkSize]
		state.command = string(runeCopy)
		if shrinkSize == 0 {
			state.removeHighlighting()
			state.mode = normalMode
		}
		state.xCommandCursor--
	case tcell.KeyDown, tcell.KeyUp:
		// Do Nothing
	case tcell.KeyLeft:
		if state.xCommandCursor > 1 {
			state.xCommandCursor--
		}
	case tcell.KeyRight:
		if state.xCommandCursor < len(state.command) {
			state.xCommandCursor++
		}
	default:
		state.command += string(ev.Rune())
		state.xCommandCursor++
	}
}

func (state *State) removeHighlighting() {
	search := state.search
	if search == nil {
		return
	}
	for i := 0; i < len(search.xPoints); i++ {
		startX, y := search.xPoints[i], search.yPoints[i]
		for x := startX; x < search.length+startX; x++ {
			r, _, _, _ := state.screen.GetContent(x, y)
			state.screen.SetContent(x, y, r, nil, terminalStyle)
		}
	}
	state.search = nil
}

func (state *State) bufferAction(ev *tcell.EventKey) {
	state.setColor(state.xCursor, state.yCursor, terminalStyle)
	switch ev.Key() {
	case tcell.KeyDown:
		state.actionDown()
	case tcell.KeyUp:
		state.actionUp()
	case tcell.KeyLeft:
		state.actionLeft()
	case tcell.KeyRight:
		state.actionRight()
	case tcell.KeyDEL:
		state.actionDelete()
	case tcell.KeyEnter:
		state.actionEnter()
	default:
		state.actionKeyPress(ev)
	}
	state.setColor(state.xCursor, state.yCursor, cursorStyle)
}

func (state *State) actionDown() {
	possible, x := state.buffer.Down(state.xCursor, state.mode == insertMode)
	if possible {
		if state.yCursor == state.screenHeight-2 {
			for y := 0; y < state.screenHeight-2; y++ {
				for x := 0; x < state.screenWidth; x++ {
					r, _, _, _ := state.screen.GetContent(x, y+1)
					state.screen.SetContent(x, y, r, nil, terminalStyle)
				}
			}
			state.putString(state.blankLine, 0, state.screenHeight-2)
			state.putString(state.buffer.GetLine(), 0, state.screenHeight-2)
		} else {
			state.yCursor++
		}
		state.xCursor = x
	}
}

func (state *State) actionUp() {
	possible, x := state.buffer.Up(state.xCursor, state.mode == insertMode)
	if possible {
		if state.yCursor == 0 {
			for y := state.screenHeight - 2; y > 0; y-- {
				for x := 0; x < state.screenWidth; x++ {
					r, _, _, _ := state.screen.GetContent(x, y-1)
					state.screen.SetContent(x, y, r, nil, terminalStyle)
				}
			}
			state.putString(state.blankLine, 0, 0)
			state.putString(state.buffer.GetLine(), 0, 0)
		} else {
			state.yCursor--
		}
		state.xCursor = x
	}
}

func (state *State) actionLeft() {
	possible := state.buffer.Left()
	if possible {
		state.xCursor--
	}
}

func (state *State) actionRight() {
	possible := state.buffer.Right(state.mode == insertMode)
	if possible {
		state.xCursor++
	}
}

func (state *State) actionDelete() {
	possible, newX, requiredUpdates := state.buffer.Remove()
	if possible {
		if state.xCursor != 0 {
			state.xCursor--
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
	state.xCursor++
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
	puts(state.screen, terminalStyle, x, y, string(r))
}

func (state *State) putString(s string, x, y int) {
	puts(state.screen, terminalStyle, x, y, s)
}
