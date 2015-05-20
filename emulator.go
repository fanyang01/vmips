package mips

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
)

// ExitStatus is the exit status of virtual machine
//go:generate stringer -type=ExitStatus
type ExitStatus int

const (
	EXIT_NORMAL ExitStatus = iota
	EXIT_INT
	EXIT_EOF
	EXIT_ERROR
)

// Emulator runs machine code virtually
type Emulator struct {
	machine *Machine
	inst    chan execInst
	next    chan int
	exit    chan ExitStatus
	err     chan error
}

type execInst struct {
	name   string
	f      instFunc
	args   []int // function arguments
	branch bool  // Is branch instruction?
	eof    bool
}

func NewEmulator() *Emulator {
	return &Emulator{
		machine: NewMachine(),
		inst:    make(chan execInst),
		next:    make(chan int),
		exit:    make(chan ExitStatus),
		err:     make(chan error),
	}
}

func (e *Emulator) LoadAndRun(raw []byte) error {
	err := e.Load(raw)
	if err != nil {
		return err
	}
	e.Run()
	return nil
}

func (e *Emulator) Run() {
	go func() {
		for {
			e.step()
		}
	}()
	go func() {
		for {
			<-e.next
			e.fetch()
		}
	}()
	e.next <- 1
}

func (e *Emulator) Wait() error {
	status := <-e.exit
	switch status {
	case EXIT_ERROR:
		return <-e.err
	case EXIT_NORMAL, EXIT_EOF:
		return nil
	default:
		return nil
	}
}

func (e *Emulator) LoadAndStart(raw []byte) error {
	err := e.Load(raw)
	if err != nil {
		return err
	}
	e.Start()
	return nil
}

func (e *Emulator) Start() {
	go func() {
		for {
			e.step()
		}
	}()
	e.next <- 1
}

func (e *Emulator) Step() {
	<-e.next
	e.fetch()
}

func (e *Emulator) Exit() {
	e.exit <- EXIT_INT
}

func (e *Emulator) step() {
	defer func() {
		if err := recover(); err != nil {
			e.err <- fmt.Errorf("%v", err)
			e.exit <- EXIT_ERROR
		}
	}()
	inst := <-e.inst
	if inst.eof {
		e.exit <- EXIT_EOF
		return
	}
	inst.f(e.machine, inst.args...)
	if e.machine.exit {
		e.exit <- EXIT_NORMAL
		return
	}
	if !inst.branch {
		e.machine.r.PC = e.machine.r.PC + 4
	}
	e.next <- 1
}

// StartOnline loads emulator with empty code
func (e *Emulator) StartOnline() error {
	err := e.Load([]byte("text:0,data:0,main:0\n"))
	if err != nil {
		return err
	}
	go func() {
		for {
			e.step()
		}
	}()
	e.next <- 1
	return nil
}

// FetchOnline fetches 4 bytes as an instruction
func (e *Emulator) FetchOnline(raw []byte) {
	if len(raw) != 4 {
		e.exit <- EXIT_ERROR
		e.err <- errors.New("bad machine code")
		return
	}
	addr := e.machine.r.PC
	err := e.machine.m.writeBytes(addr, raw)
	if err != nil {
		e.exit <- EXIT_ERROR
		e.err <- fmt.Errorf("%v", err)
		return
	}
	e.fetch()
}

// Load loads object codes into emulator
func (e *Emulator) Load(code []byte) error {
	i := bytes.IndexByte(code, '\n')
	if i < 0 {
		return errors.New("load code: no header")
	}
	if !headerPattern.Match(code[:i]) {
		return errors.New("load code: invalid header")
	}
	sub := headerPattern.FindStringSubmatch(string(code[:i]))
	code = code[i+1:]

	text, err := strconv.ParseInt(sub[1], 10, 32)
	if err != nil {
		panic(err)
	}
	// text segment can be empty, for online loading
	if int(text) > len(code) {
		return errors.New("load code: text offset out of range")
	}
	data, err := strconv.ParseInt(sub[2], 10, 32)
	if err != nil {
		panic(err)
	}
	// data segment can be empty
	if data < text || int(data) > len(code) {
		return errors.New("load code: data offset out of range")
	}

	err = e.machine.m.writeBytes(TEXT_ADDRESS, code[text:data])
	if err != nil {
		return err
	}
	e.machine.m.writeBytes(DATA_ADDRESS, code[data:])
	if err != nil {
		return err
	}

	main, err := strconv.ParseInt(sub[3], 10, 32)
	if err != nil {
		panic(err)
	}
	if main < text || main >= data {
		return errors.New("load code: main offset out of range")
	}
	e.machine.r.PC = TEXT_ADDRESS + int(main)
	// stack pointer
	e.machine.r.write(29, STACK_ADDRESS)
	return nil
}

// resolve transfer 4 bytes into a execInst structure
func resolve(s []byte) (*execInst, error) {
	if len(s) != 4 {
		return nil, errors.New("resolve: bad machine code")
	}
	raw := binary.LittleEndian.Uint32(s)
	// obtain instruction name
	opcode := (raw >> 26) & 0x3F
	var name string
	var ok bool
	switch opcode {
	case 0:
		funct := raw & 0x3F
		name, ok = rInstructions[int(funct)]
		if !ok {
			return nil, errors.New("unsupported function code")
		}
	default:
		name, ok = ijInstructions[int(opcode)]
		if !ok {
			return nil, errors.New("unsupported opcode")
		}
	}
	isBranch := false
	switch name {
	case "beq", "bne", "j", "jr", "jal":
		isBranch = true
	}

	// resolve argument
	inst := instructionTable[name]
	var args []int
	for i, j := 0, 0; i < len(inst.syntax) && j < len(inst.formats); {
		switch {
		case inst.syntax[i]&argReg != 0:
			var arg int
			switch inst.formats[j] {
			case fmtRegS:
				arg = int((raw >> 21) & 0x1F)
			case fmtRegT:
				arg = int((raw >> 16) & 0x1F)
			case fmtRegD:
				arg = int((raw >> 11) & 0x1F)
			}
			args = append(args, arg)
			j++
		case inst.syntax[i]&argInteger != 0:
			var arg int
			switch inst.formats[j] {
			case fmtShamt:
				arg = int((raw >> 6) & 0x1F)
			case fmtImmediate:
				arg = int(raw & 0xFFFF)
			case fmtAddress:
				arg = int(raw & 0x3FFFFFF)
			}
			args = append(args, arg)
			j++
		case inst.syntax[i]&argAddr != 0:
			var arg int
			switch inst.formats[j] {
			case fmtRegS:
				arg = int((raw >> 21) & 0x1F)
			case fmtRegT:
				arg = int((raw >> 16) & 0x1F)
			case fmtRegD:
				arg = int((raw >> 11) & 0x1F)
			}
			args = append(args, arg)
			j++
			if inst.formats[j] != fmtImmediate {
				panic("something wrong...")
			}
			args = append(args, int(raw&0xFFFF))
			j++
		}
		i++
	}
	return &execInst{
		name:   name,
		f:      funcTable[name],
		args:   args,
		branch: isBranch,
	}, nil
}

func (e *Emulator) fetch() {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		e.exit <- EXIT_ERROR
	// 		e.err <- fmt.Errorf("%v", err)
	// 	}
	// }()
	i, err := e.machine.m.readWord(e.machine.r.PC)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, uint32(i))
	if err != nil {
		panic(err)
	}
	inst, err := resolve(buf.Bytes())
	if err != nil {
		panic(err)
	}
	e.inst <- *inst
}
