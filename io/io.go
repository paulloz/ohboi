package io

import (
	"github.com/paulloz/ohboi/bits"
)

type IO struct {
	registers map[uint8]Register
}

type Register interface {
	Read() uint8
	Write(uint8)
}

func (io *IO) Read(address uint8) uint8 {
	register, found := io.registers[address]
	if found {
		return register.Read()
	}
	return 0xff
}

func (io *IO) ReadBit(address uint8, bit uint8) bool {
	return bits.Test(bit, io.Read(address))
}

func (io *IO) Write(address uint8, value uint8) {
	register, found := io.registers[address]
	if found {
		register.Write(value)
	}
}

func (io *IO) WriteBit(address uint8, bit uint8, value bool) {
	if value {
		io.SetBit(address, bit)
	} else {
		io.ResetBit(address, bit)
	}
}

func (io *IO) SetBit(address uint8, bit uint8) {
	io.Write(address, bits.Set(bit, io.Read(address)))
}

func (io *IO) ResetBit(address uint8, bit uint8) {
	io.Write(address, bits.Reset(bit, io.Read(address)))
}

func (io *IO) MapRegister(address uint8, getter func() uint8, setter func(uint8)) {
	io.registers[address] = &MappedRegister{getter: getter, setter: setter}
}

func (io *IO) GetRegister(address uint8) Register {
	return io.registers[address]
}

func NewIO() *IO {
	io := &IO{registers: make(map[uint8]Register, 0xff)}

	io.registers[P1] = newMemoryRegister(0, 0xc0)

	// Serial
	io.registers[SB] = newMemoryRegister(0)
	io.registers[SC] = newMemoryRegister(0, 0x7e)

	// Timing
	io.registers[DIV] = newMemoryRegister(0)
	io.registers[TIMA] = newMemoryRegister(0)
	io.registers[TMA] = newMemoryRegister(0)
	io.registers[TAC] = newMemoryRegister(0, 0xf8)

	// Interrupt
	io.registers[IF] = newMemoryRegister(0, 0xe0)
	io.registers[IE] = newMemoryRegister(0)

	// Sound
	io.registers[NR10] = newMemoryRegister(0, 0x80)
	io.registers[NR11] = newMemoryRegister(0)
	io.registers[NR12] = newMemoryRegister(0)
	io.registers[NR13] = newMemoryRegister(0)
	io.registers[NR14] = newMemoryRegister(0)
	io.registers[NR21] = newMemoryRegister(0)
	io.registers[NR22] = newMemoryRegister(0)
	io.registers[NR23] = newMemoryRegister(0)
	io.registers[NR24] = newMemoryRegister(0)
	io.registers[NR30] = newMemoryRegister(0, 0x7f)
	io.registers[NR31] = newMemoryRegister(0)
	io.registers[NR32] = newMemoryRegister(0, 0x9f)
	io.registers[NR33] = newMemoryRegister(0)
	io.registers[NR34] = newMemoryRegister(0)
	io.registers[NR41] = newMemoryRegister(0, 0xc0)
	io.registers[NR42] = newMemoryRegister(0)
	io.registers[NR43] = newMemoryRegister(0)
	io.registers[NR44] = newMemoryRegister(0, 0x3f)
	io.registers[NR50] = newMemoryRegister(0)
	io.registers[NR51] = newMemoryRegister(0)
	io.registers[NR52] = newMemoryRegister(0, 0x70)
	io.registers[WAVE] = newMemoryRegister(0)

	// Graphics
	io.registers[LDCD] = newMemoryRegister(0)
	io.registers[STAT] = newMemoryRegister(0, 0x80)
	io.registers[SCY] = newMemoryRegister(0)
	io.registers[SCX] = newMemoryRegister(0)
	io.registers[LY] = newMemoryRegister(0)
	io.registers[LYC] = newMemoryRegister(0)
	io.registers[DMA] = newMemoryRegister(0)
	io.registers[BGP] = newMemoryRegister(0)
	io.registers[OBP0] = newMemoryRegister(0)
	io.registers[OBP1] = newMemoryRegister(0)
	io.registers[WY] = newMemoryRegister(0)
	io.registers[WX] = newMemoryRegister(0)

	io.registers[BOOTROM] = newMemoryRegister(0)

	/*
		// Undocumented GBC registers for GBC
		io.registers[0x72] = newMemoryRegister(0)
		io.registers[0x73] = newMemoryRegister(0)
		io.registers[0x74] = newMemoryRegister(0, 0xff)
		io.registers[0x75] = newMemoryRegister(0, 0x8f)
		io.registers[0x76] = ZeroRegister
		io.registers[0x77] = ZeroRegister
	*/

	return io
}
