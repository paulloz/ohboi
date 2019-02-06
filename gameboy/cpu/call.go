package cpu

func (cpu *CPU) Call(nn uint16) {
	// TODO: push PC on stack
	cpu.PC = nn
}
