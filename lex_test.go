package mips

import (
	"log"
	"testing"
)

const (
	input = `
Start01:
	add $t0, $s0, $s1
	add $t1, $a0, $zero
	j Next
	lw $s0, 0X20($s1)
	addi $t0, $s2, 0x10
	move $s0, $t0
	la $s0, Next
	jr $s0
	.data
		val: .byte 'c','b'
	.text
Next:
	add $t2, $a0, $zero`
)

func TestLex(t *testing.T) {
	ch := lex(input)
	for token := range ch {
		if token.typ == tokenError {
			log.Println(token)
			t.Fail()
		}
	}
}
