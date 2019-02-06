package cpu

func (cpu *CPU) Jump(nn uint16) {
	cpu.PC = nn
}
