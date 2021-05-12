package mog

import (
	"fmt"
	"log"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gdamore/tcell/v2"
)

func TestSimpleFrame_MoveCursor(t *testing.T) {
	type fields struct {
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
		args                  args
		expectedScreenCursorX int
		expectedScreenCursorY int
	}{
		{
			name: "cannot move SimpleCursor above first line",
			fields: fields{
				buffer: []string{""},
				cursor: NewSimpleCursorAt(0, 0),
				offset: 0,
			},
			args:                  args{dirUp},
			expectedScreenCursorX: 0,
			expectedScreenCursorY: 0,
		},
		{
			name: "can move SimpleCursor up when there is another line",
			fields: fields{
				buffer: []string{"ab", "cd"},
				cursor: NewSimpleCursorAt(0, 1),
				offset: 0,
			},
			args:                  args{dirUp},
			expectedScreenCursorX: 0,
			expectedScreenCursorY: 0,
		},
		{
			name: "cannot move SimpleCursor below last line",
			fields: fields{
				buffer: []string{""},
				cursor: NewSimpleCursorAt(0, 0),
				offset: 0,
			},
			args:                  args{dirDown},
			expectedScreenCursorX: 0,
			expectedScreenCursorY: 0,
		},
		{
			name: "can move SimpleCursor down when there is another line",
			fields: fields{
				buffer: []string{"ab", "cd"},
				cursor: NewSimpleCursorAt(0, 0),
				offset: 0,
			},
			args:                  args{dirDown},
			expectedScreenCursorX: 0,
			expectedScreenCursorY: 1,
		},
		{
			name: "cannot move SimpleCursor to the left of the first character on a line",
			fields: fields{
				buffer: []string{""},
				cursor: NewSimpleCursorAt(0, 0),
				offset: 0,
			},
			args:                  args{dirLeft},
			expectedScreenCursorX: 0,
			expectedScreenCursorY: 0,
		},
		{
			name: "can move SimpleCursor left when there is another character",
			fields: fields{
				buffer: []string{"ab", "cd"},
				cursor: NewSimpleCursorAt(1, 0),
				offset: 0,
			},
			args:                  args{dirLeft},
			expectedScreenCursorX: 0,
			expectedScreenCursorY: 0,
		},
		{
			name: "moves SimpleCursor the second last character when when moving from position outside of line length",
			fields: fields{
				buffer: []string{"ab", "cd"},
				cursor: NewSimpleCursorAt(10, 0),
				offset: 0,
			},
			args:                  args{dirLeft},
			expectedScreenCursorX: 0,
			expectedScreenCursorY: 0,
		},
		{
			name: "cannot move SimpleCursor to the right of the last character on a line",
			fields: fields{
				buffer: []string{""},
				cursor: NewSimpleCursorAt(0, 0),
				offset: 0,
			},
			args:                  args{dirRight},
			expectedScreenCursorX: 0,
			expectedScreenCursorY: 0,
		},
		{
			name: "can move SimpleCursor right when there is another character",
			fields: fields{
				buffer: []string{"ab", "cd"},
				cursor: NewSimpleCursorAt(0, 0),
				offset: 0,
			},
			args:                  args{dirRight},
			expectedScreenCursorX: 1,
			expectedScreenCursorY: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			simulationScreen := tcell.NewSimulationScreen("UTF-8")
			simulationScreen.SetSize(5, 5)
			f := &SimpleFrame{
				screen: simulationScreen,
				buffer: tt.fields.buffer,
				cursor: tt.fields.cursor,
				offset: tt.fields.offset,
			}
			f.MoveCursor(tt.args.d)

			ss := f.screen.(tcell.SimulationScreen)
			cx, cy, vis := ss.GetCursor()
			assert.True(t, vis)
			assert.Equal(t, tt.expectedScreenCursorX, cx)
			assert.Equal(t, tt.expectedScreenCursorY, cy)
		})
	}
}

func TestSimpleFrame_writeBufferToScreen(t *testing.T) {
	type fields struct {
		buffer []string
		cursor Cursor
		offset int
	}
	tests := []struct {
		name                      string
		fields                    fields
		expectedContents          []string
		screenWidth, screenHeight int
	}{
		{
			name: "displays empty line + tildes + default mode with empty buffer",
			fields: fields{
				buffer: []string{""},
				cursor: NewSimpleCursor(),
				offset: 0,
			},
			screenHeight:     3,
			screenWidth:      15,
			expectedContents: []string{"               ", "~              ", " -- Insert --  "},
		},
		{
			name: "display wrapped line when buffer line is longer than screen width",
			fields: fields{
				buffer: []string{"abcdefghijabcdefghij"},
				cursor: NewSimpleCursor(),
				offset: 0,
			},
			screenHeight:     3,
			screenWidth:      15,
			expectedContents: []string{"abcdefghijabcde", "fghij          ", " -- Insert --  "},
		},
		{
			name: "display only line 1 and 2 (0 indexed) + 1 tilde line when offset is 1",
			fields: fields{
				buffer: []string{"a", "b", "c"},
				cursor: NewSimpleCursor(),
				offset: 1,
			},
			screenHeight:     3,
			screenWidth:      15,
			expectedContents: []string{"b              ", "c              ", " -- Insert --  "},
		},
		{
			name: "display only mode on last line",
			fields: fields{
				buffer: []string{"", "", "abcdefghihjkmno"},
				cursor: NewSimpleCursor(),
				offset: 0,
			},
			screenHeight:     3,
			screenWidth:      15,
			expectedContents: []string{"               ", "               ", " -- Insert --  "},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			simulationScreen := tcell.NewSimulationScreen("UTF-8")
			err := simulationScreen.Init()
			assert.Nil(t, err)
			simulationScreen.SetSize(tt.screenWidth, tt.screenHeight)

			f := &SimpleFrame{
				screen: simulationScreen,
				buffer: tt.fields.buffer,
				cursor: tt.fields.cursor,
				offset: tt.fields.offset,
				mode:   ModeInsert,
			}
			f.writeBufferToScreen()

			// To see what has been written to the contents of the screen we
			// also have to show it.
			f.screen.Show()

			// Read actual screen contents
			ss := f.screen.(tcell.SimulationScreen)
			cells, w, h := ss.GetContents()
			var actualBuffer []string
			for i := 0; i < h; i++ {
				var line []rune
				for j := 0; j < w; j++ {
					line = append(line, cells[i*w+j].Runes[0])
				}
				actualBuffer = append(actualBuffer, string(line))
			}

			assert.EqualValues(t, tt.expectedContents, actualBuffer)
		})
	}
}

func TestSimpleFrame_InsertRune(t *testing.T) {
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
	type args struct {
		r rune
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		expectedBuffer []string
	}{
		{
			name: "adds rune at cursor position",
			fields: fields{
				screen: simulationScreen,
				buffer: []string{"bb"},
				cursor: NewSimpleCursorAt(0, 0),
				offset: 0,
			},
			args:           args{r: 'a'},
			expectedBuffer: []string{"abb"},
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
			f.InsertRune(tt.args.r)
			assert.EqualValues(t, tt.expectedBuffer, f.buffer)
		})
	}
}

func Test_loadFile(t *testing.T) {
	filename := "TestNewFrameFromFile.txt"
	lockFileName := lockFilePathOf(filename)
	file, err := os.Create(filename)
	assert.Nil(t, err)
	defer func() {
		err := os.Remove(filename)
		if err != nil {
			log.Println(fmt.Errorf("couldn't remove file at end of test: %w", err))
		}
	}()
	err = file.Close()
	assert.Nil(t, err)

	f := &SimpleFrame{
		screen:       tcell.NewSimulationScreen("UTF-8"),
		buffer:       nil,
		cursor:       nil,
		offset:       0,
		filePath:     "",
		lockFilePath: "",
	}
	err = f.loadFile(filename)
	assert.Nil(t, err)
	assert.EqualValues(t, lockFileName, f.lockFilePath)
	assert.EqualValues(t, filename, f.filePath)

	info, err := os.Stat(lockFileName)
	assert.Nil(t, err)
	assert.Equal(t, false, info.IsDir())
	assert.EqualValues(t, 0, info.Size())
	assert.Equal(t, path.Base(lockFileName), info.Name())

	err = f.Close()
	assert.Nil(t, err)
	_, err = os.Stat(lockFileName)
	assert.NotNil(t, err)
}

func Test_loadFile_ReturnsErrorWhenLockFileExists(t *testing.T) {
	filename := "TestNewFrameFromFile.txt"

	lockFileName := lockFilePathOf(filename)
	file, err := os.Create(lockFileName)
	assert.Nil(t, err)
	defer func() {
		err := os.Remove(lockFileName)
		if err != nil {
			log.Println(fmt.Errorf("couldn't remove file at end of test: %w", err))
		}
	}()
	err = file.Close()
	assert.Nil(t, err)
	f := &SimpleFrame{
		screen:       tcell.NewSimulationScreen("UTF-8"),
		buffer:       nil,
		cursor:       nil,
		offset:       0,
		filePath:     "",
		lockFilePath: "",
	}

	err = f.loadFile(filename)
	assert.Error(t, err)
}
