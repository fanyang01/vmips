package mips

import (
	"container/list"
	"fmt"
	"strconv"
	"strings"
)

//go:generate stringer -type=itemType
type itemType int

const (
	itemEOF itemType = 1 << iota
	itemError
	itemInst  // instruction
	itemLabel // label
	itemDir   // directive
	// Address layout
	TEXT_ADDRESS = 0x8000000
	DATA_ADDRESS = 0x8100000
)

// parseItem can be a instruction or a label
type parseItem struct {
	typ         itemType
	instruction string
	registers   []string
	directive   string
	data        interface{} // args of directive
	imme        int         // immediate constant
	label       string
	address     int
	line        int
	err         string
}

type parser struct {
	items     chan parseItem
	tokens    <-chan token
	labels    map[string]parseItem
	itemList  *list.List
	entryAddr int
	line      int
}

type parseFn func(*parser) parseFn

func parse(input string) (*parser, <-chan parseItem) {
	p := &parser{
		items:    make(chan parseItem),
		tokens:   lex(input),
		itemList: list.New(),
		labels:   make(map[string]parseItem),
	}
	go p.run(parseStart)
	return p, p.pseudoFilter(p.labelFilter(p.items))
}

func (p *parser) run(initState parseFn) {
	for state := initState; state != nil; {
		state = state(p)
	}
	close(p.items)
}

func parseStart(p *parser) parseFn {
	t := <-p.tokens
	switch t.typ {
	case tokenInstruction:
		return p.parseInst(t.val)
	case tokenLabelDef:
		return p.parseLabel(t.val)
	case tokenEndline:
		// Skip blank lines
		p.line++
		return parseStart
	case tokenDirective:
		return p.parseDir(t.val)
	case tokenEOF:
		p.items <- parseItem{
			typ: itemEOF,
		}
		return nil
	default:
		return p.errorf("unexpected token %q(type %s)", t.val, t.typ)
	}
}

func (p *parser) errorf(format string, args ...interface{}) parseFn {
	p.items <- parseItem{
		typ:  itemError,
		err:  fmt.Sprintf("line %d: ", p.line+1) + fmt.Sprintf(format, args...),
		line: p.line,
	}
	return nil
}

func parseEndline(p *parser) parseFn {
	token := <-p.tokens
	switch token.typ {
	case tokenEndline:
		p.line++
		return parseStart
	case tokenEOF:
		p.items <- parseItem{
			typ: itemEOF,
		}
		return nil
	default:
		return p.errorf("unexpected token %q(type %q), expect Endline",
			token.val, token.typ)
	}
}

func (p *parser) parseLabel(label string) parseFn {
	p.items <- parseItem{
		typ:   itemLabel,
		label: label,
		line:  p.line,
	}
	return parseStart
}

func (p *parser) parseInst(inst string) parseFn {
	expectTokens := make(chan []tokenType, 16)
	ret := make(chan parseFn)
	go p.parseArgs(inst, expectTokens, ret)

	args := instructionTable[inst].syntax
	for i, s := range args {
		switch s {
		case argReg:
			expectTokens <- []tokenType{tokenRegister}
		case argInteger:
			expectTokens <- []tokenType{tokenInteger}
		case argLabel:
			expectTokens <- []tokenType{tokenLabel}
		case argInteger | argLabel:
			expectTokens <- []tokenType{tokenLabel, tokenInteger}
		case argAddr:
			expectTokens <- []tokenType{tokenInteger}
			expectTokens <- []tokenType{tokenLeftParenthese}
			expectTokens <- []tokenType{tokenRegister}
			expectTokens <- []tokenType{tokenRightParenthese}
		default:
			// shouldn't get here
		}
		if i < len(args)-1 {
			expectTokens <- []tokenType{tokenComma}
		}
	}
	close(expectTokens)
	return <-ret
}

func (p *parser) parseArgs(inst string, expectTokens <-chan []tokenType, ret chan<- parseFn) {
	item := parseItem{
		instruction: inst,
		typ:         itemInst,
		line:        p.line,
	}
	for types := range expectTokens {
		token := <-p.tokens
		matched := false
		for _, typ := range types {
			if token.typ != typ {
				continue
			}
			// Matched
			matched = true
			switch token.typ {
			case tokenInteger:
				if strings.HasPrefix(token.val, "0X") {
					token.val = strings.ToLower(token.val)
				}
				i, err := strconv.ParseInt(token.val, 0, 32)
				if err != nil {
					ret <- p.errorf("failed to parse integer %q: %s",
						token.val, err.Error())
					return
				}
				item.imme = int(i)
			case tokenRegister:
				item.registers = append(item.registers, token.val)
			case tokenLabel:
				item.label = token.val
			default:
				// Skip
			}
			break // This position is completed
		}
		// Not match any one
		if !matched {
			var s []string
			for _, t := range types {
				s = append(s, t.String())
			}
			ret <- p.errorf("unexpected token %q(type %q), expect %q",
				token.val, token.typ, strings.Join(s, " | "))
			return
		}
	}
	p.items <- item
	ret <- parseEndline
}

func (p *parser) parseDir(dir string) parseFn {
	item := parseItem{
		typ:       itemDir,
		directive: dir,
		line:      p.line,
	}
	var t token
	switch dir {
	case "byte", "half", "word":
		width := 32
		switch dir {
		case "byte":
			width = 8
		case "half":
			width = 16
		}
		var data []int
	END_LOOP:
		for {
			t = <-p.tokens
			switch t.typ {
			case tokenByte:
				data = append(data, int(t.val[0]))
			case tokenInteger:
				i, err := strconv.ParseInt(t.val, 0, width)
				if err != nil {
					return p.errorf("parse %q: %s", t.val, err.Error())
				}
				data = append(data, int(i))
			default:
				return p.errorf("unexpected token %q(type %q), expect %q",
					t.val, t.typ, "tokenByte | tokenInteger")
			}
			t = <-p.tokens
			switch t.typ {
			case tokenComma:
				continue
			case tokenEndline, tokenEOF:
				break END_LOOP
			default:
				return p.errorf("unexpected token %q(type %q)",
					t.val, t.typ)
			}
		}
		item.data = data
	case "align", "space":
		t = <-p.tokens
		switch t.typ {
		case tokenInteger:
			i, err := strconv.ParseInt(t.val, 0, 8)
			if err != nil {
				return p.errorf("parse %q: %s", t.val, err.Error())
			}
			item.data = i
		default:
			return p.errorf("unexpected token %q(type %q), expect %q",
				t.val, t.typ, tokenInteger)
		}
	case "ascii", "asciiz":
		t = <-p.tokens
		switch t.typ {
		case tokenString:
			item.data = t.val
		default:
			return p.errorf("unexpected token %q(type %q), expect %q",
				t.val, t.typ, tokenString)
		}
	case "globl":
		t = <-p.tokens
		switch t.typ {
		case tokenLabel:
			item.data = t.val
		default:
			return p.errorf("unexpected token %q(type %q), expect %q",
				t.val, t.typ, tokenLabel)
		}
	case "data", "text":
		// Do nothing
	default:
		return p.errorf("invalid directive %q", t.val)
	}
	p.items <- item
	if t.typ == tokenEOF {
		p.items <- parseItem{
			typ: itemEOF,
		}
		return nil
	}
	if t.typ == tokenEndline {
		return parseStart
	}
	return parseEndline
}
