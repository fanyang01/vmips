package asm

import "fmt"

func (p *parser) pseudoFilter(items <-chan parseItem) <-chan parseItem {
	result := make(chan parseItem)
	go func() {
		for item := range items {
			switch item.typ {
			case itemInst:
				if instructionTable[item.instruction].typ == "P" {
					p.translate(item, result)
				} else {
					result <- item
				}
			default:
				result <- item
			}
		}
		close(result)
	}()
	return result
}

func (p *parser) translate(i parseItem, result chan<- parseItem) {
	switch i.instruction {
	case "move":
		result <- parseItem{
			typ:         itemInst,
			instruction: "add",
			registers:   []string{i.registers[0], i.registers[1], "$zero"},
		}
	case "not":
		result <- parseItem{
			typ:         itemInst,
			instruction: "nor",
			registers:   []string{i.registers[0], i.registers[1], "$zero"},
		}
	case "clear":
		result <- parseItem{
			typ:         itemInst,
			instruction: "add",
			registers:   []string{i.registers[0], "$zero", "$zero"},
		}
	case "li":
		result <- parseItem{
			typ:         itemInst,
			instruction: "lui",
			registers:   []string{i.registers[0]},
			imme:        (i.imme >> 16) & 0x0000FFFF,
		}
		result <- parseItem{
			typ:         itemInst,
			instruction: "ori",
			registers:   []string{i.registers[0], i.registers[0]},
			imme:        i.imme & 0x0000FFFF,
		}
	case "la":
		addr := p.labels[i.label].address
		result <- parseItem{
			typ:         itemInst,
			instruction: "lui",
			registers:   []string{i.registers[0]},
			imme:        (addr >> 16) & 0x0000FFFF,
		}
		result <- parseItem{
			typ:         itemInst,
			instruction: "ori",
			registers:   []string{i.registers[0], i.registers[0]},
			imme:        addr & 0x0000FFFF,
		}
	case "bgt":
		result <- parseItem{
			typ:         itemInst,
			instruction: "slt",
			registers:   []string{"$at", i.registers[1], i.registers[0]},
		}
		result <- parseItem{
			typ:         itemInst,
			instruction: "bne",
			registers:   []string{"$at", "$zero"},
			imme:        i.imme,
		}
	case "blt":
		result <- parseItem{
			typ:         itemInst,
			instruction: "slt",
			registers:   []string{"$at", i.registers[0], i.registers[1]},
		}
		result <- parseItem{
			typ:         itemInst,
			instruction: "bne",
			registers:   []string{"$at", "$zero"},
			imme:        i.imme,
		}
	case "bge":
		result <- parseItem{
			typ:         itemInst,
			instruction: "slt",
			registers:   []string{"$at", i.registers[0], i.registers[1]},
		}
		result <- parseItem{
			typ:         itemInst,
			instruction: "beq",
			registers:   []string{"$at", "$zero"},
			imme:        i.imme,
		}
	case "ble":
		result <- parseItem{
			typ:         itemInst,
			instruction: "slt",
			registers:   []string{"$at", i.registers[1], i.registers[0]},
		}
		result <- parseItem{
			typ:         itemInst,
			instruction: "beq",
			registers:   []string{"$at", "$zero"},
			imme:        i.imme,
		}
	case "bgtu":
		result <- parseItem{
			typ:         itemInst,
			instruction: "sltu",
			registers:   []string{"$at", i.registers[1], i.registers[0]},
		}
		result <- parseItem{
			typ:         itemInst,
			instruction: "bne",
			registers:   []string{"$at", "$zero"},
			imme:        i.imme,
		}
	case "bgtz":
		result <- parseItem{
			typ:         itemInst,
			instruction: "slt",
			registers:   []string{"$at", "$zero", i.registers[0]},
		}
		result <- parseItem{
			typ:         itemInst,
			instruction: "bne",
			registers:   []string{"$at", "$zero"},
			imme:        i.imme,
		}
	case "beqz":
		result <- parseItem{
			typ:         itemInst,
			instruction: "beq",
			registers:   []string{"$zero", i.registers[0]},
			imme:        i.imme,
		}
	case "mul":
		result <- parseItem{
			typ:         itemInst,
			instruction: "mult",
			registers:   []string{i.registers[1], i.registers[2]},
		}
		result <- parseItem{
			typ:         itemInst,
			instruction: "mflo",
			registers:   []string{i.registers[0]},
			imme:        i.imme,
		}
	case "divq":
		result <- parseItem{
			typ:         itemInst,
			instruction: "div",
			registers:   []string{i.registers[1], i.registers[2]},
		}
		result <- parseItem{
			typ:         itemInst,
			instruction: "mflo",
			registers:   []string{i.registers[0]},
			imme:        i.imme,
		}
	case "rem":
		result <- parseItem{
			typ:         itemInst,
			instruction: "div",
			registers:   []string{i.registers[1], i.registers[2]},
		}
		result <- parseItem{
			typ:         itemInst,
			instruction: "mfhi",
			registers:   []string{i.registers[0]},
			imme:        i.imme,
		}
	default:
		result <- parseItem{
			typ: itemError,
			err: fmt.Sprintf("invalid pseudo instruction %q",
				i.instruction),
		}
	}
}
