package mips

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestAssemble(t *testing.T) {
	input := []string{
		// 000000 00100 01000 10000 00000 100000
		// 0x00888020
		"add $s0, $a0, $t0",
		// 100011 01000 10000 0000 0000 0000 0000
		// 0x8D100000
		"lw $s0, 0($t0)",
		`.globl main
		add $s0, $a0, $t0
		main:
		lw $s0, 0($t0)
		.data
		.ascii "hello, world"`,
	}
	input1 := []string{
		"add $s0, $a0, $t0",
		"lw $s0, 0($t0)",
	}
	expectedResult := [][]byte{
		[]byte("text:0,data:4,main:0\n\x20\x80\x88\x00"),
		[]byte("text:0,data:4,main:0\n\x00\x00\x10\x8d"),
		[]byte("text:0,data:8,main:4\n\x20\x80\x88\x00\x00\x00\x10\x8dhello, world"),
	}
	expectedResult1 := [][]byte{
		[]byte("\x20\x80\x88\x00"),
		[]byte("\x00\x00\x10\x8d"),
	}
	for i, in := range input {
		a := NewAssembler(strings.NewReader(in))
		b, err := a.Assemble()
		if err != nil {
			log.Println(err)
			t.Fail()
		}
		if !bytes.Equal(b, expectedResult[i]) {
			log.Printf("expect %q, got %q\n", string(expectedResult[i]), string(b))
			t.Fail()
		}
	}
	for i, in := range input1 {
		b, err := Assemble([]byte(in))
		if err != nil {
			log.Println(err)
			t.Fail()
		}
		if !bytes.Equal(b, expectedResult1[i]) {
			log.Printf("expect %q, got %q\n", string(expectedResult[i]), string(b))
			t.Fail()
		}
	}
}
