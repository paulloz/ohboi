package io

type IO struct {
	registers map[uint8]Register
}

type Register interface {
	Read() uint8
	Write(uint8)
}

func (io *IO) Read(address uint8) uint8 {
	return io.registers[address].Read()
}

func (io *IO) Write(address uint8, value uint8) {
	io.registers[address].Write(value)
}

func (io *IO) MapRegister(address uint8, getter func() uint8, setter func(uint8)) {
	io.registers[address] = &MappedRegister{getter: getter, setter: setter}
}

func (io *IO) GetRegister(address uint8) Register {
	return io.registers[address]
}

func NewIO() *IO {
	io := &IO{registers: make(map[uint8]Register, 0xff)}

	for i := uint16(0x0000); i <= 0x00ff; i++ {
		io.registers[uint8(i&0xff)] = newMemoryRegister(0)
	}

	return io
}
