package mips

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestEmulator(t *testing.T) {
	var (
		input = `.text
	main:
	li $v0, 4
	la $a0, hello
	syscall
	li $v0, 10
	syscall
.data
	hello: .ascii "Hello, world!"
	.byte 0x0A, 0`
	)
	fmt.Println(input)
	asm := NewAssembler(strings.NewReader(input))
	raw, err := asm.Assemble()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%q\n", string(raw))

	disasm := NewDisassembler(bytes.NewBuffer(raw))
	code, err := disasm.Disassemble()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", string(code))

	em := NewEmulator()
	err = em.LoadAndRun(raw)
	if err != nil {
		log.Fatal(err)
	}
	err = em.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
