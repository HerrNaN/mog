package mog

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gdamore/tcell/v2"
)

func TestSimpleFrame_MoveCursor(t *testing.T) {
	simulationScreen := tcell.NewSimulationScreen("UTF-8")
	simulationScreen.SetSize(5, 5)
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
				screen: simulationScreen,
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
				screen: simulationScreen,
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
				screen: simulationScreen,
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
				screen: simulationScreen,
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
				screen: simulationScreen,
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
				screen: simulationScreen,
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
				screen: simulationScreen,
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
				screen: simulationScreen,
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
				screen: simulationScreen,
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

func TestSimpleFrame_writeBufferToScreen(t *testing.T) {
	simulationScreen := tcell.NewSimulationScreen("UTF-8")
	err := simulationScreen.Init()
	assert.Nil(t, err)
	simulationScreen.SetSize(3, 3)
	type fields struct {
		screen tcell.Screen
		buffer []string
		cursor Cursor
		offset int
	}
	tests := []struct {
		name             string
		fields           fields
		expectedContents []string
	}{
		{
			name: "displays empty line + tildes on the rest of the lines with empty buffer",
			fields: fields{
				screen: simulationScreen,
				buffer: []string{""},
				cursor: NewSimpleCursor(),
				offset: 0,
			},
			expectedContents: []string{"   ", "~  ", "~  "},
		},
		{
			name: "display wrapped line when buffer line is longer than screen width",
			fields: fields{
				screen: simulationScreen,
				buffer: []string{"abcde"},
				cursor: NewSimpleCursor(),
				offset: 0,
			},
			expectedContents: []string{"abc", "de ", "~  "},
		},
		{
			name: "display only row 1 and 2 (0 indexed) when offset is 1",
			fields: fields{
				screen: simulationScreen,
				buffer: []string{"a", "b", "c"},
				cursor: NewSimpleCursor(),
				offset: 1,
			},
			expectedContents: []string{"b  ", "c  ", "~  "},
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
			f.writeBufferToScreen()

			// To see what has been written to the contents of the screen we
			// also have to show it.
			f.screen.Show()

			expectedRunes := []rune(strings.Join(tt.expectedContents, ""))

			ss := f.screen.(tcell.SimulationScreen)
			cells, w, h := ss.GetContents()
			var actualRunes []rune
			for i := 0; i < len(cells); i++ {
				actualRunes = append(actualRunes, cells[i].Runes[0])
			}

			assert.Equal(t, 3, w)
			assert.Equal(t, 3, h)
			assert.EqualValues(t, expectedRunes, actualRunes)
		})
	}
}
