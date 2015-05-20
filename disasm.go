package mips

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var (
	headerPattern = regexp.MustCompile(`^text:([0-9]+),data:([0-9]+),main:([0-9]+)$`)
)

type Disassembler struct {
	r          *bufio.Reader
	textOffset int
	dataOffset int
	mainOffset int
	eof        bool
}

// Disassemble disassemble each 4 bytes into an instruction
func Disassemble(raw []byte) ([]byte, error) {
	var result []byte
	for i := 0; i+4 <= len(raw); i += 4 {
		s, err := disasm(raw[i : i+4])
		if err != nil {
			return nil, err
		}
		result = append(result, s...)
		result = append(result, '\n')
	}
	return bytes.TrimSpace(result), nil
}

// NewDisasssembler creates a disassembler to disassemble
// object files assembled by this package
func NewDisassembler(r io.Reader) *Disassembler {
	return &Disassembler{
		r: bufio.NewReader(r),
	}
}

// Disassemble starts the disassembler
func (d *Disassembler) Disassemble() ([]byte, error) {
	err := d.parseHeader()
	if err != nil {
		return nil, err
	}
	return d.disassemble()
}

func (d *Disassembler) parseHeader() error {
	line, err := d.r.ReadString('\n')
	if err != nil {
		return err
	}
	line = line[:len(line)-1]
	if !headerPattern.MatchString(line) {
		return errors.New("invalid header")
	}

	sub := headerPattern.FindStringSubmatch(line)

	text, err := strconv.ParseInt(sub[1], 10, 32)
	if err != nil {
		panic(err.Error())
	}
	d.textOffset = int(text)
	data, err := strconv.ParseInt(sub[2], 10, 32)
	if err != nil {
		panic(err.Error())
	}
	d.dataOffset = int(data)
	main, err := strconv.ParseInt(sub[3], 10, 32)
	if err != nil {
		panic(err.Error())
	}
	d.mainOffset = int(main)
	return nil
}

func (d *Disassembler) disassemble() ([]byte, error) {
	var i int
	for ; i < d.textOffset; i++ {
		_, err := d.r.ReadByte()
		if err != nil {
			return nil, err
		}
	}

	cmp := false
	if d.dataOffset > d.textOffset {
		cmp = true
	}

	ret := []byte{}
	for ; ; i += 4 {
		if cmp && i >= d.dataOffset {
			break
		}

		var s []byte
		for j := 0; j < 4; j++ {
			b, err := d.r.ReadByte()
			if err != nil {
				if err == io.EOF {
					d.eof = true
					goto RETURN
				}
				return nil, err
			}
			s = append(s, b)
		}

		line, err := disasm(s)
		if err != nil {
			return nil, err
		}
		ret = append(ret, line...)
		ret = append(ret, '\n')
	}
RETURN:
	ret = bytes.TrimSpace(ret)
	return ret, nil
}

func disasm(s []byte) ([]byte, error) {
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

	// disassemble
	inst := instructionTable[name]
	var args []string
	for i, j := 0, 0; i < len(inst.syntax) && j < len(inst.formats); {
		var token string
		switch {
		case inst.syntax[i]&argReg != 0:
			switch inst.formats[j] {
			case fmtRegS:
				token = registerNames[(raw>>21)&0x1F]
			case fmtRegT:
				token = registerNames[(raw>>16)&0x1F]
			case fmtRegD:
				token = registerNames[(raw>>11)&0x1F]
			}
			token = "$" + token
			j++
		case inst.syntax[i]&argInteger != 0:
			switch inst.formats[j] {
			case fmtShamt:
				token = fmt.Sprintf("%d", (raw>>6)&0x1F)
			case fmtImmediate:
				token = fmt.Sprintf("%d", raw&0xFFFF)
			case fmtAddress:
				token = fmt.Sprintf("%#v", raw&0x3FFFFFF)
			}
			j++
		case inst.syntax[i]&argAddr != 0:
			var reg string
			switch inst.formats[j] {
			case fmtRegS:
				reg = registerNames[(raw>>21)&0x1F]
			case fmtRegT:
				reg = registerNames[(raw>>16)&0x1F]
			case fmtRegD:
				reg = registerNames[(raw>>11)&0x1F]
			}
			j++
			if inst.formats[j] != fmtImmediate {
				panic("something wrong...")
			}
			token = fmt.Sprintf("%d($%s)", raw&0xFFFF, reg)
			j++
		}
		args = append(args, token)
		i++
	}
	return []byte(fmt.Sprintf("%s\t%s", name, strings.Join(args, ", "))), nil
}
