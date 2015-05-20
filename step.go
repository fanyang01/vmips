package main

import (
	"os"

	"github.com/fanyang01/mips"
)

func stepRun(filename string) {
	f, err := os.Open(filename)
	checkFatalErr(err)
	defer f.Close()

	assembler := mips.NewAssembler(f)
	s, err := assembler.Assemble()
	checkFatalErr(err)

	em := mips.NewEmulator()
	err = em.LoadAndStart(s)
	checkFatalErr(err)
}
