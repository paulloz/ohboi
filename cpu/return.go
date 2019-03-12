package cpu

import (
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func EnableInterrupts(cpu *CPU, mem *memory.Memory) error {
	cpu.EnableInterrupts()
	return nil
}

func PopByte(v *uint8) MicroInstruction {
	return func(cpu *CPU, mem *memory.Memory) error {
		*v = cpu.PopByte()
		return nil
	}
}

func newRetInstruction(condition func(cpu *CPU, mem *memory.Memory) bool, jmpCombine MicroInstruction) Instruction {
	micros := []MicroInstruction{DecodeInstruction}

	if condition != nil {
		micros = append(micros, SkipIf(condition))
	}

	var jmp MicroInstruction = func(cpu *CPU, mem *memory.Memory) error {
		cpu.PC = AddressImmediateOperand.Address()
		return nil
	}
	if jmpCombine != nil {
		jmp.Combine(jmpCombine)
	}

	micros = append(micros, PopByte(&AddressImmediateOperand.lo),
		PopByte(&AddressImmediateOperand.hi),
		jmp)

	return NewInstruction(micros...)
}

func init() {
	RegisterInstructions(map[uint8]Instruction{
		op.RET: newRetInstruction(nil, nil),

		op.RETI: newRetInstruction(nil, EnableInterrupts),

		op.RET_C: newRetInstruction(func(cpu *CPU, mem *memory.Memory) bool { return cpu.GetCFlag() }, nil),

		op.RET_NC: newRetInstruction(func(cpu *CPU, mem *memory.Memory) bool { return !cpu.GetCFlag() }, nil),

		op.RET_Z: newRetInstruction(func(cpu *CPU, mem *memory.Memory) bool { return cpu.GetZFlag() }, nil),

		op.RET_NZ: newRetInstruction(func(cpu *CPU, mem *memory.Memory) bool { return !cpu.GetZFlag() }, nil),
	})
}
