package mips

/*
instruction format:

	R   opcode[6] rs[5] rt[5] rd[5] shamt[5] funct[6]

	I   opcode[6] rs[5] rt[5] immediate[16]

	J   opcode[6] address[26]
*/

type argType int
type fmtType int

type instInfo struct {
	typ     string
	syntax  []argType
	formats []fmtType
	opcode  int
	funct   int
	size    int
}

const (
	// argument type
	argReg     argType = 1 << iota // register
	argInteger                     // immediate integer
	argLabel                       // label
	argAddr                        // format of address is C($s)
	// format type
	fmtRegD fmtType = 1 << iota
	fmtRegS
	fmtRegT
	fmtShamt
	fmtImmediate
	fmtAddress
)

var (
	instructionTable = map[string]instInfo{
		"add": instInfo{
			typ:     "R",
			syntax:  []argType{argReg, argReg, argReg},
			formats: []fmtType{fmtRegD, fmtRegS, fmtRegT},
			opcode:  0,
			funct:   0x20,
		},
		"addu": instInfo{
			typ:     "R",
			syntax:  []argType{argReg, argReg, argReg},
			formats: []fmtType{fmtRegD, fmtRegS, fmtRegT},
			opcode:  0,
			funct:   0x21,
		},
		"sub": instInfo{
			typ:     "R",
			syntax:  []argType{argReg, argReg, argReg},
			formats: []fmtType{fmtRegD, fmtRegS, fmtRegT},
			opcode:  0,
			funct:   0x22,
		},
		"subu": instInfo{
			typ:     "R",
			syntax:  []argType{argReg, argReg, argReg},
			formats: []fmtType{fmtRegD, fmtRegS, fmtRegT},
			opcode:  0,
			funct:   0x23,
		},
		"and": instInfo{
			typ:     "R",
			syntax:  []argType{argReg, argReg, argReg},
			formats: []fmtType{fmtRegD, fmtRegS, fmtRegT},
			opcode:  0,
			funct:   0x24,
		},
		"or": instInfo{
			typ:     "R",
			syntax:  []argType{argReg, argReg, argReg},
			formats: []fmtType{fmtRegD, fmtRegS, fmtRegT},
			opcode:  0,
			funct:   0x25,
		},
		"xor": instInfo{
			typ:     "R",
			syntax:  []argType{argReg, argReg, argReg},
			formats: []fmtType{fmtRegD, fmtRegS, fmtRegT},
			opcode:  0,
			funct:   0x26,
		},
		"nor": instInfo{
			typ:     "R",
			syntax:  []argType{argReg, argReg, argReg},
			formats: []fmtType{fmtRegD, fmtRegS, fmtRegT},
			opcode:  0,
			funct:   0x27,
		},
		"slt": instInfo{
			typ:     "R",
			syntax:  []argType{argReg, argReg, argReg},
			formats: []fmtType{fmtRegD, fmtRegS, fmtRegT},
			opcode:  0,
			funct:   0x2A,
		},
		"sllv": instInfo{
			typ:     "R",
			syntax:  []argType{argReg, argReg, argReg},
			formats: []fmtType{fmtRegD, fmtRegS, fmtRegT},
			opcode:  0,
			funct:   0x4,
		},
		"srlv": instInfo{
			typ:     "R",
			syntax:  []argType{argReg, argReg, argReg},
			formats: []fmtType{fmtRegD, fmtRegS, fmtRegT},
			opcode:  0,
			funct:   0x6,
		},
		"srav": instInfo{
			typ:     "R",
			syntax:  []argType{argReg, argReg, argReg},
			formats: []fmtType{fmtRegD, fmtRegS, fmtRegT},
			opcode:  0,
			funct:   0x7,
		},
		// R2
		"sll": instInfo{
			typ:     "R",
			syntax:  []argType{argReg, argReg, argInteger},
			formats: []fmtType{fmtRegD, fmtRegT, fmtShamt},
			opcode:  0,
			funct:   0x0,
		},
		"srl": instInfo{
			typ:     "R",
			syntax:  []argType{argReg, argReg, argInteger},
			formats: []fmtType{fmtRegD, fmtRegT, fmtShamt},
			opcode:  0,
			funct:   0x2,
		},
		"sra": instInfo{
			typ:     "R",
			syntax:  []argType{argReg, argReg, argInteger},
			formats: []fmtType{fmtRegD, fmtRegT, fmtShamt},
			opcode:  0,
			funct:   0x3,
		},
		// R3
		"mult": instInfo{
			typ:     "R",
			syntax:  []argType{argReg, argReg},
			formats: []fmtType{fmtRegS, fmtRegT},
			opcode:  0,
			funct:   0x18,
		},
		"multu": instInfo{
			typ:     "R",
			syntax:  []argType{argReg, argReg},
			formats: []fmtType{fmtRegS, fmtRegT},
			opcode:  0,
			funct:   0x19,
		},
		"div": instInfo{
			typ:     "R",
			syntax:  []argType{argReg, argReg},
			formats: []fmtType{fmtRegS, fmtRegT},
			opcode:  0,
			funct:   0x1A,
		},
		"divu": instInfo{
			typ:     "R",
			syntax:  []argType{argReg, argReg},
			formats: []fmtType{fmtRegS, fmtRegT},
			opcode:  0,
			funct:   0x1B,
		},
		// R4
		"mfhi": instInfo{
			typ:     "R",
			syntax:  []argType{argReg},
			formats: []fmtType{fmtRegD},
			opcode:  0,
			funct:   0x10,
		},
		"mflo": instInfo{
			typ:     "R",
			syntax:  []argType{argReg},
			formats: []fmtType{fmtRegD},
			opcode:  0,
			funct:   0x12,
		},
		// R5
		"jr": instInfo{
			typ:     "R",
			syntax:  []argType{argReg},
			formats: []fmtType{fmtRegS},
			opcode:  0,
			funct:   0x8,
		},
		// R6
		"syscall": instInfo{
			typ:     "R",
			syntax:  []argType{},
			formats: []fmtType{},
			opcode:  0,
			funct:   0xC,
		},
		// I1
		"addi": instInfo{
			typ:     "I",
			syntax:  []argType{argReg, argReg, argInteger},
			formats: []fmtType{fmtRegT, fmtRegS, fmtImmediate},
			opcode:  0x8,
		},
		"addiu": instInfo{
			typ:     "I",
			syntax:  []argType{argReg, argReg, argInteger},
			formats: []fmtType{fmtRegT, fmtRegS, fmtImmediate},
			opcode:  0x9,
		},
		"andi": instInfo{
			typ:     "I",
			syntax:  []argType{argReg, argReg, argInteger},
			formats: []fmtType{fmtRegT, fmtRegS, fmtImmediate},
			opcode:  0xC,
		},
		"ori": instInfo{
			typ:     "I",
			syntax:  []argType{argReg, argReg, argInteger},
			formats: []fmtType{fmtRegT, fmtRegS, fmtImmediate},
			opcode:  0xD,
		},
		"slti": instInfo{
			typ:     "I",
			syntax:  []argType{argReg, argReg, argInteger},
			formats: []fmtType{fmtRegT, fmtRegS, fmtImmediate},
			opcode:  0xA,
		},
		// I2
		"bne": instInfo{
			typ:     "I",
			syntax:  []argType{argReg, argReg, argInteger | argLabel},
			formats: []fmtType{fmtRegS, fmtRegT, fmtImmediate},
			opcode:  0x4,
		},
		"beq": instInfo{
			typ:     "I",
			syntax:  []argType{argReg, argReg, argInteger | argLabel},
			formats: []fmtType{fmtRegS, fmtRegT, fmtImmediate},
			opcode:  0x5,
		},
		// I3
		"lw": instInfo{
			typ:     "I",
			syntax:  []argType{argReg, argAddr},
			formats: []fmtType{fmtRegT, fmtRegS, fmtImmediate},
			opcode:  0x23,
		},
		"lh": instInfo{
			typ:     "I",
			syntax:  []argType{argReg, argAddr},
			formats: []fmtType{fmtRegT, fmtRegS, fmtImmediate},
			opcode:  0x21,
		},
		"lhu": instInfo{
			typ:     "I",
			syntax:  []argType{argReg, argAddr},
			formats: []fmtType{fmtRegT, fmtRegS, fmtImmediate},
			opcode:  0x25,
		},
		"lb": instInfo{
			typ:     "I",
			syntax:  []argType{argReg, argAddr},
			formats: []fmtType{fmtRegT, fmtRegS, fmtImmediate},
			opcode:  0x20,
		},
		"lbu": instInfo{
			typ:     "I",
			syntax:  []argType{argReg, argAddr},
			formats: []fmtType{fmtRegT, fmtRegS, fmtImmediate},
			opcode:  0x24,
		},
		"sw": instInfo{
			typ:     "I",
			syntax:  []argType{argReg, argAddr},
			formats: []fmtType{fmtRegT, fmtRegS, fmtImmediate},
			opcode:  0x2B,
		},
		"sh": instInfo{
			typ:     "I",
			syntax:  []argType{argReg, argAddr},
			formats: []fmtType{fmtRegT, fmtRegS, fmtImmediate},
			opcode:  0x29,
		},
		"sb": instInfo{
			typ:     "I",
			syntax:  []argType{argReg, argAddr},
			formats: []fmtType{fmtRegT, fmtRegS, fmtImmediate},
			opcode:  0x28,
		},
		// I4
		"lui": instInfo{
			typ:     "I",
			syntax:  []argType{argReg, argInteger},
			formats: []fmtType{fmtRegT, fmtImmediate},
			opcode:  0xF,
		},
		// J
		"j": instInfo{
			typ:     "J",
			syntax:  []argType{argInteger | argLabel},
			formats: []fmtType{fmtAddress},
			opcode:  0x2,
		},
		"jal": instInfo{
			typ:     "J",
			syntax:  []argType{argInteger | argLabel},
			formats: []fmtType{fmtAddress},
			opcode:  0x3,
		},
		// Pseudo
		"mul": instInfo{
			typ:    "P",
			syntax: []argType{argReg, argReg, argReg},
			size:   2,
		},
		"divq": instInfo{
			typ:    "P",
			syntax: []argType{argReg, argReg, argReg},
			size:   2,
		},
		"rem": instInfo{
			typ:    "P",
			syntax: []argType{argReg, argReg, argReg},
			size:   2,
		},
		"bgt": instInfo{
			typ:    "P",
			syntax: []argType{argReg, argReg, argInteger | argLabel},
			size:   2,
		},
		"blt": instInfo{
			typ:    "P",
			syntax: []argType{argReg, argReg, argInteger | argLabel},
			size:   2,
		},
		"bge": instInfo{
			typ:    "P",
			syntax: []argType{argReg, argReg, argInteger | argLabel},
			size:   2,
		},
		"ble": instInfo{
			typ:    "P",
			syntax: []argType{argReg, argReg, argInteger | argLabel},
			size:   2,
		},
		"bgtu": instInfo{
			typ:    "P",
			syntax: []argType{argReg, argReg, argInteger | argLabel},
			size:   2,
		},
		"bgtz": instInfo{
			typ:    "P",
			syntax: []argType{argReg, argInteger | argLabel},
			size:   2,
		},
		"beqz": instInfo{
			typ:    "P",
			syntax: []argType{argReg, argInteger | argLabel},
			size:   1,
		},
		"move": instInfo{
			typ:    "P",
			syntax: []argType{argReg, argReg},
			size:   1,
		},
		"not": instInfo{
			typ:    "P",
			syntax: []argType{argReg, argReg},
			size:   1,
		},
		"li": instInfo{
			typ:    "P",
			syntax: []argType{argReg, argInteger},
			size:   2,
		},
		"la": instInfo{
			typ:    "P",
			syntax: []argType{argReg, argLabel},
			size:   2,
		},
		"clear": instInfo{
			typ:    "P",
			syntax: []argType{argReg},
			size:   1,
		},
		"nop": instInfo{
			typ:    "P",
			syntax: []argType{},
			size:   1,
		},
	}
	registerNames = []string{
		"zero", "at",
		"v0", "v1",
		"a0", "a1", "a2", "a3",
		"t0", "t1", "t2", "t3", "t4", "t5", "t6", "t7",
		"s0", "s1", "s2", "s3", "s4", "s5", "s6", "s7",
		"t8", "t9",
		"k0", "k1",
		"gp", "sp", "fp", "ra",
	}
	registerTable = make(map[string]int)
)

func init() {
	for i, r := range registerNames {
		registerTable[r] = i
	}
}
