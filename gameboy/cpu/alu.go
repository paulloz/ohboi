package cpu

func (cpu *CPU) And(out func(uint8), a uint8, b uint8) {
	result := b & a
	out(result)

	cpu.SetZFlag(result == 0)
	cpu.SetNFlag(false)
	cpu.SetHFlag(true)
	cpu.SetCFlag(false)
}

func (cpu *CPU) Inc(out func(uint8), in uint8) {
	new := in + 1
	out(new)

	cpu.SetZFlag(new == 0)
	cpu.SetNFlag(false)
	cpu.SetHFlag(false) // TODO: Implement HalfCarry Flag
}
