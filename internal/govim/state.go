package govim

import (
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type State struct {
	Screen           tcell.Screen
	cursorX, cursorY int

	// The buffer containing the contents currently being processed
	buffer []string

	// The line number from the buffer that should be written to the
	// first line on screen
	scrollX, scrollY int
}

type dir int

const (
	dirUp dir = iota
	dirDown
	dirRight
	dirLeft
)

// NewState creates a new program empty state
func NewState() *State {
	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	return &State{
		Screen:  s,
		cursorY: 0,
		cursorX: 0,
		buffer:  strings.Split("Hello\nMy\nName\nIs\n  HAHAHAHAHAHAHAHAHAHAHAHAHAHAHAHAHAHAHAHAH\n\n\n\n\nasdasdasd\nasd\nasd\nasd\n", "\n"),
	}
}

// Start sets up the program and starts the event loop
func (s *State) Start() {
	s.Screen.Clear()
	s.Screen.ShowCursor(s.cursorX, s.cursorY)
	s.run()
}

// Quit exits the program
func (s *State) Quit() {
	s.Screen.Fini()
	os.Exit(0)
}

func (s *State) moveCursor(d dir) {
	_, h := s.Screen.Size()
	switch d {
	case dirUp:
		if s.cursorY > 0 {
			s.cursorY--
		} else if s.scrollY > 0 {
			s.scrollY--
		}
	case dirDown:
		if s.cursorY == h-1 {
			if s.scrollY+s.cursorY < len(s.buffer)-1 {
				s.scrollY++
			}
		} else if s.cursorY < h-1 {
			if s.scrollY+s.cursorY < len(s.buffer)-1 {
				s.cursorY++
			}
		} else {
			panic("cursor below screen")
		}
	case dirLeft:
		if s.cursorX > 0 {
			s.cursorX--
		}
	case dirRight:
		if s.cursorX < len(s.buffer[s.cursorY]) {
			s.cursorX++
		}
	}
	if s.cursorX > len(s.buffer[s.cursorY]) {
		x, y := s.bufferPosToViewPos(len(s.buffer[s.cursorY]), s.cursorY, true)
		s.Screen.ShowCursor(x, y)
		return
	}
	x, y := s.bufferPosToViewPos(s.cursorX, s.cursorY, true)
	s.Screen.ShowCursor(x, y)
}

func (s *State) showCursor() {
	s.Screen.ShowCursor(s.cursorX, s.cursorY)
}

func (s *State) run() {
	for {
		// Update screen
		s.showBuffer()

		// Poll event
		ev := s.Screen.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Screen.Sync()
			s.moveCursorIntoBounds()
		case *tcell.EventKey:
			s.handleEventKey(*ev)
		default:
			log.Print(ev)
		}
	}
}

func (s *State) handleEventKey(ev tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyRune:
		s.setContent(ev.Rune())
		s.moveCursor(dirRight)
	case tcell.KeyEscape:
		s.Quit()
	case tcell.KeyUp:
		s.moveCursor(dirUp)
	case tcell.KeyDown:
		s.moveCursor(dirDown)
	case tcell.KeyRight:
		s.moveCursor(dirRight)
	case tcell.KeyLeft:
		s.moveCursor(dirLeft)
	case tcell.KeyBackspace2:
		s.moveCursor(dirLeft)
		s.setContent(ev.Rune())
	case tcell.KeyDelete:
		s.setContent(ev.Rune())
	}
}

func (s *State) setContent(r rune) {
	s.Screen.SetContent(s.cursorX, s.cursorY, r, nil, tcell.StyleDefault)
}

func (s *State) showBuffer() {
	s.Screen.Clear()
	for bufY := range s.buffer {
		for bufX, r := range s.buffer[bufY] {
			x, y := s.bufferPosToViewPos(bufX, bufY, true)
			s.Screen.SetContent(x, y, r, nil, tcell.StyleDefault)
		}
	}
	s.Screen.Show()
}

//func fitBufInto(buf []string, w int, h int, offs int) []string {
//	//var newBuf []string
//	//for
//}

func (s *State) writeLine(lineNum int, line string) {
	w, _ := s.Screen.Size()
	for x, r := range line {
		if x > w {
			return
		}
		s.Screen.SetContent(x, lineNum, r, nil, tcell.StyleDefault)
	}
}

func (s *State) moveCursorIntoBounds() {
	w, h := s.Screen.Size()
	if s.cursorY > h-1 {
		s.cursorY = h - 1
	}
	if s.cursorX > w-1 {
		s.cursorX = w - 1
	}
	s.showCursor()
}

func (s *State) bufferPosToViewPos(bufX, bufY int, wrap bool) (int, int) {
	w, _ := s.Screen.Size()
	x := bufX % w

	offs := 0
	for i := 0; i < bufY-s.scrollY; i++ {
		offs += 1 + len(s.buffer[s.scrollY+i])/w
	}

	y := bufX/w + offs

	return x, y
}
