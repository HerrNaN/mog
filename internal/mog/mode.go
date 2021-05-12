package mog

type Mode struct {
	Name      string
	ShortName string
	Letter    rune
}

const (
	ModeLetterLength    = 1
	ModeShortNameLength = 3
)

var (
	ModeInsert = Mode{
		Name:      "Insert",
		ShortName: "Ins",
		Letter:    'I',
	}
)
