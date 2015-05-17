package mips

import (
	"bytes"
	"log"
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
		`add $s0, $a0, $t0
		lw $s0, 0($t0)`,
	}
	expectedResult := [][]byte{
		[]byte{0x20, 0x80, 0x88, 0x00},
		[]byte{0x00, 0x00, 0x10, 0x8D},
		[]byte{0x20, 0x80, 0x88, 0x00, 0x00, 0x00, 0x10, 0x8D},
	}
	for i, in := range input {
		b, err := Assemble(in)
		if err != nil {
			log.Println(err)
			t.Fail()
		}
		if !bytes.Equal(b, expectedResult[i]) {
			log.Printf("expect % x, got % x\n", expectedResult[i], b)
			t.Fail()
		}
	}
}
