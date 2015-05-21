package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/fanyang01/mips"
)

type Cmd int

const (
	cmdEmpty Cmd = 0
	cmdError Cmd = 1 << iota
	cmdStep
	cmdListSrc
	cmdMemory
	cmdReg
	cmdRun
	cmdQuit
)

var (
	registers = []string{
		"zero", "at",
		"v0", "v1",
		"a0", "a1", "a2", "a3",
		"t0", "t1", "t2", "t3", "t4", "t5", "t6", "t7",
		"s0", "s1", "s2", "s3", "s4", "s5", "s6", "s7",
		"t8", "t9",
		"k0", "k1",
		"gp", "sp", "fp", "ra",
		"PC", "HI", "LO",
	}
)

type Command struct {
	cmd  Cmd
	args interface{}
}

func stepRun(filename string) {
	s, err := ioutil.ReadFile(filename)
	checkFatalErr(err)

	em := mips.NewEmulator()
	err = em.LoadAndStart(s)
	checkFatalErr(err)

	var cache *Command
	for {
		fmt.Printf("(mips) ")
		cmd := scanCommand()
	LABEL:
		switch cmd.cmd {
		case cmdError:
			continue
		case cmdEmpty:
			if cache != nil {
				cmd = *cache
				goto LABEL
			}
			fmt.Fprintln(os.Stderr, "Please specify a command")
			continue
		case cmdListSrc:
			listSrc(em, cmd.args.([]int))
		case cmdStep:
			step(em, cmd.args.([]int))
		case cmdMemory:
			showMemory(em, cmd.args.([]int))
		case cmdReg:
			showReg(em, cmd.args.([]string))
		case cmdRun:
			runToEnd(em)
		case cmdQuit:
			return
		}
		cache = &cmd
	}
}

func step(em *mips.Emulator, args []int) {
	defer func() {
		if err := recover(); err != nil {
			logger.Println(err)
		}
	}()
	count := 1
	if len(args) > 0 && args[0] > 1 {
		count = args[0]
	}
	for i := 0; i < count; i++ {
		s, err := em.FetchSource(1)
		checkErr(err)
		fmt.Println(string(s))
		err = em.Step()
		checkErr(err)
	}
}

func showMemory(em *mips.Emulator, args []int) {
	defer func() {
		if err := recover(); err != nil {
			logger.Println(err)
		}
	}()
	if len(args) < 1 {
		panic("Please specify at least one address")
	}
	for _, addr := range args {
		word, err := em.ShowMemory(addr)
		checkErr(err)
		fmt.Printf("%#0x: %#x(%d)\n", addr, word, word)
	}
}

func showReg(em *mips.Emulator, args []string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Println(err)
		}
	}()
	if len(args) == 0 {
		for _, reg := range registers {
			word, err := em.ShowReg(reg)
			checkErr(err)
			fmt.Printf("%s:\t%#0x(%d)\n", reg, word, word)
		}
		return
	}
	for _, reg := range args {
		word, err := em.ShowReg(reg)
		checkErr(err)
		fmt.Printf("%s:\t%#0x(%d)\n", reg, word, word)
	}
}

func runToEnd(em *mips.Emulator) {
	err := em.Step()
	for ; err == nil; err = em.Step() {
	}
	logger.Println(err)
}

func listSrc(em *mips.Emulator, args []int) {
	defer func() {
		if err := recover(); err != nil {
			logger.Println(err)
		}
	}()
	count := 1
	if len(args) > 0 && args[0] > 1 {
		count = args[0]
	}
	s, err := em.FetchSource(count)
	checkErr(err)
	fmt.Println(string(s))
	checkErr(err)
}

func scanCommand() (cmd Command) {
	defer func() {
		if err := recover(); err != nil {
			cmd.cmd = cmdError
			logger.Println(err)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	s, err := reader.ReadString('\n')
	checkErr(err)
	s = strings.TrimSpace(s)
	tokens := strings.Fields(s)
	if len(tokens) == 0 {
		cmd.cmd = cmdEmpty
		return
	}
	switch tokens[0] {
	case "l", "list":
		cmd.cmd = cmdListSrc
	case "s", "step":
		cmd.cmd = cmdStep
	case "r", "reg":
		cmd.cmd = cmdReg
	case "m", "mem", "memory":
		cmd.cmd = cmdMemory
	case "run":
		cmd.cmd = cmdRun
	case "q", "quit":
		cmd.cmd = cmdQuit
	default:
		panic(s + ": Invalid command")
	}

	args := tokens[1:]
	switch cmd.cmd {
	case cmdMemory, cmdStep, cmdListSrc:
		cmd.args = []int{}
		for _, a := range args {
			n, err := strconv.ParseInt(a, 0, 32)
			checkErr(err)
			cmd.args = append(cmd.args.([]int), int(n))
		}
	default:
		cmd.args = args
	}
	return
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
