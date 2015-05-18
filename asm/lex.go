package asm

import (
	"bufio"
	"fmt"
	"strings"
	"unicode/utf8"
)

// tokenType represents type of lex tokens
//go:generate stringer -type=tokenType
type tokenType int

const (
	eof        rune      = -1
	tokenError tokenType = 1 << iota
	tokenEOF
	tokenInstruction     // mips instruction
	tokenInteger         // immediate integer
	tokenRegister        // register token
	tokenComma           // comma, ','
	tokenColon           // colon, ':'
	tokenByte            // 'A'
	tokenString          // "Hello"
	tokenDirective       // ".data"
	tokenLeftParenthese  // '('
	tokenRightParenthese // ')'
	tokenLabel           // label reference
	tokenLabelDef        // label definition, "Next:"
	tokenEndline         // end line
)

// token represents a token, it holds type and value of lex items
type token struct {
	typ tokenType
	val string
}

// lexer holds the state of scanner
type lexer struct {
	name   string
	r      *bufio.Reader // read input from here
	buf    []byte        // buffer for current scanned string
	length int           // length of buffer
	eof    bool          // reach end of input
	tokens chan token    // channel of scanned tokens
}

// stateFn is a state of a state machine
// that returns next state as a function
type stateFn func(*lexer) stateFn

// String returns the string representation of token
func (i token) String() string {
	switch i.typ {
	case tokenEOF:
		return "EOF"
	case tokenError:
		return i.val
	}
	return fmt.Sprintf("%q", i.val)
}

// lex launches a state machine as a goruntine
func lex(r *bufio.Reader) chan token {
	l := &lexer{
		r:      r,
		tokens: make(chan token),
	}
	go l.run(lexInline)
	return l.tokens
}

// run lexes the input until the state is nil
func (l *lexer) run(initState stateFn) {
	for state := initState; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

// emit sends a token to channel
func (l *lexer) emit(t tokenType) {
	l.tokens <- token{
		typ: t,
		val: string(l.buf),
	}
	l.buf = []byte{}
}

// next read next rune in input, increase current position
func (l *lexer) next() rune {
	r, n, _ := l.r.ReadRune()
	if n > 0 {
		l.buf = append(l.buf, []byte(string(r))...)
		return r
	}
	l.eof = true
	return eof
}

// ignore skips over pending input before current position
func (l *lexer) ignore() {
	l.buf = []byte{}
}

// backup steps back one rune
// Can be called only once per call to next
func (l *lexer) backup() {
	if !l.eof {
		l.r.UnreadRune()
		_, n := utf8.DecodeLastRune(l.buf)
		l.buf = l.buf[:len(l.buf)-n]
	} else {
		l.eof = false
	}
}

// peek returns next rune in input but not consume it
func (l *lexer) peek() rune {
	r, _, err := l.r.ReadRune()
	if err != nil {
		return eof
	}
	l.r.UnreadRune()
	return r
}

// accept consume next rune if it's from valid set
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consume a run of runes from valid set
func (l *lexer) acceptRun(valid string) bool {
	run := false
	for strings.IndexRune(valid, l.next()) >= 0 {
		run = true
	}
	l.backup()
	return run
}

// curValue returns current value
func (l *lexer) curValue() string {
	return string(l.buf)
}

// errorf emit a error token and terminate the scan
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- token{
		typ: tokenError,
		val: fmt.Sprintf(format, args...),
	}
	return nil
}
