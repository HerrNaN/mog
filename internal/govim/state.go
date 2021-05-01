package govim

import (
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

type State struct {
	Screen           tcell.Screen
	cursorX, cursorY int
}

type dir int

const (
	dirUp dir = iota
	dirDown
	dirRight
	dirLeft
)

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
		Screen: s,
	}
}

func (s *State) moveCursor(d dir) {
	w, h := s.Screen.Size()
	switch d {
	case dirUp:
		if s.cursorY > 0 {
			s.cursorY--
		}
	case dirDown:
		if s.cursorY < h-1 {
			s.cursorY++
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

func (s *State) Start() {
	s.Screen.Clear()
	s.Screen.ShowCursor(s.cursorX, s.cursorY)
	s.run()
}

func (s *State) run() {
	for {
		// Update screen
		s.Screen.Show()

		// Poll event
		ev := s.Screen.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Screen.Sync()
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

func (s *State) Quit() {
	s.Screen.Fini()
	os.Exit(0)
}
