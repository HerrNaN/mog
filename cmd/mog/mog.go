package main

import (
	"log"
	"os"

	"elyria.io/govim/internal/mog"
)

func main() {
	if len(os.Args) == 1 {
		mog.NewProgram().Start()
	} else {
		p, err := mog.NewProgramAt(os.Args[1])
		if err != nil {
			log.Fatalf("%+v", err)
		}
		p.Start()
	}
}
