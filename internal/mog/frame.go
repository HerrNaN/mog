package mog

import (
	"fmt"
	"log"
	"os"
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
	screen       tcell.Screen
	buffer       []string
	cursor       Cursor
	offset       int
	filePath     string
	lockFilePath string
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
		screen:   s,
		buffer:   []string{""},
		cursor:   NewSimpleCursor(),
		filePath: "",
	}
}

func NewFrame(bs []byte) *SimpleFrame {
	f := EmptyFrame()
	buf := strings.Split(string(bs), "\n")
	f.buffer = buf
	return f
}

func NewFrameFromFile(filename string) *SimpleFrame {
	f := EmptyFrame()
	err := f.loadFile(filename)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	return f
}

func (f *SimpleFrame) loadFile(filePath string) error {
	lockFilePath := lockFilePathOf(filePath)
	if _, err := os.Stat(lockFilePath); err == nil {
		return fmt.Errorf("file open in another frame")
	}

	bs, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	file, err := os.Create(lockFilePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("%+v", err)
		}
	}(file)

	f.loadBuffer(bs)
	f.filePath = filePath
	f.lockFilePath = lockFilePath
	return nil
}

func (f *SimpleFrame) loadBuffer(bs []byte) {
	f.buffer = strings.Split(string(bs), "\n")
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

func (f *SimpleFrame) InsertRune(r rune) {
	if f.cursor.XPos() >= len(f.buffer[f.cursor.YPos()]) {
		var toX = len(f.currentLine())-1
		if f.currentLine() == "" {
			toX = 0
		}
		f.cursor.MoveTo(toX, f.cursor.YPos())
	}
	f.buffer[f.cursor.YPos()] = f.buffer[f.cursor.YPos()][:f.cursor.XPos()] + string(r) + f.buffer[f.cursor.YPos()][f.cursor.XPos():]
}

func (f *SimpleFrame) Show() {
	f.screen.Clear()
	f.writeBufferToScreen()
	f.screen.Show()
	f.showCursor()
}

func (f *SimpleFrame) writeBufferToScreen() {
	_, h := f.screen.Size()
	lastPrintedScreenLine := -1
	for bufY := range f.buffer {
		if f.buffer[bufY] == "" {
			lastPrintedScreenLine++
			continue
		}
		for bufX, r := range f.buffer[bufY] {
			if bufY < f.offset || bufY > f.offset+h {
				continue
			}
			x, y := f.bufferPosToViewPos(bufX, bufY)
			f.screen.SetContent(x, y, r, nil, tcell.StyleDefault)
			lastPrintedScreenLine = y
		}
	}
	for i := lastPrintedScreenLine + 1; i < h; i++ {
		f.screen.SetContent(0, i, '~', nil, tcell.StyleDefault)
	}
}

func (f *SimpleFrame) showCursor() {
	x, y := f.cursorScreenPos()
	f.screen.ShowCursor(x, y)
}

func (f *SimpleFrame) cursorScreenPos() (int, int) {
	if f.cursor.XPos() > len(f.currentLine()) {
		if f.currentLine() == "" {
			return f.bufferPosToViewPos(0, f.cursor.YPos())
		}
		return f.bufferPosToViewPos(len(f.currentLine())-1, f.cursor.YPos())
	}
	return f.bufferPosToViewPos(f.cursor.XPos(), f.cursor.YPos())
}

func (f *SimpleFrame) Close() error {
	f.screen.Fini()
	err := os.Remove(f.lockFilePath)
	if err != nil {
		return err
	}
	return nil
}

func (f *SimpleFrame) PollEvent() tcell.Event {
	return f.screen.PollEvent()
}

func (f *SimpleFrame) HandleEvent(e tcell.Event) bool {
	switch ev := e.(type) {
	case *tcell.EventResize:
		f.screen.Sync()
		f.Show()
	case *tcell.EventKey:
		return f.handleEventKey(*ev)
	default:
		log.Print(ev)
	}
	return false
}

func (f *SimpleFrame) handleEventKey(ev tcell.EventKey) bool {
	switch ev.Key() {
	case tcell.KeyEscape:
		err := f.Close()
		if err != nil {
			log.Fatalf("%+v", err)
		}
		return true
	case tcell.KeyUp:
		f.MoveCursor(dirUp)
	case tcell.KeyDown:
		f.MoveCursor(dirDown)
	case tcell.KeyRight:
		f.MoveCursor(dirRight)
	case tcell.KeyLeft:
		f.MoveCursor(dirLeft)
	case tcell.KeyRune:
		f.handleEventRune(ev.Rune())
	}
	return false
}

func (f *SimpleFrame) handleEventRune(r rune) {
	f.InsertRune(r)
	f.MoveCursor(dirRight)
}
