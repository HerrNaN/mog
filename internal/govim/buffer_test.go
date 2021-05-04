package govim

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuffer_MoveCursor(t *testing.T) {
	type fields struct {
		cursorX int
		cursorY int
		buffer  []string
	}
	tests := []struct {
		name     string
		fields   fields
		d        dir
		expected fields
	}{
		{
			name: "cannot move cursor above first line",
			fields: fields{
				cursorX: 0,
				cursorY: 0,
				buffer:  []string{""},
			},
			d: dirUp,
			expected: fields{
				cursorX: 0,
				cursorY: 0,
				buffer:  []string{""},
			},
		},
		{
			name: "cannot move cursor below last line",
			fields: fields{
				cursorX: 0,
				cursorY: 0,
				buffer:  []string{""},
			},
			d: dirDown,
			expected: fields{
				cursorX: 0,
				cursorY: 0,
				buffer:  []string{""},
			},
		},
		{
			name: "cannot move cursor to the left of the first character on a line",
			fields: fields{
				cursorX: 0,
				cursorY: 0,
				buffer:  []string{""},
			},
			d: dirLeft,
			expected: fields{
				cursorX: 0,
				cursorY: 0,
				buffer:  []string{""},
			},
		},
		{
			name: "cannot move cursor to the right of the last character on a line",
			fields: fields{
				cursorX: 0,
				cursorY: 0,
				buffer:  []string{""},
			},
			d: dirRight,
			expected: fields{
				cursorX: 0,
				cursorY: 0,
				buffer:  []string{""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Buffer{
				cursorX: tt.fields.cursorX,
				cursorY: tt.fields.cursorY,
				buffer:  tt.fields.buffer,
			}
			b.MoveCursor(tt.d)
			assert.Equal(t, tt.expected.cursorY, b.cursorY)
			assert.Equal(t, tt.expected.cursorX, b.cursorX)
			assert.Equal(t, tt.expected.buffer, b.buffer)
		})
	}
}
