package govim

import (
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

type Program struct {
	buf    *Buffer
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
		buf:    BufferFrom("Hello\nMy\nName\nIs\n  HAHAHAHAHAHAHAHAHAHAHAHAHAHAHAHAHAHAHAHAH\n\n\n\n\nasdasdasd\nasd\nasd\nasd\n"),
	}
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
	switch d {
	case dirUp:
		if buf.cursorY > 0 {
			buf.cursorY--
		} else if buf.offset > 0 {
			buf.offset--
		}
	case dirDown:
		if buf.cursorY == h-1 {
			if buf.offset+buf.cursorY < len(buf.buffer)-1 {
				buf.offset++
			}
		} else if buf.cursorY < h-1 {
			if buf.offset+buf.cursorY < len(buf.buffer)-1 {
				buf.cursorY++
			}
		} else {
			panic("cursor below screen")
		}
	case dirLeft:
		if buf.cursorX > len(buf.buffer[buf.cursorY]) {
			buf.cursorX = len(buf.buffer[buf.cursorY])
		}
		if buf.cursorX > 0 {
			buf.cursorX--
		}
	case dirRight:
		if buf.cursorX < len(buf.buffer[buf.cursorY]) {
			buf.cursorX++
		}
	}
	if buf.cursorX > len(buf.buffer[buf.cursorY]) {
		x, y := bufferPosToViewPos(buf, w, len(buf.buffer[buf.cursorY]), buf.cursorY)
		s.ShowCursor(x, y)
		return
	}
	x, y := bufferPosToViewPos(buf, w, buf.cursorX, buf.cursorY)
	s.ShowCursor(x, y)
}

func (p *Program) showBufferIn(buf *Buffer, s tcell.Screen) {
	w, _ := s.Size()
	s.Clear()
	for bufY := range buf.buffer {
		for bufX, r := range buf.buffer[bufY] {
			x, y := bufferPosToViewPos(buf, w, bufX, bufY)
			s.SetContent(x, y, r, nil, tcell.StyleDefault)
		}
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

func bufferPosToViewPos(buf *Buffer, w int, bufX, bufY int) (int, int) {
	x := bufX % w

	offs := 0
	for i := 0; i < bufY-buf.offset; i++ {
		offs += 1 + len(buf.buffer[buf.offset+i])/w
	}

	y := bufX/w + offs

	return x, y
}
