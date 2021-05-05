package govim

import (
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type Cursor interface {
	XPos() int
	YPos() int
	MoveLeft()
	MoveRight()
	MoveUp()
	MoveDown()
	MoveTo(x, y int)
}

type TextBuffer interface {
}

type SimpleFrame struct {
	screen tcell.Screen
	buffer []string
	cursor Cursor
	offset int
}

func EmptyFrame() *SimpleFrame {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	return &SimpleFrame{
		screen: s,
		buffer: []string{""},
		cursor: NewSimpleCursor(),
	}
}

func NewFrame(bs []byte) *SimpleFrame {
	f := EmptyFrame()
	buf := strings.Split(string(bs), "\n")
	f.buffer = buf
	return f
}

// MoveCursor moves the Cursor in the given direction.
// Moving to the left of the first character on a line, to the right of
// the last character on a line, above the first line or below the the
// last line all result in no movement. Moving up or down with a x value
// larger than the new line will not affect the x value of the SimpleCursor.
// This value is instead adjusted to the length of the line when moving to
// the left.
//
// If the cursor is at the top of the screen but not on the first line of
// the content in the buffer the screen scrolls upwards. An equivalent
// mechanic is implemented when moving down when on the bottom of the screen.
func (f *SimpleFrame) MoveCursor(d dir) {
	switch d {
	case dirUp:
		if f.cursor.YPos() > 0 {
			f.cursor.MoveUp()
		}
	case dirDown:
		if f.cursor.YPos() < len(f.buffer)-1 {
			f.cursor.MoveDown()
		}
	case dirLeft:
		if f.cursor.XPos() > len(f.currentLine()) {
			f.cursor.MoveTo(len(f.currentLine())-1, f.cursor.YPos())
		}
		if f.cursor.XPos() > 0 {
			f.cursor.MoveLeft()
		}
	case dirRight:
		if f.cursor.XPos() < len(f.currentLine())-1 {
			f.cursor.MoveRight()
		}
	}
	_, h := f.screen.Size()
	if f.cursor.YPos()-f.offset == h {
		f.offset++
	}
	if f.cursor.YPos()-f.offset == -1 {
		f.offset--
	}
	f.showCursor()
}

func (f *SimpleFrame) bufferPosToViewPos(bufX, bufY int) (int, int) {
	w, _ := f.screen.Size()
	x := bufX % w

	offs := 0
	for i := 0; i < bufY-f.offset; i++ {
		offs += 1 + len(f.buffer[f.offset+i])/w
	}

	y := bufX/w + offs

	return x, y
}

func (f *SimpleFrame) currentLine() string {
	return f.buffer[f.cursor.YPos()]
}

func (f *SimpleFrame) Show() {
	f.screen.Clear()
	f.writeBufferToScreen()
	f.screen.Show()
	f.showCursor()
}

func (f *SimpleFrame) writeBufferToScreen() {
	_, h := f.screen.Size()
	for bufY := range f.buffer {
		for bufX, r := range f.buffer[bufY] {
			if bufY < f.offset || bufY > f.offset+h {
				continue
			}
			x, y := f.bufferPosToViewPos(bufX, bufY)
			f.screen.SetContent(x, y, r, nil, tcell.StyleDefault)
		}
	}
	for i := len(f.buffer) - f.offset; i < h; i++ {
		f.screen.SetContent(0, i, '~', nil, tcell.StyleDefault)
	}
}

func (f *SimpleFrame) showCursor() {
	x, y := f.cursorScreenPos()
	f.screen.ShowCursor(x, y)
}

func (f *SimpleFrame) cursorScreenPos() (int, int) {
	if f.cursor.XPos() > len(f.currentLine()) {
		if len(f.currentLine()) == 0 {
			return f.bufferPosToViewPos(0, f.cursor.YPos())
		}
		return f.bufferPosToViewPos(len(f.currentLine())-1, f.cursor.YPos())
	}
	return f.bufferPosToViewPos(f.cursor.XPos(), f.cursor.YPos())
}

func (f *SimpleFrame) Close() {
	f.screen.Fini()
}

func (f *SimpleFrame) PollEvent() tcell.Event {
	return f.screen.PollEvent()
}

func (f *SimpleFrame) Sync() {
	f.screen.Sync()
}
