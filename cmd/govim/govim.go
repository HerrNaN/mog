package main

import (
	"os"

	"elyria.io/govim/internal/govim"
)

func main() {
	if len(os.Args) == 1 {
		govim.NewProgram().Start()
	} else {
		govim.NewProgramAt(os.Args[1]).Start()
	}
}
