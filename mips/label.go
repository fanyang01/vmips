package mips

import "fmt"

// labelFilter need to read all items in
func (p *parser) labelFilter(items <-chan parseItem) <-chan parseItem {
	// NOTE: Can't parallelize it!
	p.readAllLabel(items)
	return p.replaceLabel()
}

func (p *parser) readAllLabel(items <-chan parseItem) {
	textAddress := TEXT_ADDRESS
	dataAddress := DATA_ADDRESS
	addr := &textAddress
LOOP:
	for item := range items {
		item.address = *addr
		switch item.typ {
		case itemLabel:
			if l, ok := p.labels[item.label]; ok {
				i := parseItem{
					typ: itemError,
					err: fmt.Sprintf("line %d: label %q defined twice"+
						"(has defined at line %d)",
						item.line, item.label, l.line),
				}
				p.itemList.PushBack(i)
				return
			}
			p.labels[item.label] = item
		case itemDir:
			switch item.directive {
			case "text":
				addr = &textAddress
			case "data":
				addr = &dataAddress
			case "byte":
				*addr += len(item.data.([]int))
			case "half":
				*addr += len(item.data.([]int)) << 1
			case "word":
				*addr += len(item.data.([]int)) << 2
			case "space":
				*addr += item.data.(int)
			case "align":
				width := uint(item.data.(int))
				if rem := *addr % (1 << width); rem != 0 {
					*addr += 1<<width - rem
				}
			case "ascii":
				*addr += len(item.data.(string))
			case "asciiz":
				*addr += len(item.data.(string)) + 1
			case "globl":
				item.label = item.data.(string)
			}
		case itemInst:
			inst := instructionTable[item.instruction]
			if inst.typ == "P" {
				*addr += inst.size << 2
			} else {
				*addr += 4
			}
		case itemError:
			p.itemList.Init()
			p.itemList.PushBack(item)
			break LOOP
		}
		p.itemList.PushBack(item)
	}
}

func (p *parser) replaceLabel() <-chan parseItem {
	result := make(chan parseItem)
	go func() {
	LOOP:
		for e := p.itemList.Front(); e != nil; e = e.Next() {
			item := e.Value.(parseItem)
			switch item.typ {
			case itemLabel:
				continue
			case itemInst:
				if item.label != "" {
					if l, ok := p.labels[item.label]; ok {
						if instructionTable[item.instruction].typ == "J" {
							item.imme = l.address >> 2
						} else {
							item.imme = (l.address - (item.address + 4)) >> 2
						}
					} else {
						result <- parseItem{
							typ: itemError,
							err: fmt.Sprintf("label %q not defined",
								item.label),
						}
						break LOOP
					}
				}
				result <- item
			case itemDir:
				switch item.directive {
				case "globl":
					if l, ok := p.labels[item.label]; ok {
						item.address = l.address
					} else {
						result <- parseItem{
							typ: itemError,
							err: fmt.Sprintf("label %q not defined",
								item.label),
						}
						break LOOP
					}
				}
				result <- item
			case itemError:
				result <- item
				break LOOP
			default:
				result <- item
			}
		}
		close(result)
	}()
	return result
}
