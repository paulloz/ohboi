package io

const (
	// Joy pad
	P1 = 0x00

	// Serial
	SB = 0x01
	SC = 0x02

	// Timing
	DIV  = 0x04
	TIMA = 0x05
	TMA  = 0x06
	TAC  = 0x07

	// Interrupt
	IF = 0x0f
	IE = 0xff

	// Sound
	NR10 = 0x10
	NR11 = 0x11
	NR12 = 0x12
	NR13 = 0x13
	NR14 = 0x14
	NR21 = 0x16
	NR22 = 0x17
	NR23 = 0x18
	NR24 = 0x19
	NR30 = 0x1a
	NR31 = 0x1b
	NR32 = 0x1c
	NR33 = 0x1d
	NR34 = 0x1e
	NR41 = 0x20
	NR42 = 0x21
	NR43 = 0x22
	NR44 = 0x23
	NR50 = 0x24
	NR51 = 0x25
	NR52 = 0x26
	WAVE = 0x30

	// Graphics
	LDCD = 0x40
	STAT = 0x41
	SCY  = 0x42
	SCX  = 0x43
	LY   = 0x44
	LYC  = 0x45
	DMA  = 0x46
	BGP  = 0x47
	OBP0 = 0x48
	OBP1 = 0x49
	WY   = 0x4a
	WX   = 0x4b
)

type MemoryRegister struct {
	value  uint8
	unused uint8
}

func (r *MemoryRegister) Read() uint8 {
	return r.value | r.unused
}

func (r *MemoryRegister) Write(value uint8) {
	r.value = value &^ r.unused
}

func newMemoryRegister(value uint8, unused ...uint8) *MemoryRegister {
	if len(unused) > 0 {
		return &MemoryRegister{value: value, unused: unused[0]}
	}
	return &MemoryRegister{value: value}
}

type MappedRegister struct {
	getter func() uint8
	setter func(uint8)
}

func (r *MappedRegister) Read() uint8 {
	return r.getter()
}

func (r *MappedRegister) Write(v uint8) {
	r.setter(v)
}

func NewMappedRegister(getter func() uint8, setter func(uint8)) *MappedRegister {
	return &MappedRegister{getter: getter, setter: setter}
}
