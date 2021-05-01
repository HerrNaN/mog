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

type Dir int

const (
	DirUp Dir = iota
	DirDown
	DirRight
	DirLeft
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

func (s *State) moveCursor(d Dir) {
	w, h := s.Screen.Size()
	switch d {
	case DirUp:
		if s.cursorY > 0 {
			s.cursorY--
		}
	case DirDown:
		if s.cursorY < h-1 {
			s.cursorY++
		}
	case DirLeft:
		if s.cursorX > 0 {
			s.cursorX--
		}
	case DirRight:
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
	case tcell.KeyEscape:
		s.Quit()
	case tcell.KeyUp:
		s.moveCursor(DirUp)
	case tcell.KeyDown:
		s.moveCursor(DirDown)
	case tcell.KeyRight:
		s.moveCursor(DirRight)
	case tcell.KeyLeft:
		s.moveCursor(DirLeft)
	}
}

func (s *State) Quit() {
	s.Screen.Fini()
	os.Exit(0)
}
