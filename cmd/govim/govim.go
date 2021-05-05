package main

import (
	"log"
	"os"

	"elyria.io/govim/internal/govim"
)

func main() {
	if len(os.Args) == 1 {
		govim.NewProgram().Start()
	} else {
		p, err := govim.NewProgramAt(os.Args[1])
		if err != nil {
			log.Fatalf("%+v", err)
		}
		p.Start()
	}
}
