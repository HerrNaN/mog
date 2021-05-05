package govim

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gdamore/tcell/v2"
)

func TestSimpleFrame_MoveCursor(t *testing.T) {
	type fields struct {
		screen tcell.Screen
		buffer []string
		cursor Cursor
		offset int
	}
	type args struct {
		d dir
	}
	tests := []struct {
		name     string
		fields   fields
		d        dir
		expected fields
	}{
		{
			name: "cannot move SimpleCursor above first line",
			fields: fields{
				screen: tcell.NewSimulationScreen("UTF-8"),
				buffer: []string{""},
				cursor: NewSimpleCursorAt(0, 0),
				offset: 0,
			},
			d: dirUp,
			expected: fields{
				screen: tcell.NewSimulationScreen("UTF-8"),
				buffer: []string{""},
				cursor: NewSimpleCursorAt(0, 0),
				offset: 0,
			},
		},
		{
			name: "can move SimpleCursor up when there is another line",
			fields: fields{
				screen: tcell.NewSimulationScreen("UTF-8"),
				buffer: []string{"ab", "cd"},
				cursor: NewSimpleCursorAt(0, 1),
				offset: 0,
			},
			d: dirUp,
			expected: fields{
				screen: tcell.NewSimulationScreen("UTF-8"),
				buffer: []string{"ab", "cd"},
				cursor: NewSimpleCursorAt(0, 0),
				offset: 0,
			},
		},
		//{
		//	name: "cannot move SimpleCursor below last line",
		//	fields: fields{
		//		cursorX: 0,
		//		cursorY: 0,
		//		buffer:  []string{""},
		//	},
		//	d: dirDown,
		//	expected: fields{
		//		cursorX: 0,
		//		cursorY: 0,
		//		buffer:  []string{""},
		//	},
		//},
		//{
		//	name: "can move SimpleCursor down when there is another line",
		//	fields: fields{
		//		cursorX: 0,
		//		cursorY: 0,
		//		buffer:  []string{"ab", "cd"},
		//	},
		//	d: dirDown,
		//	expected: fields{
		//		cursorX: 0,
		//		cursorY: 1,
		//		buffer:  []string{"ab", "cd"},
		//	},
		//},
		//{
		//	name: "cannot move SimpleCursor to the left of the first character on a line",
		//	fields: fields{
		//		cursorX: 0,
		//		cursorY: 0,
		//		buffer:  []string{""},
		//	},
		//	d: dirLeft,
		//	expected: fields{
		//		cursorX: 0,
		//		cursorY: 0,
		//		buffer:  []string{""},
		//	},
		//},
		//{
		//	name: "can move SimpleCursor left when there is another character",
		//	fields: fields{
		//		cursorX: 1,
		//		cursorY: 0,
		//		buffer:  []string{"ab", "cd"},
		//	},
		//	d: dirLeft,
		//	expected: fields{
		//		cursorX: 0,
		//		cursorY: 0,
		//		buffer:  []string{"ab", "cd"},
		//	},
		//},
		//{
		//	name: "moves SimpleCursor the second last character when when moving from position outside of line length",
		//	fields: fields{
		//		cursorX: 10,
		//		cursorY: 0,
		//		buffer:  []string{"ab", "cd"},
		//	},
		//	d: dirLeft,
		//	expected: fields{
		//		cursorX: 1,
		//		cursorY: 0,
		//		buffer:  []string{"ab", "cd"},
		//	},
		//},
		//{
		//	name: "cannot move SimpleCursor to the right of the last character on a line",
		//	fields: fields{
		//		cursorX: 0,
		//		cursorY: 0,
		//		buffer:  []string{""},
		//	},
		//	d: dirRight,
		//	expected: fields{
		//		cursorX: 0,
		//		cursorY: 0,
		//		buffer:  []string{""},
		//	},
		//},
		//{
		//	name: "can move SimpleCursor right when there is another character",
		//	fields: fields{
		//		cursorX: 0,
		//		cursorY: 0,
		//		buffer:  []string{"ab", "cd"},
		//	},
		//	d: dirRight,
		//	expected: fields{
		//		cursorX: 1,
		//		cursorY: 0,
		//		buffer:  []string{"ab", "cd"},
		//	},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &SimpleFrame{
				screen: tt.fields.screen,
				buffer: tt.fields.buffer,
				cursor: tt.fields.cursor,
				offset: tt.fields.offset,
			}
			f.MoveCursor(tt.d)
			assert.Equal(t, f.screen, tt.expected.screen)
			assert.Equal(t, f.buffer, tt.expected.buffer)
			assert.Equal(t, f.cursor, tt.expected.cursor)
		})
	}
}
