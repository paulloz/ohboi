package cpu

func (cpu *CPU) SetZFlag(v bool) {
	if v {
		cpu.AF.SetLo(cpu.AF.Lo() | 0x80) // 0x80 -> 1000 0000
	} else {
		cpu.AF.SetLo(cpu.AF.Lo() & 0x7F) // 0x7F -> 0111 1111
	}
}

func (cpu *CPU) SetNFlag(v bool) {
	if v {
		cpu.AF.SetLo(cpu.AF.Lo() | 0x40) // 0x40 -> 0100 0000
	} else {
		cpu.AF.SetLo(cpu.AF.Lo() & 0xBF) // 0xBF -> 1011 1111
	}
}

func (cpu *CPU) SetHFlag(v bool) {
	if v {
		cpu.AF.SetLo(cpu.AF.Lo() | 0x20) // 0x20 -> 0010 0000
	} else {
		cpu.AF.SetLo(cpu.AF.Lo() & 0xDF) // 0xDF -> 1101 1111
	}
}

func (cpu *CPU) SetCFlag(v bool) {
	if v {
		cpu.AF.SetLo(cpu.AF.Lo() | 0x10) // 0x10 -> 0001 0000
	} else {
		cpu.AF.SetLo(cpu.AF.Lo() & 0xEF) // 0xEF -> 1110 1111
	}
}
