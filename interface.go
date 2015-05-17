package mips

func Assemble(input string) ([]byte, error) {
	_, ch := parse(input)
	return assemble(ch)
}
