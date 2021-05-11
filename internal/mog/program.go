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

	// HandleEvent returns true if the Frame closed as a result of this event
	// and otherwise false.
	HandleEvent(tcell.Event) bool
	Close() error

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

func NewProgramFromFile(filename string) (*Program, error) {
	return &Program{
		frame: NewFrameFromFile(filename),
	}, nil
}

// Start sets up the program and starts the event loop
func (p *Program) Start() {
	p.run()
}

// Quit exits the program
func (p *Program) Quit() {
	err := p.frame.Close()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	os.Exit(0)
}

func (p *Program) Show() {
	p.frame.Show()
}

func (p *Program) run() {
	for {
		p.Show()
		ev := p.frame.PollEvent()

		if closed := p.frame.HandleEvent(ev); closed {
			break
		}
	}
}
