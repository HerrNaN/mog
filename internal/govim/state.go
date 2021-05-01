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
	buffer string

	// The line number from the buffer that should be written to the
	// first line on screen
	scrollIdx int
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
		buffer:  "Hello\nMy\nName\nIs\n  HAHAH\n\n\n\n\nasdasdasd\nasd\nasd\nasd",
	}
}

func (s *State) moveCursor(d dir) {
	w, h := s.Screen.Size()
	switch d {
	case dirUp:
		if s.cursorY > 0 {
			s.cursorY--
		} else if s.scrollIdx > 0 {
			s.scrollIdx--
		}
	case dirDown:
		if s.scrollIdx+s.cursorY < len(strings.Split(s.buffer, "\n")) {
			if s.cursorY == h-1 {
				s.scrollIdx++
			} else {
				s.cursorY++
			}
		}
	case dirLeft:
		if s.cursorX > 0 {
			s.cursorX--
		}
	case dirRight:
		if s.cursorX < w-1 {
			s.cursorX++
		}
	}
	s.Screen.ShowCursor(s.cursorX, s.cursorY)
}

func (s *State) showCursor() {
	s.Screen.ShowCursor(s.cursorX, s.cursorY)
}

// Start sets up the program and starts the event loop
func (s *State) Start() {
	s.Screen.Clear()
	s.Screen.ShowCursor(s.cursorX, s.cursorY)
	s.run()
}

func (s *State) run() {
	for {
		// Update screen
		//s.Screen.Show()
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
	w, h := s.Screen.Size()
	lineN := 0
	var toShow []string
	buf := strings.Split(s.buffer, "\n")
	for _, l := range buf[s.scrollIdx:] {
		if lineN >= h {
			break
		}
		if len(l) < w {
			toShow = append(toShow, l)
			lineN++
		} else {
			chunk := l
			for {
				if len(chunk) < w {
					toShow = append(toShow, chunk)
					lineN++
					break
				}
				toShow = append(toShow, chunk[:w])
				chunk = chunk[w:]
				lineN++
			}
		}
	}

	for i, l := range toShow {
		s.writeLine(i, l)
	}

	// All lines after end of buffer are marked with a ~ in the left margin
	for i := 0; i < h-len(toShow); i++ {
		s.Screen.SetContent(0, len(toShow)+i+1, '~', nil, tcell.StyleDefault)
	}
	s.Screen.Show()
}

func (s *State) writeLine(lineNum int, line string) {
	for x, r := range line {
		s.Screen.SetContent(x, lineNum, r, nil, tcell.StyleDefault)
	}
}

// Quit exits the program
func (s *State) Quit() {
	s.Screen.Fini()
	os.Exit(0)
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
