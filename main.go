package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/fanyang01/mips"
)

type Mode int

const (
	noMode  Mode = 0
	asmMode Mode = 1 << iota
	disasmMode
	runMode
	asmRunMode
	stepRunMode
)

var (
	asmM     = flag.Bool("a", false, "Assemble")
	disasmM  = flag.Bool("d", false, "Disassemble")
	runM     = flag.Bool("r", false, "Run")
	asmRunM  = flag.Bool("R", false, "Assemble and run")
	stepRunM = flag.Bool("s", false, "Step run")
	outFile  = flag.String("o", "a.out", "Output file")
	logger   = log.New(os.Stderr, "", 0)
)

func main() {
	flag.Parse()
	mode := parseMode()
	switch mode {
	case asmMode:
		asmFile(flag.Arg(0))
	case disasmMode:
		disasmFile(flag.Arg(0))
	case runMode:
		runFile(flag.Arg(0))
	case asmRunMode:
		asmAndRun(flag.Arg(0))
	case stepRunMode:
		stepRun(flag.Arg(0))
	}
}

func asmFile(filename string) {
	f, err := os.Open(filename)
	checkFatalErr(err)
	defer f.Close()

	out, err := os.OpenFile(*outFile, os.O_WRONLY|os.O_CREATE, 0644)
	checkFatalErr(err)
	defer out.Close()
	w := bufio.NewWriter(out)
	defer w.Flush()

	assembler := mips.NewAssembler(f)
	s, err := assembler.Assemble()
	checkFatalErr(err)
	_, err = w.Write(s)
	checkFatalErr(err)
}

func disasmFile(filename string) {
	f, err := os.Open(filename)
	checkFatalErr(err)
	defer f.Close()

	disassembler := mips.NewDisassembler(f)
	s, err := disassembler.Disassemble()
	checkFatalErr(err)

	fmt.Println(string(s))
}

func runFile(filename string) {
	s, err := ioutil.ReadFile(filename)
	checkFatalErr(err)

	em := mips.NewEmulator()
	err = em.LoadAndRun(s)
	checkFatalErr(err)

	err = em.Wait()
	checkFatalErr(err)
}

func asmAndRun(filename string) {
	f, err := os.Open(filename)
	checkFatalErr(err)
	defer f.Close()

	assembler := mips.NewAssembler(f)
	s, err := assembler.Assemble()
	checkFatalErr(err)

	em := mips.NewEmulator()
	err = em.LoadAndRun(s)
	checkFatalErr(err)

	err = em.Wait()
	checkFatalErr(err)
}

func parseMode() Mode {
	mode := noMode
	if *asmM {
		mode |= asmMode
	}
	if *disasmM {
		mode |= disasmMode
	}
	if *runM {
		mode |= runMode
	}
	if *asmRunM {
		mode |= asmRunMode
	}
	if *stepRunM {
		mode |= stepRunMode
	}
	// If no flag is specified, enter assembler mode
	if mode == noMode {
		mode = asmMode
	}
	switch mode {
	case asmMode, disasmMode, runMode, asmRunMode, stepRunMode:
		if flag.NArg() < 1 {
			logger.Fatal("Please specify a file to process")
		}
	default:
		flag.PrintDefaults()
	}
	return mode
}

func checkFatalErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
