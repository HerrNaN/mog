package govim

import (
	"strings"
)

type Buffer struct {
	cursorX, cursorY int

	// The buffer containing the contents currently being processed
	buffer []string
}

func EmptyBuffer() *Buffer {
	return &Buffer{
		buffer: []string{""},
	}
}

func BufferFromRaw(raw []string) *Buffer {
	return &Buffer{
		buffer: raw,
	}
}

func BufferFrom(s string) *Buffer {
	buf := strings.Split(s, "\n")
	return BufferFromRaw(buf)
}

func (b *Buffer) MoveCursor(d dir) {
	switch d {
	case dirUp:
		if b.cursorY > 0 {
			b.cursorY--
		}
	case dirDown:
		if b.cursorY < len(b.buffer)-1 {
			b.cursorY++
		}
	case dirLeft:
		if b.cursorX > len(b.buffer[b.cursorY]) {
			b.cursorX = len(b.buffer[b.cursorY])
		}
		if b.cursorX > 0 {
			b.cursorX--
		}
	case dirRight:
		if b.cursorX < len(b.buffer[b.cursorY]) {
			b.cursorX++
		}
	}
}
