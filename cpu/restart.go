package cpu

import (
	op "github.com/paulloz/ohboi/cpu/opcodes"
	"github.com/paulloz/ohboi/memory"
)

func newReset(offset uint8) Instruction {
	return NewInstruction(
		DecodeInstruction,
		NoopInstruction,
		func(cpu *CPU, mem *memory.Memory) error {
			cpu.PushByte(uint8(cpu.PC >> 8))
			return nil
		},
		func(cpu *CPU, mem *memory.Memory) error {
			cpu.PushByte(uint8(cpu.PC))
			cpu.PC = uint16(offset)
			return nil
		},
	)
}

func init() {
	RegisterInstructions(map[uint8]Instruction{
		op.RST_00H: newReset(0x00),
		op.RST_08H: newReset(0x08),
		op.RST_10H: newReset(0x10),
		op.RST_18H: newReset(0x18),
		op.RST_20H: newReset(0x20),
		op.RST_28H: newReset(0x28),
		op.RST_30H: newReset(0x30),
		op.RST_38H: newReset(0x38),
	})
}
