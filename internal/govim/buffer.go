package govim

import (
	"strings"
)

type Buffer struct {
	cursorX, cursorY int

	// The buffer containing the contents currently being processed
	buffer []string

	// The line number from the buffer that should be written to the
	// first line on screen
	offset int
}

func EmptyBuffer() *Buffer {
	return &Buffer{}
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
