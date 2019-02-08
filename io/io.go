package io

type IO struct {
	registers map[uint8]*IORegister
}

func (io *IO) Read(address uint8) uint8 {
	return io.registers[address].value
}

func (io *IO) Write(address uint8, value uint8) {
	io.registers[address].value = value
}

func NewIO() *IO {
	io := &IO{registers: make(map[uint8]*IORegister, 0xff)}

	for i := uint16(0x0000); i <= 0x00ff; i++ {
		io.registers[uint8(i&0xff)] = newIORegister(0)
	}

	return io
}
