package cpu

import "github.com/paulloz/ohboi/bits"

const (
	NFlag     = 0x40
	ZFlag     = 0x80
	HFlag     = 0x20
	CarryFlag = 0x10
)

func (cpu *CPU) GetZFlag() bool {
	return bits.Test(7, cpu.F.Get())
}

func (cpu *CPU) SetZFlag(v bool) {
	if v {
		cpu.F.Set(bits.Set(7, cpu.F.Get()))
	} else {
		cpu.F.Set(bits.Reset(7, cpu.F.Get()))
	}
}

func (cpu *CPU) GetNFlag() bool {
	return bits.Test(6, cpu.F.Get())
}

func (cpu *CPU) SetNFlag(v bool) {
	if v {
		cpu.F.Set(bits.Set(6, cpu.F.Get()))
	} else {
		cpu.F.Set(bits.Reset(6, cpu.F.Get()))
	}
}

func (cpu *CPU) GetHFlag() bool {
	return bits.Test(5, cpu.F.Get())
}

func (cpu *CPU) SetHFlag(v bool) {
	if v {
		cpu.F.Set(bits.Set(5, cpu.F.Get()))
	} else {
		cpu.F.Set(bits.Reset(5, cpu.F.Get()))
	}
}

func (cpu *CPU) GetCFlag() bool {
	return bits.Test(4, cpu.F.Get())
}

func (cpu *CPU) SetCFlag(v bool) {
	if v {
		cpu.F.Set(bits.Set(4, cpu.F.Get()))
	} else {
		cpu.F.Set(bits.Reset(4, cpu.F.Get()))
	}
}
