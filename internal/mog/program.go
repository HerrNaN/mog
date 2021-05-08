package mog

import (
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

type Frame interface {
	MoveCursor(dir)
	Show()
	PollEvent() tcell.Event
	Sync()
	Close()

	// InsertRune inserts a rune at the current cursor position.
	// Inserting a run 'a' into a line 'bb' at position 0 would
	// yield 'abb'.
	InsertRune(r rune)

}

type Program struct {
	frame Frame
}

type dir int

const (
	dirUp dir = iota
	dirDown
	dirRight
	dirLeft
)

// NewProgram creates a new program with an empty buffer
func NewProgram() *Program {
	return &Program{
		frame: EmptyFrame(),
	}
}

func NewProgramAt(filename string) (*Program, error) {
	bs, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("%+v", err)
		return nil, err
	}
	return &Program{
		frame: NewFrame(bs),
	}, nil
}

// Start sets up the program and starts the event loop
func (p *Program) Start() {
	p.run()
}

// Quit exits the program
func (p *Program) Quit() {
	p.frame.Close()
	os.Exit(0)
}

func (p *Program) Show() {
	p.frame.Show()
}

func (p *Program) run() {
	for {
		// Update screen
		p.Show()

		// Poll event
		ev := p.frame.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			p.frame.Sync()
			p.frame.Show()
		case *tcell.EventKey:
			p.handleEventKey(*ev)
		default:
			log.Print(ev)
		}
	}
}

func (p *Program) handleEventKey(ev tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEscape:
		p.Quit()
	case tcell.KeyUp:
		p.frame.MoveCursor(dirUp)
	case tcell.KeyDown:
		p.frame.MoveCursor(dirDown)
	case tcell.KeyRight:
		p.frame.MoveCursor(dirRight)
	case tcell.KeyLeft:
		p.frame.MoveCursor(dirLeft)
	case tcell.KeyRune:
		p.handleEventRune(ev.Rune())
	}
}

func (p *Program) handleEventRune(r rune) {
	p.frame.InsertRune(r)
	p.frame.MoveCursor(dirRight)
}
