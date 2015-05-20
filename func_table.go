package mips

import (
	"bytes"
	"fmt"
)

type instFunc func(*Machine, ...int)

var (
	funcTable = map[string]instFunc{
		"add": func(m *Machine, args ...int) {
			m.r.write(args[0], m.r.read(args[1])+m.r.read(args[2]))
		},
		"addu": func(m *Machine, args ...int) {
			m.r.write(args[0], m.r.read(args[1])+m.r.read(args[2]))
		},
		"sub": func(m *Machine, args ...int) {
			m.r.write(args[0], m.r.read(args[1])-m.r.read(args[2]))
		},
		"subu": func(m *Machine, args ...int) {
			m.r.write(args[0], m.r.read(args[1])-m.r.read(args[2]))
		},
		"addi": func(m *Machine, args ...int) {
			m.r.write(args[0], m.r.read(args[1])+args[2])
		},
		"addiu": func(m *Machine, args ...int) {
			m.r.write(args[0], m.r.read(args[1])+int(uint(args[2])))
		},
		"mult": func(m *Machine, args ...int) {
			mult := m.r.read(args[0]) * m.r.read(args[1])
			m.r.HI = (mult >> 16) & 0xFFFF
			m.r.LO = mult & 0xFFFF
		},
		"multu": func(m *Machine, args ...int) {
			mult := uint(m.r.read(args[0])) * uint(m.r.read(args[1]))
			m.r.HI = int((mult >> 16) & 0xFFFF)
			m.r.LO = int(mult & 0xFFFF)
		},
		"div": func(m *Machine, args ...int) {
			m.r.HI = m.r.read(args[0]) % m.r.read(args[1])
			m.r.LO = m.r.read(args[0]) / m.r.read(args[1])
		},
		"divu": func(m *Machine, args ...int) {
			m.r.HI = int(uint(m.r.read(args[0])) % uint(m.r.read(args[1])))
			m.r.LO = int(uint(m.r.read(args[0])) / uint(m.r.read(args[1])))
		},
		"lw": func(m *Machine, args ...int) {
			i, err := m.m.readWord(m.r.read(args[1]) + args[2])
			checkInstErr(err)
			m.r.write(args[0], i)
		},
		"lh": func(m *Machine, args ...int) {
			i, err := m.m.readHalf(m.r.read(args[1]) + args[2])
			checkInstErr(err)
			m.r.write(args[0], i)
		},
		"lhu": func(m *Machine, args ...int) {
			i, err := m.m.readHalf(m.r.read(args[1]) + args[2])
			checkInstErr(err)
			m.r.write(args[0], int(uint(i)))
		},
		"lb": func(m *Machine, args ...int) {
			i, err := m.m.read(m.r.read(args[1]) + args[2])
			checkInstErr(err)
			m.r.write(args[0], int(i))
		},
		"lbu": func(m *Machine, args ...int) {
			i, err := m.m.read(m.r.read(args[1]) + args[2])
			checkInstErr(err)
			m.r.write(args[0], int(uint(i)))
		},
		"sw": func(m *Machine, args ...int) {
			err := m.m.writeWord(m.r.read(args[1])+args[2], m.r.read(args[0]))
			checkInstErr(err)
		},
		"sh": func(m *Machine, args ...int) {
			err := m.m.writeHalf(m.r.read(args[1])+args[2], m.r.read(args[0])&0xFFFF)
			checkInstErr(err)
		},
		"sb": func(m *Machine, args ...int) {
			err := m.m.write(m.r.read(args[1])+args[2], byte(m.r.read(args[0])&0xFF))
			checkInstErr(err)
		},
		"lui": func(m *Machine, args ...int) {
			m.r.write(args[0], args[1]<<16)
		},
		"mfhi": func(m *Machine, args ...int) {
			m.r.write(args[0], m.r.HI)
		},
		"mflo": func(m *Machine, args ...int) {
			m.r.write(args[0], m.r.LO)
		},
		"and": func(m *Machine, args ...int) {
			m.r.write(args[0], m.r.read(args[1])&m.r.read(args[2]))
		},
		"andi": func(m *Machine, args ...int) {
			m.r.write(args[0], m.r.read(args[1])&(args[2]&0x0000FFFF))
		},
		"or": func(m *Machine, args ...int) {
			m.r.write(args[0], m.r.read(args[1])|m.r.read(args[2]))
		},
		"ori": func(m *Machine, args ...int) {
			m.r.write(args[0], m.r.read(args[1])|(args[2]&0x0000FFFF))
		},
		"xor": func(m *Machine, args ...int) {
			m.r.write(args[0], m.r.read(args[1])^m.r.read(args[2]))
		},
		"nor": func(m *Machine, args ...int) {
			m.r.write(args[0], ^(m.r.read(args[1]) & m.r.read(args[2])))
		},
		"slt": func(m *Machine, args ...int) {
			if m.r.read(args[1]) < m.r.read(args[2]) {
				m.r.write(args[0], 1)
			} else {
				m.r.write(args[0], 0)
			}
		},
		"slti": func(m *Machine, args ...int) {
			if m.r.read(args[1]) < args[2] {
				m.r.write(args[0], 1)
			} else {
				m.r.write(args[0], 0)
			}
		},
		"sll": func(m *Machine, args ...int) {
			m.r.write(args[0], m.r.read(args[1])<<uint(args[2]))
		},
		"srl": func(m *Machine, args ...int) {
			m.r.write(args[0], int(uint(m.r.read(args[1]))>>uint(args[2])))
		},
		"sra": func(m *Machine, args ...int) {
			m.r.write(args[0], m.r.read(args[1])>>uint(args[2]))
		},
		"sllv": func(m *Machine, args ...int) {
			m.r.write(args[0], m.r.read(args[1])<<uint(m.r.read(args[2])))
		},
		"srlv": func(m *Machine, args ...int) {
			m.r.write(args[0], int(uint(m.r.read(args[1]))>>uint(m.r.read(args[2]))))
		},
		"srav": func(m *Machine, args ...int) {
			m.r.write(args[0], m.r.read(args[1])>>uint(m.r.read(args[2])))
		},
		"beq": func(m *Machine, args ...int) {
			if m.r.read(args[0]) == m.r.read(args[1]) {
				m.r.PC = m.r.PC + 4 + args[2]<<2
			} else {
				m.r.PC = m.r.PC + 4
			}
		},
		"bne": func(m *Machine, args ...int) {
			if m.r.read(args[0]) != m.r.read(args[1]) {
				m.r.PC = m.r.PC + 4 + args[2]<<2
			} else {
				m.r.PC = m.r.PC + 4
			}
		},
		"j": func(m *Machine, args ...int) {
			m.r.PC = (m.r.PC & 0xF00000000) | ((args[0] << 2) & 0x0FFFFFFF)
		},
		"jr": func(m *Machine, args ...int) {
			m.r.PC = m.r.read(args[0])
		},
		"jal": func(m *Machine, args ...int) {
			m.r.write(31, m.r.PC+4)
			m.r.PC = (m.r.PC & 0xF00000000) | ((args[0] << 2) & 0x0FFFFFFF)
		},
		"syscall": systemCall,
	}
)

func systemCall(m *Machine, args ...int) {
	a0 := registerTable["a0"]
	a1 := registerTable["a1"]
	// a2 := registerTable["a2"]
	v0 := registerTable["v0"]
	switch m.r.read(v0) {
	case 1: // print integer
		fmt.Printf("%d", m.r.read(a0))
	case 4: // print null-terminate string
		buf := new(bytes.Buffer)
		addr := m.r.read(a0)
		b, err := m.m.read(addr)
		for ; err == nil && b != 0; b, err = m.m.read(addr) {
			err := buf.WriteByte(b)
			checkInstErr(err)
			addr++
		}
		checkInstErr(err)
		fmt.Printf("%s", buf.String())
	case 5: // read integer
		var i int
		_, err := fmt.Scanf("%d", &i)
		checkInstErr(err)
		m.r.write(v0, i)
	case 8:
		var s string
		_, err := fmt.Scanf("%s", &s)
		checkInstErr(err)
		addr := m.r.read(a0)
		max := m.r.read(a1)
		if max < len(s) {
			s = s[:max]
		}
		err = m.m.writeBytes(addr, []byte(s))
		checkInstErr(err)
	case 10:
		m.exit = true
	default:
	}
}

func checkInstErr(err error) {
	if err != nil {
		panic(err)
	}
}
