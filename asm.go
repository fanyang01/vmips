package mips

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

func assemble(items <-chan parseItem) (b []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("runtime panic: %v", r)
		}
	}()
	var section **bytes.Buffer
	textSection := new(bytes.Buffer)
	dataSection := new(bytes.Buffer)
	section = &textSection
LOOP:
	for item := range items {
		switch item.typ {
		case itemError:
			return nil, errors.New(item.err)
		case itemEOF:
			break LOOP
		case itemInst:
			_, err := (*section).Write(asmInst(item))
			if err != nil {
				return nil, err
			}
		case itemDir:
			switch item.directive {
			case "text":
				section = &textSection
			case "data":
				section = &dataSection
			case "globl":
				// TODO
			default:
				_, err := (*section).Write(asmDir(item))
				if err != nil {
					return nil, err
				}
			}
		}
	}
	b = textSection.Bytes()
	b = append(b, dataSection.Bytes()...)
	return
}

func asmInst(item parseItem) []byte {
	var raw int
	inst := instructionTable[item.instruction]
	raw |= inst.opcode << 26
	if inst.typ == "R" {
		raw |= inst.funct
	}
	for i, f := range inst.formats {
		switch f {
		case fmtRegS:
			raw |= registerTable[item.registers[i][1:]] << 21
		case fmtRegT:
			raw |= registerTable[item.registers[i][1:]] << 16
		case fmtRegD:
			raw |= registerTable[item.registers[i][1:]] << 11
		case fmtShamt:
			if item.imme < 0 || item.imme > 31 {
				panic(fmt.Sprintf("shift amount %d out of range", item.imme))
			}
			raw |= item.imme << 6
		case fmtImmediate:
			if item.imme >= 1<<16 || item.imme <= -(1<<16) {
				panic(fmt.Sprintf("immediate number %d out of range",
					item.imme))
			}
			raw |= item.imme
		case fmtAddress:
			if item.imme >= 1<<26 || item.imme < 0 {
				panic(fmt.Sprintf("immediate number %d out of range", item.imme))
			}
			raw |= item.imme
		default:
			// shouldn't get here
		}
	}
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, uint32(raw))
	if err != nil {
		panic(err.Error())
	}
	return buf.Bytes()
}

func asmDir(item parseItem) []byte {
	switch item.directive {
	case "space":
		return make([]byte, item.data.(int))
	case "align":
		width := uint(item.data.(int))
		if rem := item.address % (1 << width); rem != 0 {
			return make([]byte, 1<<width-rem)
		}
		return []byte{}
	case "byte":
		var s []byte
		for _, n := range item.data.([]int) {
			s = append(s, byte(n&0xFF))
		}
		return s
	case "half":
		var s []byte
		for _, n := range item.data.([]int) {
			s = append(s, byte(n&0xFF), byte((n>>8)&0xFF))
		}
		return s
	case "word":
		var s []byte
		for _, n := range item.data.([]int) {
			s = append(s, byte(n&0xFF), byte((n>>8)&0xFF),
				byte(n>>16&0xFF), byte((n>>24)&0xFF))
		}
		return s
	case "ascii":
		return []byte(item.data.(string))
	case "asciiz":
		return append([]byte(item.data.(string)), byte(0))
	}
	return []byte{}
}
