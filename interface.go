package mips

import "bufio"

type Assembler struct {
	r *bufio.Reader
}

func NewAssembler(r *bufio.Reader) *Assembler {
	return &Assembler{
		r: r,
	}
}

func (a *Assembler) Assemble() ([]byte, error) {
	ch := parse(a.r)
	return assemble(ch)
}
