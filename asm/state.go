package asm

import "unicode"

// lexInline lexes identifiers, numbers, registers and comment
func lexInline(l *lexer) stateFn {
	for {
		r := l.peek()
		switch r {
		case eof, '\n':
			return lexEndline
		case '#':
			return lexComment
		case '-', '+':
			return lexNumber
		case '$':
			return lexRegister
		case '.':
			return lexDirective
		case '\'':
			return lexByte
		case '"':
			return lexString
		case ',':
			l.next()
			l.emit(tokenComma)
		case ':':
			l.next()
			l.emit(tokenColon)
		case '(':
			l.next()
			l.emit(tokenLeftParenthese)
		case ')':
			l.next()
			l.emit(tokenRightParenthese)
		default:
			switch {
			case isLetter(r):
				return lexIdentifier
			case unicode.IsSpace(r):
				l.next()
				l.ignore()
			case unicode.IsDigit(r):
				return lexNumber
			default:
				return l.errorf("bad syntax: %q", l.curValue())
			}
		}
	}
}

// lexComment skips comment
func lexComment(l *lexer) stateFn {
	// skip commit
	for r0 := l.next(); r0 != '\n' && r0 != eof; r0 = l.next() {
	}
	l.backup()
	l.ignore()
	return lexEndline
}

// lexEndline emits endline or eof, then go to next line if possible
func lexEndline(l *lexer) stateFn {
	r := l.next()
	switch r {
	case eof:
		l.emit(tokenEOF)
		return nil
	case '\n':
		l.emit(tokenEndline)
		return lexInline
	default:
		return l.errorf("state error at %q", l.curValue())
	}
}

// lexIdentifier lexes instructions and labels
// Any identifier not in instruction set is treated as label
func lexIdentifier(l *lexer) stateFn {
	var r rune
	for r = l.next(); isLetterDigit(r); r = l.next() {
	}
	l.backup()
	switch {
	case unicode.IsSpace(r) || r == '#' || r == eof:
		if _, ok := instructionTable[l.curValue()]; ok {
			l.emit(tokenInstruction)
		} else {
			l.emit(tokenLabel)
		}
	case r == ':':
		l.emit(tokenLabelDef)
		l.next()
		l.ignore()
	default:
		l.next()
		return l.errorf("invalid identifier %q", l.curValue())
	}
	return lexInline
}

// lexNumber lexes numbers in decimal or hex format
func lexNumber(l *lexer) stateFn {
	// Optional leading sign.
	l.accept("+-")
	// Is it hex?
	digits := "0123456789"
	if l.accept("0") && l.accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}
	l.acceptRun(digits)
	l.emit(tokenInteger)
	return lexInline
}

// lexRegister lexes 32 mips registers
func lexRegister(l *lexer) stateFn {
	r := l.next() // '$'
	for r = l.next(); isLetterDigit(r); r = l.next() {
	}
	l.backup()
	if _, ok := registerTable[l.curValue()[1:]]; ok {
		l.emit(tokenRegister)
		return lexInline
	}
	return l.errorf("invalid register name: %q", l.curValue())
}

// lexString lexes double-quoted string
func lexString(l *lexer) stateFn {
	// skip leading '"'
	r := l.next()
	l.ignore()
	for r = l.next(); r != '"' && isPrintByte(r); r = l.next() {
	}
	if r != '"' {
		return l.errorf("bad string syntax: %q, expect \"", l.curValue())
	}
	l.backup()
	l.emit(tokenString)
	l.next()
	l.ignore()
	return lexInline
}

// lexByte lexes single-quoted character
func lexByte(l *lexer) stateFn {
	// skip leading '''
	r := l.next()
	l.ignore()
	r = l.next()
	if r == '\'' || !isPrintByte(r) {
		return l.errorf("invalid byte: %q", l.curValue())
	}
	r = l.next()
	if r != '\'' {
		return l.errorf("bad byte syntax: %q, expect '", l.curValue())
	}
	l.backup()
	l.emit(tokenByte)
	l.next()
	l.ignore()
	return lexInline
}

// lexDirective lexes directive
func lexDirective(l *lexer) stateFn {
	// skip leading '.'
	r := l.next()
	l.ignore()
	r = l.peek()
	if !isLetter(r) {
		return l.errorf("invalid directive %q, must start with a letter",
			l.curValue())
	}
	for r = l.next(); isLetterDigit(r); r = l.next() {
	}
	l.backup()
	if unicode.IsSpace(r) || r == '#' || r == eof {
		l.emit(tokenDirective)
	} else {
		l.next()
		return l.errorf("invalid directive syntax: %q", l.curValue())
	}
	return lexInline
}

func isLetter(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isLetterDigit(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

func isPrintByte(r rune) bool {
	return r >= 0x20 && r <= 0x7e
}
