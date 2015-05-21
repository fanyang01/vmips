package mips

import "errors"

type addrSeg int

const (
	TEXT_ADDRESS   = 0x8000000
	DATA_ADDRESS   = 0x8100000
	STACK_ADDRESS  = 0x7FFFFFFF
	MIN_STACK_ADDR = 0x7F000000
	MAX_DATA_ADDR  = 0x8F00000
	textSegment    = 1 << iota
	dataSegment
	stackSegment
)

/*
   +-----------+   STACK_ADDRESS
   |   stack   |
   +-----------+
   |  unmaped  |
   |   ...     |
   +-----------+
   |   data    |
   +-----------+   DATA_ADDRESS
   |   text    |
   +-----------+   TEXT_ADDRESS
*/

type virtualMemory struct {
	text, data, stack []byte
}

type registerFile struct {
	general    [32]int
	HI, LO, PC int
}

type Machine struct {
	m    *virtualMemory
	r    *registerFile
	exit bool
}

func NewMachine() *Machine {
	return &Machine{
		m: &virtualMemory{
			text:  make([]byte, 1<<12),
			data:  make([]byte, 1<<12),
			stack: make([]byte, 1<<12),
		},
		r: new(registerFile),
	}
}

func (rf *registerFile) read(id int) int {
	id &= 0x1F
	if id == 0 {
		return 0
	}
	return rf.general[id]
}

func (rf *registerFile) write(id int, value int) {
	id &= 0x1F
	if id == 0 {
		return
	}
	rf.general[id] = value
}

func (m *virtualMemory) read(addr int) (byte, error) {
	actual, seg, err := m.transfer(addr)
	if err != nil {
		return 0, err
	}
	switch seg {
	case textSegment:
		return m.text[actual], nil
	case dataSegment:
		return m.data[actual], nil
	case stackSegment:
		return m.stack[actual], nil
	default:
		// shouldn't get here
		panic("something wrong...")
	}
}

func (m *virtualMemory) write(addr int, value byte) error {
	actual, seg, err := m.transfer(addr)
	if err != nil {
		return err
	}
	switch seg {
	case textSegment:
		m.text[actual] = value
	case dataSegment:
		m.data[actual] = value
	case stackSegment:
		m.stack[actual] = value
	default:
		panic("something wrong...")
	}
	return nil
}

func (m *virtualMemory) readHalf(addr int) (int, error) {
	b1, err := m.read(addr)
	if err != nil {
		return 0, err
	}
	b2, err := m.read(addr + 1)
	if err != nil {
		return 0, err
	}
	i := int(b2)
	i <<= 8
	i &= int(b1)
	return i, nil
}

func (m *virtualMemory) writeHalf(addr int, value int) error {
	err := m.write(addr, byte(value&0xFF))
	if err != nil {
		return err
	}
	return m.write(addr+1, byte((value>>8)&0xFF))
}

func (m *virtualMemory) readWord(addr int) (int, error) {
	b0, err := m.read(addr)
	if err != nil {
		return 0, err
	}
	b1, err := m.read(addr + 1)
	if err != nil {
		return 0, err
	}
	b2, err := m.read(addr + 2)
	if err != nil {
		return 0, err
	}
	b3, err := m.read(addr + 3)
	if err != nil {
		return 0, err
	}
	i := int(b3)
	i = i << 8
	i |= int(b2)
	i = i << 8
	i |= int(b1)
	i = i << 8
	i |= int(b0)
	return i, nil
}

func (m *virtualMemory) writeWord(addr int, value int) error {
	err := m.write(addr, byte(value&0xFF))
	if err != nil {
		return err
	}
	err = m.write(addr+1, byte((value>>8)&0xFF))
	if err != nil {
		return err
	}
	err = m.write(addr+2, byte((value>>16)&0xFF))
	if err != nil {
		return err
	}
	return m.write(addr+3, byte((value>>24)&0xFF))
}

func (m *virtualMemory) writeBytes(addr int, s []byte) error {
	for i := 0; i < len(s); i++ {
		err := m.write(addr+i, s[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *virtualMemory) transfer(virtual int) (int, addrSeg, error) {
	switch {
	case virtual >= TEXT_ADDRESS && virtual < DATA_ADDRESS:
		actual := virtual - TEXT_ADDRESS
		for actual >= len(m.text) {
			text := make([]byte, len(m.text)<<1)
			copy(text, m.text)
			m.text = text
		}
		return actual, textSegment, nil
	case virtual >= DATA_ADDRESS && virtual < MAX_DATA_ADDR:
		actual := virtual - DATA_ADDRESS
		for actual >= len(m.data) {
			data := make([]byte, len(m.data)<<1)
			copy(data, m.data)
			m.data = data
		}
		return actual, dataSegment, nil
	case virtual >= MIN_STACK_ADDR && virtual <= STACK_ADDRESS:
		actual := STACK_ADDRESS - virtual
		for actual >= len(m.stack) {
			stack := make([]byte, len(m.stack)<<1)
			copy(stack, m.stack)
			m.stack = stack
		}
		return actual, stackSegment, nil
	default:
		return 0, 0, errors.New("Segmentfault")
	}
}
