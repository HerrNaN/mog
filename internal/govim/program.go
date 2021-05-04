package govim

import (
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

type Program struct {
	buf *Buffer

	// The line number from the buffer that should be written to the
	// first line on screen
	offset int

	screen tcell.Screen
}

type dir int

const (
	dirUp dir = iota
	dirDown
	dirRight
	dirLeft
)

// NewState creates a new program empty state
func NewProgram() *Program {
	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	return &Program{
		screen: s,
		buf:    EmptyBuffer(),
	}
}

func NewProgramAt(filename string) *Program {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	buf, err := openFileIn(filename)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	return &Program{
		screen: s,
		buf:    buf,
	}
}

func openFileIn(filename string) (*Buffer, error) {
	bs, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return BufferFrom(string(bs)), nil
}

// Start sets up the program and starts the event loop
func (p *Program) Start() {
	p.screen.Clear()
	p.showCursorIn(p.buf, p.screen)
	p.run()
}

// Quit exits the program
func (p *Program) Quit() {
	p.screen.Fini()
	os.Exit(0)
}

func (p *Program) Show() {
	p.showBufferIn(p.buf, p.screen)
}

func (p *Program) showCursorIn(buf *Buffer, s tcell.Screen) {
	s.ShowCursor(buf.cursorX, buf.cursorY)
}

func (p *Program) moveCursor(buf *Buffer, s tcell.Screen, d dir) {
	w, h := s.Size()
	buf.MoveCursor(d)
	if buf.cursorY-p.offset == h {
		p.offset++
		p.showBufferIn(buf, s)
	}
	if buf.cursorY-p.offset == -1 {
		p.offset--
		p.showBufferIn(buf, s)
	}
	if buf.cursorX > len(buf.buffer[buf.cursorY]) {
		x, y := bufferPosToViewPos(buf, p.offset, w, len(buf.buffer[buf.cursorY]), buf.cursorY)
		s.ShowCursor(x, y)
		return
	}
	x, y := bufferPosToViewPos(buf, p.offset, w, buf.cursorX, buf.cursorY)
	s.ShowCursor(x, y)
}

func (p *Program) showBufferIn(buf *Buffer, s tcell.Screen) {
	w, h := s.Size()
	s.Clear()
	for bufY := range buf.buffer {
		for bufX, r := range buf.buffer[bufY] {
			if bufY < p.offset || bufY > p.offset+h {
				continue
			}
			x, y := bufferPosToViewPos(buf, p.offset, w, bufX, bufY)
			s.SetContent(x, y, r, nil, tcell.StyleDefault)
		}
	}
	for i := len(buf.buffer) - p.offset; i < h; i++ {
		s.SetContent(0, i, '~', nil, tcell.StyleDefault)
	}
	s.Show()
}

func (p *Program) run() {
	for {
		// Update screen
		p.Show()

		// Poll event
		ev := p.screen.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			p.screen.Sync()
			p.showBufferIn(p.buf, p.screen)
		case *tcell.EventKey:
			p.handleEventKey(*ev)
		default:
			log.Print(ev)
		}
	}
}

func (p *Program) handleEventKey(ev tcell.EventKey) {
	switch ev.Key() {
	//case tcell.KeyRune:
	//	p.setContent(ev.Rune())
	//	p.moveCursor(p.buf, p.screen, dirRight)
	case tcell.KeyEscape:
		p.Quit()
	case tcell.KeyUp:
		p.moveCursor(p.buf, p.screen, dirUp)
	case tcell.KeyDown:
		p.moveCursor(p.buf, p.screen, dirDown)
	case tcell.KeyRight:
		p.moveCursor(p.buf, p.screen, dirRight)
	case tcell.KeyLeft:
		p.moveCursor(p.buf, p.screen, dirLeft)
		//case tcell.KeyBackspace2:
		//	p.moveCursor(p.buf, p.screen, dirLeft)
		//	p.setContent(ev.Rune())
		//case tcell.KeyDelete:
		//	p.setContent(ev.Rune())
	}
}

func (p *Program) setContent(buf Buffer, s tcell.Screen, r rune) {
	s.SetContent(buf.cursorX, buf.cursorY, r, nil, tcell.StyleDefault)
}

func bufferPosToViewPos(buf *Buffer, bufferOffset, w int, bufX, bufY int) (int, int) {
	x := bufX % w

	offs := 0
	for i := 0; i < bufY-bufferOffset; i++ {
		offs += 1 + len(buf.buffer[bufferOffset+i])/w
	}

	y := bufX/w + offs

	return x, y
}
