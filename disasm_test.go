package mips

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestDisassemble(t *testing.T) {
	input := [][]byte{
		[]byte("text:0,data:4,main:0\n\x20\x80\x88\x00"),
		[]byte("text:0,data:4,main:0\n\x00\x00\x10\x8d"),
	}
	expectedResult := [][]byte{
		[]byte("add\t$s0, $a0, $t0"),
		[]byte("lw\t$s0, 0($t0)"),
	}
	for i, in := range input {
		d := NewDisassembler(strings.NewReader(string(in)))
		b, err := d.Disassemble()
		if err != nil {
			log.Println(err)
			t.Fail()
		}
		if !bytes.Equal(b, expectedResult[i]) {
			log.Printf("expected %q, got %q\n", string(expectedResult[i]), string(b))
			t.Fail()
		}
	}
}
