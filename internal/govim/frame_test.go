package govim

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gdamore/tcell/v2"
)

func TestSimpleFrame_MoveCursor(t *testing.T) {
	ss := tcell.NewSimulationScreen("UTF-8")
	ss.SetSize(5, 5)
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
		name                  string
		fields                fields
		d                     dir
		expectedScreenCursorX int
		expectedScreenCursorY int
	}{
		{
			name: "cannot move SimpleCursor above first line",
			fields: fields{
				screen: ss,
				buffer: []string{""},
				cursor: NewSimpleCursorAt(0, 0),
				offset: 0,
			},
			d:                     dirUp,
			expectedScreenCursorX: 0,
			expectedScreenCursorY: 0,
		},
		{
			name: "can move SimpleCursor up when there is another line",
			fields: fields{
				screen: ss,
				buffer: []string{"ab", "cd"},
				cursor: NewSimpleCursorAt(0, 1),
				offset: 0,
			},
			d:                     dirUp,
			expectedScreenCursorX: 0,
			expectedScreenCursorY: 0,
		},
		{
			name: "cannot move SimpleCursor below last line",
			fields: fields{
				screen: ss,
				buffer: []string{""},
				cursor: NewSimpleCursorAt(0, 0),
				offset: 0,
			},
			d:                     dirDown,
			expectedScreenCursorX: 0,
			expectedScreenCursorY: 0,
		},
		{
			name: "can move SimpleCursor down when there is another line",
			fields: fields{
				screen: ss,
				buffer: []string{"ab", "cd"},
				cursor: NewSimpleCursorAt(0, 0),
				offset: 0,
			},
			d:                     dirDown,
			expectedScreenCursorX: 0,
			expectedScreenCursorY: 1,
		},
		{
			name: "cannot move SimpleCursor to the left of the first character on a line",
			fields: fields{
				screen: ss,
				buffer: []string{""},
				cursor: NewSimpleCursorAt(0, 0),
				offset: 0,
			},
			d:                     dirLeft,
			expectedScreenCursorX: 0,
			expectedScreenCursorY: 0,
		},
		{
			name: "can move SimpleCursor left when there is another character",
			fields: fields{
				screen: ss,
				buffer: []string{"ab", "cd"},
				cursor: NewSimpleCursorAt(1, 0),
				offset: 0,
			},
			d:                     dirLeft,
			expectedScreenCursorX: 0,
			expectedScreenCursorY: 0,
		},
		{
			name: "moves SimpleCursor the second last character when when moving from position outside of line length",
			fields: fields{
				screen: ss,
				buffer: []string{"ab", "cd"},
				cursor: NewSimpleCursorAt(10, 0),
				offset: 0,
			},
			d:                     dirLeft,
			expectedScreenCursorX: 0,
			expectedScreenCursorY: 0,
		},
		{
			name: "cannot move SimpleCursor to the right of the last character on a line",
			fields: fields{
				screen: ss,
				buffer: []string{""},
				cursor: NewSimpleCursorAt(0, 0),
				offset: 0,
			},
			d:                     dirRight,
			expectedScreenCursorX: 0,
			expectedScreenCursorY: 0,
		},
		{
			name: "can move SimpleCursor right when there is another character",
			fields: fields{
				screen: ss,
				buffer: []string{"ab", "cd"},
				cursor: NewSimpleCursorAt(0, 0),
				offset: 0,
			},
			d:                     dirRight,
			expectedScreenCursorX: 1,
			expectedScreenCursorY: 0,
		},
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

			ss := f.screen.(tcell.SimulationScreen)
			cx, cy, vis := ss.GetCursor()
			assert.True(t, vis)
			assert.Equal(t, tt.expectedScreenCursorX, cx)
			assert.Equal(t, tt.expectedScreenCursorY, cy)
		})
	}
}
